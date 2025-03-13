/*
Copyright © 2025 NAME HERE rkuzin.2003@gmail.com
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/kuzin57/shad-networks/generated"
	"github.com/spf13/cobra"
)

var getFlags struct {
	graphID string
	file    string
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get graph",
	Long:  `get graph to SVG image. Name of file to save image to is passed to arguments`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return getGraph(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.PersistentFlags().StringVarP(&getFlags.graphID, "graph", "g", "", "graph ID")
	getCmd.PersistentFlags().StringVarP(&getFlags.file, "file", "f", "", "file to save an image")
}

func getGraph(ctx context.Context) error {
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("get client: %w", err)
	}

	graphResponse, err := client.Get(ctx, &generated.GetGraphRequest{
		GraphId: getFlags.graphID,
	})
	if err != nil {
		return fmt.Errorf("get graph error: %w", err)
	}

	return scrollAndSaveImage(ctx, client, getFlags.file, graphResponse.ScrollId, graphResponse.B64Image)
}
