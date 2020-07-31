package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

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

func main() {
	if len(os.Args) != 2 {
		usage()
		return
	}

	ctx := context.Background()

	client, err := authenticatedClient(ctx)
	if err != nil {
		panic(err)
	}

	org := os.Args[1]

	repos, _, err := client.Repositories.List(ctx, org, nil)
	if err != nil {
		panic(err)
	}

	rawLabels, err := ioutil.ReadFile("labels.json")
	if err != nil {
		panic(err)
	}

	labels := []github.Label{}
	err = json.Unmarshal(rawLabels, &labels)
	if err != nil {
		panic(err)
	}

next_repo:
	for _, repo := range repos {
		fmt.Printf("\nUpdating repository %s/%s\n", org, repo.GetName())

		for _, label := range labels {
			fmt.Printf("* %s: %s", label.GetName(), label.GetDescription())
			if _, resp, err := client.Issues.CreateLabel(ctx, org, repo.GetName(), &label); err != nil {
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
