package cmd

import (
	"context"
	"os"
	"os/signal"

	"github.com/edwintcloud/go-stock-scanner/internal/config"
	"github.com/edwintcloud/go-stock-scanner/internal/scanner"
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
)

func ScanCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scan",
		Short: "Run stock market scanner",
		Run: func(cmd *cobra.Command, args []string) {
			runScan()
		},
	}
	return cmd
}

func runScan() {
	context, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	scanner, err := scanner.NewScanner(config)
	if err != nil {
		log.Fatal(err)
	}
	err = scanner.Start(context)
	if err != nil {
		log.Fatal(err)
	}

	// blocks until interrupt signal is received
	<-context.Done()
}
