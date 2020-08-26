package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/ghodss/yaml"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

func authenticatedClient(ctx context.Context) (*github.Client, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN is not set")
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return client, nil
}

func usage() {
	fmt.Printf("Usage: %s ORGANIZATION\n", os.Args[0])
}

const (
	configFilename = "labels.json"
)

type Config struct {
	Repositories []string       `json:"repos"`
	Labels       []github.Label `json:"labels"`
}

func parseConfig(file string, config *Config) error {
	rawLabels, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(rawLabels, config)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		usage()
		return
	}

	config := Config{}
	err := parseConfig(configFilename, &config)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	client, err := authenticatedClient(ctx)
	if err != nil {
		panic(err)
	}

	org := os.Args[1]

	repoNames := config.Repositories

	if config.Repositories == nil {
		repos, _, err := client.Repositories.List(ctx, org, nil)
		if err != nil {
			panic(err)
		}
		for _, r := range repos {
			repoNames = append(repoNames, r.GetName())
		}
	}

next_repo:
	for _, repo := range repoNames {
		fmt.Printf("\nUpdating repository %s/%s\n", org, repo)

		for _, label := range config.Labels {
			fmt.Printf("* %s: %s", label.GetName(), label.GetDescription())
			if _, resp, err := client.Issues.CreateLabel(ctx, org, repo, &label); err != nil {
				if resp.StatusCode == http.StatusForbidden {
					fmt.Printf(" (403, skipping repo)\n")
					continue next_repo
				}

				if resp.StatusCode == http.StatusUnprocessableEntity {
					if err, ok := err.(*github.ErrorResponse); ok {
						if len(err.Errors) == 1 && err.Errors[0].Code == "already_exists" {
							fmt.Printf(" (already exists)\n")
							continue
						}
					}
				}
				panic(err)
			}
			fmt.Printf(" (created)\n")
		}

	}
}
