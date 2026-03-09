package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/hidnt/fconv/cmd"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	srv := InitializeService()
	cmd.NewRootCmd(srv).ExecuteContext(ctx)
}
