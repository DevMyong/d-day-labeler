package main

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/go-github/v45/github"
)

func main() {
	ctx := context.Background()
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN env not set")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// GITHUB_REPOSITORY format: "owner/repo"
	repoFull := os.Getenv("GITHUB_REPOSITORY")
	if repoFull == "" {
		log.Fatal("GITHUB_REPOSITORY env not set")
	}
	parts := strings.Split(repoFull, "/")
	if len(parts) != 2 {
		log.Fatal("GITHUB_REPOSITORY env has invalid format")
	}
	owner := parts[0]
	repo := parts[1]

	// Get open PRs
	pulls, _, err := client.PullRequests.List(ctx, owner, repo, &github.PullRequestListOptions{
		State: "open",
	})
	if err != nil {
		log.Fatalf("Failed to list PRs: %v", err)
	}

	// Regex to match D-0, D-1, D-2, ...
	re := regexp.MustCompile(`^D-(\d)$`)
	for _, pr := range pulls {
		prNumber := pr.GetNumber()
		labels, _, err := client.Issues.ListLabelsByIssue(ctx, owner, repo, prNumber, nil)
		if err != nil {
			log.Printf("Failed to get labels for PR #%d: %v", prNumber, err)
			continue
		}
		for _, label := range labels {
			name := label.GetName()
			match := re.FindStringSubmatch(name)
			if match != nil {
				day, err := strconv.Atoi(match[1])
				if err != nil {
					log.Printf("Label %s failed to convert to number: %v", name, err)
					continue
				}
				if day == 0 {
					log.Printf("PR #%d's label %s is already D-0.", prNumber, name)
					continue
				}
				newLabel := fmt.Sprintf("D-%d", day-1)
				// delete old label
				_, err = client.Issues.RemoveLabelForIssue(ctx, owner, repo, prNumber, name)
				if err != nil {
					log.Printf("Failed to remove label %s from PR #%d: %v", name, prNumber, err)
					continue
				}
				// add new label
				_, _, err = client.Issues.AddLabelsToIssue(ctx, owner, repo, prNumber, []string{newLabel})
				if err != nil {
					log.Printf("Failed to add label %s to PR #%d: %v", newLabel, prNumber, err)
					continue
				}
				log.Printf("PR #%d label %s -> %s updated.", prNumber, name, newLabel)
			}
		}
	}
}
