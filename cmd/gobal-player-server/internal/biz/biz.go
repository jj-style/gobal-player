package biz

import (
	"github.com/google/wire"
	"github.com/jj-style/gobal-player/cmd/gobal-player-server/internal/biz/globalplayer"
)

var ProviderSet = wire.NewSet(globalplayer.NewUseCase)
