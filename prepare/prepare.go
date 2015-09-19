package main

import (
	_ "fmt"
	"github.com/google/go-github/github"
	_ "github.com/phayes/hookserve/hookserve"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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
	filepath := "testit/v1"
	opt := &github.RepositoryContentGetOptions{"b04c4c9ba6d9548df666913e6b2f6a164ad03cfe"}

	url, _, err := client.Repositories.GetArchiveLink(owner, repo, github.Zipball, opt)
	if err != nil {
		log.Fatal(err)
	}

	downloadArchive := func() {
		log.Printf("Downloading: %s\n", url)
		cmd := exec.Command("curl", "-o", "abc.zip", "-O", url.String())
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
	extractArchive := func() {
		log.Printf("Extracting")
		cmd := exec.Command("unzip", "abc.zip", "-d", "unique_dir/")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}

	downloadArchive()
	os.Mkdir("unique_dir", 0777)
	extractArchive()

	name, err := exec.Command("ls", "unique_dir").CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Found the name: %s\n", name)

	dl := func() {
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
	dl()
}

func main() {
	log.Println(os.Getenv("GH_TOKEN"))
	download()
}
