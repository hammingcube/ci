package main

import (
	"encoding/json"
	"fmt"
	"github.com/CloudCom/firego"
	"github.com/google/go-github/github"
	"github.com/maddyonline/ci/prepare"
	"github.com/phayes/hookserve/hookserve"
	"github.com/streamrail/concurrent-map"
	"golang.org/x/oauth2"
	"log"
	"os"
	"strings"
	"time"
)

const (
	PENDING = "pending"
	SUCCESS = "success"
	ERROR   = "error"
	FAILURE = "failure"
)

var mymap cmap.ConcurrentMap

func sPtr(s string) *string { return &s }

func wait() {
	ticker := time.NewTicker(time.Second)
	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
		}
	}()
	time.Sleep(time.Second * 15)
	ticker.Stop()
	fmt.Println("Ticker stopped")
}

func build(commit *hookserve.Event) {
	key := strings.Join([]string{commit.Owner, commit.Repo, commit.Commit}, ",")
	if _, ok := mymap.Get(key); ok {
		return
	} else {
		mymap.Set(key, "pending")
	}
	defer mymap.Remove(key)
	fmt.Println("Building: ", commit.Owner, commit.Repo, commit.Branch, commit.Commit)

	currentStatus := PENDING

	buildURL := fmt.Sprintf("https://builds.firebaseio.com/%s/%s/%s", commit.Owner, commit.Repo, commit.Branch)
	f := firego.New(buildURL)
	f.Auth(os.Getenv("FIREBASE_SECRET"))
	v := map[string]string{
		"commit": commit.Commit,
		"status": currentStatus,
	}
	pushedFirego, err := f.Push(v)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", pushedFirego)

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GH_TOKEN")},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	repoStatus, _, err := client.Repositories.CreateStatus(commit.Owner, commit.Repo, commit.Commit,
		&github.RepoStatus{
			State:       sPtr(currentStatus),
			TargetURL:   sPtr(buildURL),
			Description: sPtr("The build started"),
			Context:     sPtr("ci/builds"),
		})

	if err != nil {
		fmt.Println(err)
	}
	jsonBytes := prepare.Main(commit.Owner, commit.Repo)
	var statusVal map[string]string
	json.Unmarshal(jsonBytes, &statusVal)
	if statusVal != nil {
		fmt.Println("Got the following status: ", statusVal["status"])
		currentStatus = statusVal["status"]
	}
	wait()
	pushedFirego.Update(map[string]string{"status": currentStatus})
	repoStatus, _, err = client.Repositories.CreateStatus(commit.Owner, commit.Repo, commit.Commit,
		&github.RepoStatus{
			State:       sPtr(currentStatus),
			TargetURL:   sPtr(buildURL),
			Description: sPtr("The build succeeded"),
			Context:     sPtr("ci/builds"),
		})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(repoStatus)
}

func main() {
	// Create a new map.
	mymap = cmap.New()

	log.Println(os.Getenv("GH_TOKEN"))
	server := hookserve.NewServer()
	server.Port = 8120
	server.Secret = "absolutesecret"
	server.GoListenAndServe()
	fmt.Printf("Listening on %d\n", server.Port)
	for {
		select {
		case commit := <-server.Events:
			fmt.Println("Got: ", commit.Owner, commit.Repo, commit.Branch, commit.Commit)
			go func() { build(&commit) }()
		default:
			time.Sleep(100)
			//fmt.Println("No activity...")
		}
	}
}
