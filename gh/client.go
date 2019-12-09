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
func GetPRDetailsWithContext(ctx context.Context, cfg *Config, ids ...string) (global error) {
	owner, repo, err := cfg.GetOwner()
	if err != nil {
		return err
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.Token},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	fmt.Printf("%+v", client)

	// pull requests ids
	for _, id := range ids {
		// owner , repo , id
		number, err := strconv.Atoi(id)
		if err != nil {

			if global != nil {
				global = fmt.Errorf("%v,%w", global, err)
			} else {
				global = err
			}

			continue
		}

		// worst case if we can get the info
		pr, response, err := client.PullRequests.Get(ctx, owner, repo, number)
		if err != nil {
			if global != nil {
				global = fmt.Errorf("%v, %w", global, err)
			} else {
				global = err
			}
		}

		log.Printf("response: %+v\n", response)
		log.Printf("%+v\n", pr)

		// now get all files (statistics)

	}

	return
}
