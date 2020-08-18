package main

import (
	"context"
	"encoding/json"
	"net/http"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"github.com/gorilla/mux"
)

// Container holds info about docker container
type Container struct {
	ID      string `json:"Id"`
	Image   string
	ImageID string
	Command string
	Created int64
	State   string
	Status  string
}

func showContainers(w http.ResponseWriter, r *http.Request) {
	cli, err := docker.NewEnvClient()

	dockerJSON := []Container{}

	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})

	if err != nil {
		panic(err)
	}

	for _, c := range containers {
		dockerJSON = append(dockerJSON, Container{c.ID[:10], c.Image, c.ImageID, c.Command, c.Created, c.State, c.Status})
	}

	response := map[string][]Container{
		"data": dockerJSON,
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/containers", showContainers).Methods("GET")

	http.ListenAndServe(":3000", router)
}
