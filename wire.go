//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/hidnt/fconv/internal/adapters"
	"github.com/hidnt/fconv/internal/service"
)

func InitializeService() service.ConverterService {
	wire.Build(
		adapters.AdaptersSet,
		service.NewConverterService,
	)
	return nil
}
