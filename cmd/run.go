package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vrecan/death/v3"
	"github.com/willgorman/mastodon-bsky/pkg/mastodon"
	"github.com/willgorman/mastodon-bsky/pkg/sync"
)

var runCmd = &cobra.Command{
	Use: "run",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("RUN")
		dataPath := viper.GetString("data_path")
		if dataPath == "" {
			return errors.New("missing data_path")
		}
		log.Println(dataPath)
		data, err := sync.CreateDatastore(dataPath)
		if err != nil {
			return fmt.Errorf("opening database %s: %w", dataPath, err)
		}

		// TODO: (willgorman) create mastodon/bsky source/sink
		process := sync.New(data, mastodon.NewFakeSource(), nil)

		// TODO: (willgorman) error logging
		go process.Run(context.TODO())
		d := death.NewDeath(syscall.SIGINT, syscall.SIGTERM)
		return d.WaitForDeath()
	},
}
