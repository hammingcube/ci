package main

import (
	"encoding/json"
	_ "fmt"
	"github.com/google/go-github/github"
	_ "github.com/phayes/hookserve/hookserve"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"
	_ "time"
)

func sPtr(s string) *string { return &s }

func downloadArchive(url *url.URL, dest string) {
	log.Printf("Downloading: %s\n", url)
	cmd := exec.Command("curl", "-o", dest, "-O", url.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
func extractArchive(zipFile, outputDir string) {
	log.Printf("Extracting...\n")
	cmd := exec.Command("unzip", zipFile, "-d", outputDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func doIt(client *github.Client, owner, repo string, opt *github.RepositoryContentGetOptions) {
	url, _, err := client.Repositories.GetArchiveLink(owner, repo, github.Zipball, opt)
	if err != nil {
		log.Fatal(err)
	}
	downloadArchive(url, "abc.zip")
	os.Mkdir("unique_dir", 0777)
	extractArchive("abc.zip", "unique_dir")
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
	client := github.NewClient(nil)
	opt := &github.RepositoryContentGetOptions{"master"}
	doIt(client, "maddyonline", "fun-with-algo", opt)
	doIt(client, "maddyonline", "epibook.github.io", opt)
}
