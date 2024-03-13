package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/jj-style/gobal-player/cmd/gobal-player-tui/internal/utils/text"
	"github.com/jj-style/gobal-player/pkg/audioplayer"
	"github.com/jj-style/gobal-player/pkg/globalplayer"
	"github.com/jj-style/gobal-player/pkg/globalplayer/models"
	"github.com/jj-style/gobal-player/pkg/resty"
	"github.com/rivo/tview"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
)

type app struct {
	tv     *tview.Application
	player audioplayer.Player
	gp     globalplayer.GlobalPlayer
	hc     *http.Client

	stations     []models.StationBrand
	stationsList *tview.List

	shows     []models.CatchupInfo
	showsList *tview.List

	catchups    []models.Episode
	catchupList *tview.List

	stationSlug string
	showId      string

	// panel -> {"key": "description"}
	kbdShortcuts map[int]map[string]string
	helpText     *tview.TextView

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

func NewApp(gp globalplayer.GlobalPlayer, player audioplayer.Player, hc *http.Client) *app {
	tv := tview.NewApplication().EnableMouse(false)

	a := &app{
		tv:           tv,
		gp:           gp,
		player:       player,
		hc:           hc,
		streaming:    streamingData{},
		kbdShortcuts: make(map[int]map[string]string),
		helpText:     tview.NewTextView(),
	}

	// ensure log.Fatal stops tui so doesn't mess up terminal
	log.StandardLogger().ExitFunc = func(code int) {
		a.tv.Stop()
		os.Exit(code)
	}

	go a.prefetch()

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
		case k == tcell.KeyEnter:
			station := a.stations[stList.GetCurrentItem()]
			a.stream(a.stationsList, lo.Map(a.stations, func(item models.StationBrand, _ int) streamItem { return streamItem{Name: item.Name, Id: item.ID} }), station.NationalStation.StreamURL)
		}
		return nil
	})
	a.stationsList = stList

	showList := tview.NewList().ShowSecondaryText(false).
		SetChangedFunc(func(idx int, mainText, secondaryText string, shortcut rune) {
			a.showId = a.shows[idx].ID
			go func() {
				a.tv.QueueUpdateDraw(func() {
					a.getCatchupList(a.stationSlug, a.shows[idx].ID)
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

	cuList := tview.NewList().ShowSecondaryText(false)
	cuList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if cuList.GetItemCount() == 0 {
			return nil
		}

		r := event.Rune()
		k := event.Key()
		switch {
		case r == 'j', k == tcell.KeyDown:
			cuList.SetCurrentItem((cuList.GetCurrentItem() + 1) % cuList.GetItemCount())
		case r == 'k', k == tcell.KeyUp:
			cuList.SetCurrentItem((cuList.GetCurrentItem() - 1) % cuList.GetItemCount())
		case r == 'd':
			mainText, _ := cuList.GetItemText(cuList.GetCurrentItem())
			curr := cuList.GetCurrentItem()
			cuList.SetItemText(cuList.GetCurrentItem(), mainText+" [blue](downloading...)", "")
			go func() {
				if err := resty.DownloadFile(a.hc, fmt.Sprintf("%s.m4a", mainText), a.catchups[curr].StreamURL); err != nil {
					log.Fatal(err)
				}
				a.tv.QueueUpdateDraw(func() {
					cuList.SetItemText(curr, mainText, "")
				})
			}()
		case k == tcell.KeyEnter:
			ep := a.catchups[cuList.GetCurrentItem()]
			a.stream(a.catchupList, lo.Map(a.catchups, func(item models.Episode, _ int) streamItem {
				mainText := fmt.Sprintf("%s - %s - %s", item.Title, item.StartDate.Format("Mon 2006-01-02"), item.Availability)
				return streamItem{Name: mainText, Id: item.ID}
			}), ep.StreamURL)
		}
		return nil
	})
	a.catchupList = cuList

	// create keyboard shortcut help
	a.kbdShortcuts[1] = map[string]string{"\u21B5": "play/pause"}
	a.kbdShortcuts[2] = map[string]string{}
	a.kbdShortcuts[3] = map[string]string{"\u21B5": "play/pause", "d": "download"}
	a.helpText.SetText(text.FormatHelp(a.kbdShortcuts[1]))
}

// initialise the TUI app
func (a *app) initTui() {
	a.initViews()

	a.render()

	// base input handlers
	a.tv.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			// Exit the application
			a.tv.Stop()
			return nil
		}
		switch event.Rune() {
		case 49: // 1
			a.tv.SetFocus(a.stationsList)
			a.helpText.SetText(text.FormatHelp(a.kbdShortcuts[1]))
			return nil
		case 50: // 2
			a.tv.SetFocus(a.showsList)
			a.helpText.SetText(text.FormatHelp(a.kbdShortcuts[2]))
			return nil
		case 51: // 2
			a.tv.SetFocus(a.catchupList)
			a.helpText.SetText(text.FormatHelp(a.kbdShortcuts[3]))
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
	epsFlex.AddItem(a.catchupList, 0, 1, true)

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
	a.stations = make([]models.StationBrand, 0)

	stations, err := a.gp.GetStations()
	if err != nil {
		log.Fatal(err)
	}

	for _, station := range stations.PageProps.Feature.Blocks[0].Brands {
		a.stations = append(a.stations, station)
		text := station.Name
		a.stationsList.AddItem(text, "", 0, nil)
	}
}

// update the list of shows for a station
func (a *app) getShowsList(slug string) {
	a.showsList.Clear()
	a.shows = make([]models.CatchupInfo, 0)

	if slug != "" {
		shows, err := a.gp.GetCatchup(slug)
		if err != nil {
			log.Fatal(err)
		}

		for _, show := range shows.PageProps.CatchupInfo {
			a.shows = append(a.shows, show)
			a.showsList.AddItem(show.Title, "", 0, nil)
		}
	}
}

// update the catchup list of the given station and show
func (a *app) getCatchupList(slug, id string) {
	a.catchupList.Clear()
	a.catchups = make([]models.Episode, 0)

	if slug != "" {
		shows, err := a.gp.GetCatchupShows(slug, id)
		if err != nil {
			log.Fatal(err)
		}

		for _, show := range shows.PageProps.CatchupInfo.Episodes {
			a.catchups = append(a.catchups, show)
			text := fmt.Sprintf("%s - %s - %s", show.Title, show.StartDate.Format("Mon 2006-01-02"), show.Availability)
			a.catchupList.AddItem(text, "", 0, nil)
		}
	}
}

// pre-fetch all data so it's in the cache
func (a *app) prefetch() {
	stations, err := a.gp.GetStations()
	if err != nil {
		return
	}
	brands := stations.PageProps.Feature.Blocks[0].Brands
	if len(brands) == 0 {
		return
	}
	for _, st := range brands {
		cu, _ := a.gp.GetCatchup(st.Slug)
		if len(cu.PageProps.CatchupInfo) > 0 {
			for _, cu := range cu.PageProps.CatchupInfo {
				go func() {
					_, _ = a.gp.GetCatchupShows(st.Slug, cu.ID)
				}()
			}
		}
	}
}
