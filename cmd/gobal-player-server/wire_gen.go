// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/jj-style/gobal-player/cmd/gobal-player-server/internal/biz/globalplayer"
	"github.com/jj-style/gobal-player/cmd/gobal-player-server/internal/config"
	"github.com/jj-style/gobal-player/cmd/gobal-player-server/internal/server"
	"github.com/jj-style/gobal-player/cmd/gobal-player-server/internal/service"
)

// Injectors from wire.go:

func InitializeServer(config2 *config.Config) (*server.Server, error) {
	cache := server.NewCache(config2)
	globalPlayer, err := server.NewGlobalPlayer(config2, cache)
	if err != nil {
		return nil, err
	}
	useCase := globalplayer.NewUseCase(globalPlayer)
	serviceService := service.NewService(useCase)
	serverServer := server.NewServer(serviceService)
	return serverServer, nil
}
