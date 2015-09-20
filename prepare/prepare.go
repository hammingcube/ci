package main

import (
	"encoding/json"
	_ "fmt"
	"github.com/google/go-github/github"
	_ "github.com/phayes/hookserve/hookserve"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	_ "time"
)

func sPtr(s string) *string { return &s }

func download() {

	client := github.NewClient(nil)
	owner := "maddyonline"
	repo := "fun-with-algo"
	opt := &github.RepositoryContentGetOptions{"master"}

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
	dirs, err := ioutil.ReadDir("unique_dir")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Found following directories:\n%v\n", dirs[0].Name())
	filename := path.Join("unique_dir", dirs[0].Name(), "runtests.json")
	data, err := ioutil.ReadFile(filename)
	var v map[string]interface{}
	json.Unmarshal(data, &v)
	log.Printf("Read:%s\n", v)
}

func downloadFile(client *github.Client, owner, repo, filepath string, opt *github.RepositoryContentGetOptions) {
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
	download()
}
