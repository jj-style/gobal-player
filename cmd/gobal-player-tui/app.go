package main

import (
	"context"
	"fmt"
	"maps"
	"net/http"
	"os"

	"dario.cat/mergo"
	"github.com/gdamore/tcell/v2"
	"github.com/jj-style/gobal-player/cmd/gobal-player-tui/internal/utils/text"
	"github.com/jj-style/gobal-player/pkg/audioplayer"
	"github.com/jj-style/gobal-player/pkg/globalplayer"
	"github.com/jj-style/gobal-player/pkg/globalplayer/models"
	"github.com/jj-style/gobal-player/pkg/resty"
	"github.com/rivo/tview"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

type app struct {
	tv     *tview.Application
	player audioplayer.Player
	gp     globalplayer.GlobalPlayer
	hc     *http.Client
	cache  resty.Cache[[]byte]

	stations     []*models.Station
	stationsList *tview.List

	shows     []*models.Show
	showsList *tview.List

	episodes     []*models.Episode
	episodesList *tview.List

	stationSlug string
	showId      string

	globalKbdShortcuts map[string]string
	kbdShortcuts       map[int]map[string]string
	currentPane        int
	helpText           *tview.TextView

	streaming streamingData
}

type streamingData struct {
	item   streamItem
	stop   chan struct{}
	change func()
}

type streamItem struct {
	Name string
	Id   string
}

func NewApp(gp globalplayer.GlobalPlayer, player audioplayer.Player, hc *http.Client, cache resty.Cache[[]byte]) *app {
	tv := tview.NewApplication().EnableMouse(false)

	a := &app{
		tv:                 tv,
		gp:                 gp,
		player:             player,
		hc:                 hc,
		cache:              cache,
		streaming:          streamingData{},
		kbdShortcuts:       make(map[int]map[string]string),
		globalKbdShortcuts: map[string]string{"r": "refresh"},
		currentPane:        1,
	}

	// ensure log.Fatal stops tui so doesn't mess up terminal
	log.StandardLogger().ExitFunc = func(code int) {
		a.tv.Stop()
		os.Exit(code)
	}

	a.initTui()
	return a
}

func (a *app) Run() error {
	return a.tv.Run()
}

func (a *app) stream(list *tview.List, items []streamItem, url string) {
	// change previous text if needed
	if a.streaming.change != nil {
		a.streaming.change()
	}
	currIdx := list.GetCurrentItem()
	currItem := items[currIdx]
	if a.streaming.item.Id != currItem.Id {
		// start playing new audio
		log.WithField("url", url).Debug("start streaming")
		stop, err := a.player.Play(url)
		if err != nil {
			log.Fatal(err)
		}

		// update state
		a.streaming.item = streamItem{Id: currItem.Id, Name: currItem.Name}
		a.streaming.stop = stop
		list.SetItemText(currIdx, fmt.Sprintf("%s*", currItem.Name), "")
		a.streaming.change = func() { list.SetItemText(currIdx, currItem.Name, "") }
	} else {
		// stop currently playing
		log.Debug("stop streaming")
		a.streaming.item = streamItem{}
		close(a.streaming.stop)
		list.SetItemText(currIdx, currItem.Name, "")
	}
}

// initialise the views in the applications
// setup the widgets and configure their event handlers
func (a *app) initViews() {
	stList := tview.NewList().ShowSecondaryText(false).
		SetChangedFunc(func(idx int, mainText, secondaryText string, shortcut rune) {
			a.stationSlug = a.stations[idx].Slug
			go func() {
				a.tv.QueueUpdateDraw(func() {
					a.getShowsList(a.stations[idx].Slug)
				})
			}()
		})
	stList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		r := event.Rune()
		k := event.Key()
		switch {
		case r == 'j', k == tcell.KeyDown:
			stList.SetCurrentItem((stList.GetCurrentItem() + 1) % stList.GetItemCount())
		case r == 'k', k == tcell.KeyUp:
			stList.SetCurrentItem((stList.GetCurrentItem() - 1) % stList.GetItemCount())
		case r == 'f':
			viper.Set("favourite", (a.stations[stList.GetCurrentItem()].Slug))
		case k == tcell.KeyEnter:
			station := a.stations[stList.GetCurrentItem()]
			a.stream(a.stationsList, lo.Map(a.stations, func(item *models.Station, _ int) streamItem { return streamItem{Name: item.Name, Id: item.Id} }), station.StreamUrl)
		}
		return nil
	})
	a.stationsList = stList

	showList := tview.NewList().ShowSecondaryText(false).
		SetChangedFunc(func(idx int, mainText, secondaryText string, shortcut rune) {
			a.showId = a.shows[idx].Id
			go func() {
				a.tv.QueueUpdateDraw(func() {
					a.getEpisodes(a.stationSlug, a.shows[idx].Id)
				})
			}()
		})
	showList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		r := event.Rune()
		k := event.Key()
		switch {
		case r == 'j', k == tcell.KeyDown:
			showList.SetCurrentItem((showList.GetCurrentItem() + 1) % showList.GetItemCount())
		case r == 'k', k == tcell.KeyUp:
			showList.SetCurrentItem((showList.GetCurrentItem() - 1) % showList.GetItemCount())
		}
		return nil
	})
	a.showsList = showList

	epList := tview.NewList().ShowSecondaryText(false)
	epList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if epList.GetItemCount() == 0 {
			return nil
		}

		r := event.Rune()
		k := event.Key()
		switch {
		case r == 'j', k == tcell.KeyDown:
			epList.SetCurrentItem((epList.GetCurrentItem() + 1) % epList.GetItemCount())
		case r == 'k', k == tcell.KeyUp:
			epList.SetCurrentItem((epList.GetCurrentItem() - 1) % epList.GetItemCount())
		case r == 'd':
			mainText, _ := epList.GetItemText(epList.GetCurrentItem())
			curr := epList.GetCurrentItem()
			epList.SetItemText(epList.GetCurrentItem(), mainText+" [blue](downloading...)", "")
			go func() {
				if err := resty.DownloadFile(a.hc, fmt.Sprintf("%s.m4a", mainText), a.episodes[curr].StreamUrl); err != nil {
					log.Fatal(err)
				}
				a.tv.QueueUpdateDraw(func() {
					epList.SetItemText(curr, mainText, "")
				})
			}()
		case k == tcell.KeyEnter:
			ep := a.episodes[epList.GetCurrentItem()]
			a.stream(a.episodesList, lo.Map(a.episodes, func(item *models.Episode, _ int) streamItem {
				mainText := fmt.Sprintf("%s - %s - %s", item.Name, item.Aired.Format("Mon 2006-01-02"), item.Availability)
				return streamItem{Name: mainText, Id: item.Id}
			}), ep.StreamUrl)
		}
		return nil
	})
	a.episodesList = epList

	// create keyboard shortcut help
	a.kbdShortcuts[1] = map[string]string{"\u21B5": "play/pause", "f": "favourite"}
	a.kbdShortcuts[2] = map[string]string{}
	a.kbdShortcuts[3] = map[string]string{"\u21B5": "play/pause", "d": "download"}

	a.helpText = tview.NewTextView()
}

// initialise the TUI app
func (a *app) initTui() {
	a.initViews()

	a.prefetch()

	a.render()

	// base input handlers
	a.tv.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			// Exit the application
			a.tv.Stop()
			return nil
		}

		switch event.Rune() {
		case 49: // 1
			a.tv.SetFocus(a.stationsList)
			a.currentPane = 1
			a.updateHelpText()
			return nil
		case 50: // 2
			a.tv.SetFocus(a.showsList)
			a.currentPane = 2
			a.updateHelpText()
			return nil
		case 51: // 2
			a.tv.SetFocus(a.episodesList)
			a.currentPane = 3
			a.updateHelpText()
			return nil
		case 'r': // refresh feeds
			a.prefetch()
			return nil
		}
		return event
	})
}

// render widgets in the TUI
func (a *app) render() {
	stationsFlex := tview.NewFlex()
	stationsFlex.Box.SetBorder(true).SetTitle("[1] Stations")
	stationsFlex.AddItem(a.stationsList, 0, 1, true)
	a.getStationList()

	showsFlex := tview.NewFlex()
	showsFlex.Box.SetBorder(true).SetTitle("[2] Shows")
	showsFlex.AddItem(a.showsList, 0, 1, true)

	epsFlex := tview.NewFlex()
	epsFlex.Box.SetBorder(true).SetTitle("[3] Episodes")
	epsFlex.AddItem(a.episodesList, 0, 1, true)

	a.updateHelpText()

	panels := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(stationsFlex, 0, 1, false).
			AddItem(showsFlex, 0, 1, false),
			0, 1, false).
		AddItem(epsFlex, 0, 1, false)

	root := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(panels, 0, 99, true).
		AddItem(a.helpText, 0, 1, false)

	a.tv.SetRoot(root, true)
	a.tv.SetFocus(a.stationsList)
}

// update the list of stations
func (a *app) getStationList() {
	a.stationsList.Clear()

	stations, err := a.gp.GetStations()
	if err != nil {
		log.Fatal(err)
	}
	a.stations = stations

	for idx, station := range stations {
		text := station.Name
		a.stationsList.AddItem(text, "", 0, nil)

		// if station is favourited, select in the UI
		if station.Slug == viper.GetString("favourite") {
			go func() {
				a.tv.QueueUpdateDraw(func() {
					a.stationsList.SetCurrentItem(idx)
				})
			}()
		}
	}
}

// updates the help text to the keyboard shortcuts for the active pane
func (a *app) updateHelpText() {
	shortcuts := maps.Clone(a.globalKbdShortcuts)
	if err := mergo.Merge(&shortcuts, a.kbdShortcuts[a.currentPane]); err != nil {
		log.Fatal(err)
	}
	a.helpText.SetText(text.FormatHelp(shortcuts))
}

// update the list of shows for a station
func (a *app) getShowsList(slug string) {
	a.showsList.Clear()

	if slug != "" {
		shows, err := a.gp.GetShows(slug)
		if err != nil {
			log.Fatal(err)
		}
		a.shows = shows

		for _, show := range shows {
			a.showsList.AddItem(show.Name, "", 0, nil)
		}
	}
}

// update the catchup list of the given station and show
func (a *app) getEpisodes(slug, id string) {
	a.episodesList.Clear()

	if slug != "" {
		episodes, err := a.gp.GetEpisodes(slug, id)
		if err != nil {
			log.Fatal(err)
		}
		a.episodes = episodes

		for _, show := range episodes {
			text := fmt.Sprintf("%s - %s - %s", show.Name, show.Aired.Format("Mon 2006-01-02"), show.Availability)
			a.episodesList.AddItem(text, "", 0, nil)
		}
	}
}

// pre-fetch all data so it's in the cache
func (a *app) prefetch() {
	// clear the cache to force re-fetch
	a.cache.Clear(context.TODO())

	// fetch everything!
	var g errgroup.Group
	stations, err := a.gp.GetStations()
	if err != nil {
		return
	}

	for _, st := range stations {
		st := st
		g.Go(func() error {
			shows, err := a.gp.GetShows(st.Slug)
			if err != nil {
				return err
			}
			for _, sh := range shows {
				sh := sh
				g.Go(func() error {
					_, err := a.gp.GetEpisodes(st.Slug, sh.Id)
					return err
				})
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		log.Warnf("error fetching all data: %v", err)
	}
}
