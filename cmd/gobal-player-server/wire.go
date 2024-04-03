//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/jj-style/gobal-player/cmd/gobal-player-server/internal/biz"
	"github.com/jj-style/gobal-player/cmd/gobal-player-server/internal/config"
	"github.com/jj-style/gobal-player/cmd/gobal-player-server/internal/server"
	"github.com/jj-style/gobal-player/cmd/gobal-player-server/internal/service"
)

func InitializeServer(config *config.Config) (*server.Server, func(), error) {
	panic(wire.Build(server.GlobalPlayerProvider, service.NewService, biz.ProviderSet, server.NewServer))
}
