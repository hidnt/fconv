package service

import (
	"context"
	"errors"
	"fmt"
)

var NoAdapterFoundErr error = errors.New("no adapter found for conversion")

type Converter interface {
	Convert(ctx context.Context, srcPath, dstPath, dstExt string) error
	Supports(ctx context.Context, srcExt, dstExt string) bool
}

type ConverterService interface {
	Convert(ctx context.Context, srcPath, dstPath, srcExt, dstExt string) error
}

func NewConverterService(adapters []Converter) ConverterService {
	return &converterSrv{
		adapters: adapters,
	}
}

type converterSrv struct {
	adapters []Converter
}

func (s converterSrv) Convert(ctx context.Context, srcPath, dstPath, srcExt, dstExt string) error {
	if srcExt == dstExt {
		return nil
	}
	for _, a := range s.adapters {
		if a.Supports(ctx, srcExt, dstExt) {
			return a.Convert(ctx, srcPath, dstPath, dstExt)
		}
	}
	return fmt.Errorf("%w %s to %s", NoAdapterFoundErr, srcExt, dstExt)
}
