package cmd

import (
	"strings"

	"github.com/hidnt/fconv/internal/fconv"
	"github.com/hidnt/fconv/internal/models"
	"github.com/hidnt/fconv/internal/service"
	"github.com/spf13/cobra"
)

func NewRootCmd(srv service.ConverterService) *cobra.Command {
	var cfg models.Config
	var rootCmd = &cobra.Command{
		Use:   "fconv [FILES...] --to <EXTENSION> [OPTIONS]",
		Short: "File type convert",
		Long:  ``,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			app := fconv.New(cfg, srv)
			app.Fconv(ctx, args)
		},
	}

	rootCmd.Flags().StringVarP(&cfg.DstExt, "to", "t", "", "Target extension (required)")
	rootCmd.MarkFlagRequired("to")
	cfg.DstExt = strings.TrimPrefix(strings.ToLower(cfg.DstExt), ".")
	rootCmd.Flags().StringVarP(&cfg.DstFolder, "output", "o", "", "Destination folder for saving")
	rootCmd.Flags().BoolVarP(&cfg.NeedRecursion, "recursive", "r", false, "Recursive directory traversal")
	rootCmd.Flags().IntVarP(&cfg.LevelOfRec, "level", "L", -1, "Level of recursion")
	if !cfg.NeedRecursion {
		cfg.LevelOfRec = 1
	}
	rootCmd.Flags().BoolVarP(&cfg.Delete, "delete", "d", false, "Delete files after convertion")
	rootCmd.Flags().BoolVarP(&cfg.Force, "force", "f", false, "Overwrite the target file if it already exists")

	return rootCmd
}
