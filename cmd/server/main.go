package main

import (
	"fmt"
	"log"
	"os"

	"github.com/matthieukhl/align-back/config"
	"github.com/spf13/cobra"
)

var cfgFile string

func main() {
	cmd := &cobra.Command{
		Use:   "align-back",
		Short: "Pilates Management System API Server",
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}

	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config/config.yaml)")
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() {
	// Intialize configuration
	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load configuration")
	}

	// Initialize logger
	logger.InitLogger(cfg.LogLevel)
}
