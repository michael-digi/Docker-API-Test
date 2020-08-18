package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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

func checkAPIKey(next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		apiKey := req.Header.Get("X-Api-Key")

		if apiKey == "" {
			return
		}
		next(res, req)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		log.Println("Method:", req.Method)
		log.Println("Route:", req.RequestURI)
		log.Println("Body:", req.Body)
		log.Println("Host:", req.Host)
		log.Println("Remote Address:", req.RemoteAddr)

		next.ServeHTTP(res, req)
	})
}

func test(res http.ResponseWriter, req *http.Request) {
	fmt.Println("You hit test")
}

func showContainers(res http.ResponseWriter, req *http.Request) {
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

	json.NewEncoder(res).Encode(response)
}

func main() {
	router := mux.NewRouter()

	router.Use(loggingMiddleware)

	router.HandleFunc("/containers", checkAPIKey(showContainers)).Methods("GET")
	//testing to see if loggingMiddleware working for all routes
	router.HandleFunc("/test", test).Methods("GET")

	http.ListenAndServe(":3000", router)
}
