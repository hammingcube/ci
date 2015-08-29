package main

import (
	"fmt"
	"github.com/google/go-github/github"
	"github.com/phayes/hookserve/hookserve"
	"golang.org/x/oauth2"
	"log"
	"os"
	"time"
)

func sPtr(s string) *string { return &s }

func wait() {
	ticker := time.NewTicker(time.Second)
	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
		}
	}()
	time.Sleep(time.Second * 180)
	ticker.Stop()
	fmt.Println("Ticker stopped")
}

func build(commit *hookserve.Event) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GH_TOKEN")},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	repoStatus, _, err := client.Repositories.CreateStatus(commit.Owner, commit.Repo, commit.Commit,
		&github.RepoStatus{
			State:       sPtr("pending"),
			TargetURL:   sPtr("https://www.google.com"),
			Description: sPtr("The build started"),
			Context:     sPtr("ci/builds"),
		})

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(repoStatus)
	wait()
	repoStatus, _, err = client.Repositories.CreateStatus(commit.Owner, commit.Repo, commit.Commit,
		&github.RepoStatus{
			State:       sPtr("success"),
			TargetURL:   sPtr("https://www.google.com"),
			Description: sPtr("The build succeeded"),
			Context:     sPtr("ci/builds"),
		})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(repoStatus)
}

func main() {
	log.Println(os.Getenv("GH_TOKEN"))
	server := hookserve.NewServer()
	server.Port = 8120
	server.Secret = "absolutesecret"
	server.GoListenAndServe()
	fmt.Printf("Listening on %d\n", server.Port)
	for {
		select {
		case commit := <-server.Events:
			fmt.Println(commit.Owner, commit.Repo, commit.Branch, commit.Commit)
			build(&commit)
		default:
			time.Sleep(100)
			//fmt.Println("No activity...")
		}
	}
}
