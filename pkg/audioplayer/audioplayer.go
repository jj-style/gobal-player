package audioplayer

import (
	vlc "github.com/adrg/libvlc-go/v3"
)

// A Player can play audio
type Player interface {
	// Play plays the audio stream at the given URL.
	//
	// Returns:
	//
	// - channel to stop the audio stream
	//
	// - error if it could not play the audio
	Play(url string) (chan struct{}, error)
}

// vlcPlayer plays audio using vlc
type vlcPlayer struct {
	player *vlc.Player
}

func NewPlayer() (Player, func(), error) {
	if err := vlc.Init("--no-video", "--quiet"); err != nil {
		return nil, func() {}, err
	}
	player, err := vlc.NewPlayer()
	if err != nil {
		return nil, func() {}, err
	}

	return &vlcPlayer{player}, func() {
		player.Release()
		player.Stop()
		vlc.Release()
	}, nil
}

func (p *vlcPlayer) Play(url string) (chan struct{}, error) {
	quit := make(chan struct{})
	go func() {
		media, err := p.player.LoadMediaFromURL(url)
		if err != nil {
			panic(err)
		}
		defer media.Release()

		// Retrieve player event manager.
		manager, err := p.player.EventManager()
		if err != nil {
			panic(err)
		}

		eventCallback := func(event vlc.Event, userData interface{}) {
			close(quit)
		}

		eventID, err := manager.Attach(vlc.MediaPlayerEndReached, eventCallback, nil)
		if err != nil {
			panic(err)
		}
		defer manager.Detach(eventID)

		// Start playing the media.
		err = p.player.Play()
		if err != nil {
			panic(err)
		}
		<-quit
		p.player.Stop()
	}()
	return quit, nil
}
