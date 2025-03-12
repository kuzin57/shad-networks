package main

import (
	"context"
	"encoding/base64"
	"os"

	"github.com/kuzin57/shad-networks/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var (
		ctx     = context.Background()
		graphID = "7121d59c-0fe0-444e-872f-557ae748b1f6"
	)

	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	client := generated.NewGraphClient(conn)

	graphResponse, err := client.Get(ctx, &generated.GetGraphRequest{
		GraphId: graphID,
	})
	if err != nil {
		panic(err)
	}

	var (
		scrollID = graphResponse.ScrollId
		b64Image = graphResponse.B64Image
	)

	for {
		scrollResponse, err := client.Scroll(ctx, &generated.ScrollRequest{
			ScrollId: scrollID,
		})
		if err != nil {
			panic(err)
		}

		b64Image = append(b64Image, scrollResponse.B64Image...)

		if scrollResponse.IsOver {
			break
		}
	}

	image, err := base64.StdEncoding.DecodeString(string(b64Image))
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile("graph.svg", image, 0666); err != nil {
		panic(err)
	}
}
