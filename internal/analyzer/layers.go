package analyzer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/PrashantMohite1/docklens/internal/docker"
	"github.com/moby/moby/client"
)

func Get_img_layer(imageName string) {

	ctx := context.Background()
	apiClient, err := docker.CreateClient()

	if err != nil {
		log.Fatal(err)
	}
	defer apiClient.Close()

	var raw bytes.Buffer

	_, err = apiClient.ImageInspect(ctx, imageName, client.ImageInspectWithRawResponse(&raw))
	if err != nil {
		log.Fatal(err)
	}

	// raw is compact JSON from Docker. Indent it so it's human-readable.
	var pretty bytes.Buffer
	if err := json.Indent(&pretty, raw.Bytes(), "", "  "); err != nil {
		log.Fatal(err)
	}

	layers := pretty.String()

	fmt.Println(layers)

}
