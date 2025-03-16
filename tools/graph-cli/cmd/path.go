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

var pathFlags struct {
	graphID string
	file    string
	source  int
	target  int
	amount  int
}

var pathCmd = &cobra.Command{
	Use:   "path",
	Short: "find shortest path",
	Long: `find shortest path in a graph between <source> and <target> vertices.
			Returns an image of the path. Filename for image should be provided`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		return findPaths(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(pathCmd)

	pathCmd.PersistentFlags().StringVarP(&pathFlags.graphID, "graph", "g", "", "graph ID")
	pathCmd.PersistentFlags().StringVarP(&pathFlags.file, "file", "f", "", "file to save an image")
	pathCmd.PersistentFlags().IntVarP(&pathFlags.source, "source", "s", 0, "source vertex")
	pathCmd.PersistentFlags().IntVarP(&pathFlags.target, "target", "t", 0, "target vertex")
	pathCmd.PersistentFlags().IntVarP(&pathFlags.amount, "amount", "a", 1, "amount of paths")
}

func findPaths(ctx context.Context) error {
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("get client: %w", err)
	}

	graphResponse, err := client.FindPaths(ctx, &generated.FindPathsRequest{
		GraphId: pathFlags.graphID,
		From:    uint32(pathFlags.source),
		To:      uint32(pathFlags.target),
		Amount:  uint32(pathFlags.amount),
	})
	if err != nil {
		return fmt.Errorf("get graph error: %w", err)
	}

	return scrollAndSaveImage(ctx, client, pathFlags.file, graphResponse.ScrollId, graphResponse.B64Image)
}
