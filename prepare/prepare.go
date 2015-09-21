package prepare

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	_ "github.com/phayes/hookserve/hookserve"
	"io/ioutil"
	"log"
	"math/rand"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
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

func doIt(client *github.Client, owner, repo string, opt *github.RepositoryContentGetOptions) (string, string, string) {
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
	var runTestsConfig map[string]string
	var problem, mySolnDir string
	json.Unmarshal(data, &runTestsConfig)
	log.Printf("Read:%s\n", runTestsConfig)
	arr := strings.Split(runTestsConfig["runtests"], ",")
	if len(arr) > 1 {
		problem = arr[0]
		mySolnDir = arr[1]
	}
	fmt.Printf("problem: %s, mysolnDir: %s\n", problem, mySolnDir)
	return path.Join("unique_dir", dirs[0].Name()), problem, mySolnDir
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

func Main() []byte {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randDir := fmt.Sprintf("rand_dir_%d", r.Uint32())
	fmt.Printf("Using %s as working directory.\n", randDir)
	os.Mkdir(randDir, 0777)
	err := os.Chdir(randDir)
	if err != nil {
		log.Fatal(err)
	}
	client := github.NewClient(nil)
	opt := &github.RepositoryContentGetOptions{"master"}
	mySolnRepo, problem, mySolnDir := doIt(client, "maddyonline", "fun-with-algo", opt)
	problemsRepo, _, _ := doIt(client, "maddyonline", "epibook.github.io", opt)
	fmt.Println("In main:")
	fmt.Println(problemsRepo, problem, mySolnRepo, mySolnDir)
	return createExe(problemsRepo, problem, mySolnRepo, mySolnDir)
}
