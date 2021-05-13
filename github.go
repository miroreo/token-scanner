package main

import (
	"github.com/google/go-github/github"
	"github.com/bradleyfalzon/ghinstallation"
	"net/http"
	"os"
	"context"
	// "strconv"
	"fmt"
	"time"
)

func githubInit() *github.Client {
	// Wrap the shared transport for use with the integration ID 1 authenticating with installation ID 99.
	itr, err := ghinstallation.New(http.DefaultTransport, 115168, 16900202, []byte(os.Getenv("GITHUB_PRIVATE_KEY")))
	if err != nil {
		handleError(err)
	}
	fmt.Println(itr.BaseURL)
	client := github.NewClient(&http.Client{Transport: itr})
	return client
}

func pushToken(client *github.Client, message string, token string) error {
	committerName := "Token Scanner Bot"
	committerEmail := "tokenscanner.roblockhead@example.com"
	commitMessage := "Push token from message " + message
	contResp, _, err := client.Repositories.CreateFile(context.Background(), "RoBlockHead", "TokenDisabler", "tokens/" + message + ".txt", &github.RepositoryContentFileOptions{
		Content: []byte(fmt.Sprintf("Hey there, this file was created automatically because a bot found your token in public. By creating this file, the Discord bots will immediately reset the token.\nHopefully, we caught it in time to make sure that you weren't compromised.\nThis is a dangerous thing to leak, as your bot can be controlled by anyone if they have the token! Please keep your token safe. \nReporter: %s\nTime Detected: %s\n\n%s", "Token Scanner Bot (replit/@RoBlockHead)", (time.Now()).String, token)),
		Message: &commitMessage,
		Author: &github.CommitAuthor{
			Name: &committerName,
			Email: &committerEmail,
		},
		Committer: &github.CommitAuthor{
			Name: &committerName,
			Email: &committerEmail,
		},
	})
	if err != nil{
		fmt.Println(err)
		return err
	} else {
		fmt.Printf("Token Pushed: %v\n", contResp.Commit.SHA)
		return nil
	}
}