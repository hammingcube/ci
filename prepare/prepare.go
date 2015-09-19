package main

import (
	_ "fmt"
	"github.com/google/go-github/github"
	_ "github.com/phayes/hookserve/hookserve"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"os"
	_ "time"
)

func sPtr(s string) *string { return &s }

func download() {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GH_TOKEN")},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	owner := "maddyonline"
	repo := "fun-with-algo"
	filepath := "testit/v1/Can_string_be_palindrome_gen.cc"
	opt := &github.RepositoryContentGetOptions{"b04c4c9ba6d9548df666913e6b2f6a164ad03cfe"}
	r, err := client.Repositories.DownloadContents(owner, repo, filepath, opt)
	if err != nil {
		log.Fatal(err)
	}
	d1, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("dl_data", d1, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.Println(os.Getenv("GH_TOKEN"))
	download()
}
