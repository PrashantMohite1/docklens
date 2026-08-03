package analyzer

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

func errorCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Run_in_container(imageName string, imgpath string) {
	// Create a new client that handles common environment variables
	// for configuration (DOCKER_HOST, DOCKER_API_VERSION), and does
	// API-version negotiation to allow downgrading the API version
	// when connecting with an older daemon version.
	apiClient, err := client.New(client.FromEnv)
	errorCheck(err)

	cmd := []string{"/bin/sh", "-c", "sha256sum " + imgpath}

	cfg := &container.Config{
		Image:        imageName,
		Cmd:          cmd,
		OpenStdin:    true,
		AttachStdin:  false,
		AttachStdout: false,
		AttachStderr: false,
		Tty:          true,
	}
	creatopt := client.ContainerCreateOptions{
		Config: cfg,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	createCont, err := apiClient.ContainerCreate(ctx, creatopt)

	if err != nil {
		fmt.Println("Error creating container: ", err)
	}

	// start the container
	_, err = apiClient.ContainerStart(ctx, createCont.ID, client.ContainerStartOptions{})
	if err != nil {
		fmt.Println("Error starting container: ", err)
	}

	// Wait for the command to finish
	waitCh := apiClient.ContainerWait(ctx, createCont.ID, client.ContainerWaitOptions{
		Condition: container.WaitConditionNotRunning,
	})
	select {
	case err := <-waitCh.Error:
		if err != nil {
			log.Fatal(err)
		}
	case <-waitCh.Result:
	}

	log, err := apiClient.ContainerLogs(ctx, createCont.ID, client.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})

	errorCheck(err)

	out1, _ := io.ReadAll(log)
	fmt.Print("out1 : ", string(out1))

	_, err = apiClient.ContainerRemove(ctx, createCont.ID, client.ContainerRemoveOptions{})
	errorCheck(err)
}
