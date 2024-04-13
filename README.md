# gobal-player

[![CI](https://github.com/jj-style/gobal-player/actions/workflows/ci.yml/badge.svg)](https://github.com/jj-style/gobal-player/actions/workflows/ci.yml)

Unofficial collection of apps, packages and APIs for global player radio.

## gobal-player-tui
TUI application for streaming live radio, and streaming or downloading catchup episodes.

[![Packaging status](https://repology.org/badge/vertical-allrepos/gobal-player-tui.svg)](https://repology.org/project/gobal-player-tui/versions)

### Install

#### Arch Linux
Published in the [AUR: gobal-player-tui](https://aur.archlinux.org/packages/gobal-player-tui).
`yay -S gobal-player-tui`

#### From source
```bash
go install github.com/jj-style/gobal-player/cmd/gobal-player-tui@v0.1.19
gobal-player-tui
```

![tui-demo](.github/assets/tui.gif)

## gobal-player-server
RESTful server with friendly APIs to global player, and more features coming soon!

### Features
- RESTful APIs for exploring stations, shows and episodes
- RSS feeds:
  - RSS feed for all shows in a station
  - RSS feed for all episodes for a given show


### Docker

Run with docker: `docker run --rm -it -p 8080:8080 ghcr.io/jj-style/gobal-player/gobal-player-server:v0.1.19`

## Developing

Please install [pre-commit](https://pre-commit.com/#install), and [just](https://github.com/casey/just?tab=readme-ov-file#installation) and run `just hooks` to initialize your git pre-commit hooks.
