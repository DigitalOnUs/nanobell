package gh

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
)

//Config basic params
type Config struct {
	Token string
	Repo  string
}

//FeatureEntry feature vector
type FeatureEntry struct {
	Title           string  `json:"title"`
	Comments        int     `json:"comments"` // will be 0
	OpenDays        int     `json:"open-days"`
	Adds            int     `json:"adds"`
	Dels            int     `json:"dels"`
	Files           int     `json:"files"`
	Number          int     `json:"number"`
	CouplingAverage float64 `json:"coupling average"`
}

var (
	//ErrInvalidRepoName Awful case where we got nothing
	ErrInvalidRepoName = errors.New("Invalid repo name ")
)

//GetOwner error for stuff
func (cfg *Config) GetOwner() (owner, repo string, err error) {

	args := strings.Split(cfg.Repo, "/")
	if len(args) < 2 {
		return "", "", fmt.Errorf(" %w: %s", ErrInvalidRepoName, cfg.Repo)
	}

	return args[0], args[len(args)-1], nil
}

//GetPRDetailsWithContext ctx, token , .... ids pr
func GetPRDetailsWithContext(ctx context.Context, cfg *Config, errors chan error, ids ...string) <-chan FeatureEntry {
	owner, repo, err := cfg.GetOwner()
	queue := make(chan FeatureEntry)

	if err != nil {
		errors <- err
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.Token},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	go func(c chan FeatureEntry) {
		defer close(c)
		defer close(errors)
		// pull requests ids
		for _, id := range ids {
			// owner , repo , id
			number, err := strconv.Atoi(id)
			if err != nil {
				errors <- err
				continue
			}

			// worst case if we can get the info
			pr, response, err := client.PullRequests.Get(ctx, owner, repo, number)
			if err != nil {
				errors <- err
				continue
			}

			log.Printf("%+v\n", response)

			entry := FeatureEntry{
				Title:    pr.GetTitle(),
				Comments: pr.GetComments(),
			}

			/* TO DO GET FILE LIST FROM BELLATOR*/
			c <- entry
		}
	}(queue)

	return queue
}
