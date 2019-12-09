package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/DigitalOnUs/nanobell/gh"
)

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

	cfg := gh.Config{envs[TOKEN], envs[REPO]}

	if err := gh.GetPRDetailsWithContext(ctx, &cfg, envs[PULL]); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
