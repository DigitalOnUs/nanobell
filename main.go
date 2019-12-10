package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/DigitalOnUs/nanobell/gh"
)

// expected output
// title	comments	open-days	adds	dels	files	number	coupling average

// environment data
var (
	envs map[string]string
)

const (
	//TOKEN for the app to run requests in GITHUB
	TOKEN = "GITHUB_TOKEN"
	// ENV BY DRONE
	REPO = "DRONE_REPO"
	// PULL REQUEST ID that is done against a branch
	PULL = "DRONE_PULL_REQUEST"
)

func init() {
	// validations
	envs = make(map[string]string)
	reqs := []string{TOKEN, PULL, REPO}
	for _, param := range reqs {
		val := os.Getenv(param)
		if val == "" {
			log.Println("Missing param ", param)
			os.Exit(1)
		}
		envs[param] = val
	}
}

func main() {
	// fetching pull request
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		getOutOfHere := <-c
		log.Printf("going out with ... user cancel %s \n", getOutOfHere)
		cancel()
		os.Exit(1)
	}()

	cfg := gh.Config{Token: envs[TOKEN], Repo: envs[REPO]}

	errors := make(chan error)
	// ideally this is ready to support multiple requests in pipeline
	pool := gh.GetPRDetailsWithContext(ctx, &cfg, errors, envs[PULL])

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		for err := range errors {
			fmt.Println(err)
		}
		wg.Done()
	}()

	go func() {
		for entry := range pool {
			fmt.Println(entry)
		}
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("ya se acabo")

}
