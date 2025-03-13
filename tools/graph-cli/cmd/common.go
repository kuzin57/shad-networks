package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/kuzin57/shad-networks/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getClient() (generated.GraphClient, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("error while creating grpc client")
	}

	return generated.NewGraphClient(conn), nil
}

func scrollAndSaveImage(
	ctx context.Context,
	client generated.GraphClient,
	imageFile,
	scrollID string,
	b64Image []byte,
) error {
	for {
		scrollResponse, err := client.Scroll(ctx, &generated.ScrollRequest{
			ScrollId: scrollID,
		})
		if err != nil {
			return fmt.Errorf("scroll error: %w", err)
		}

		b64Image = append(b64Image, scrollResponse.B64Image...)

		if scrollResponse.IsOver {
			break
		}
	}

	image, err := base64.StdEncoding.DecodeString(string(b64Image))
	if err != nil {
		return fmt.Errorf("decode from b64: %w", err)
	}

	if err := os.WriteFile(imageFile, image, 0666); err != nil {
		return fmt.Errorf("write to file: %w", err)
	}

	return nil
}
