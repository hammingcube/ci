package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

const SCRIPT = "local_build.sh"
const TEMPL = "docker run --rm -v {{.Path}}:/app -v {{.Destination}}:/dest -w /app {{.Container}} sh {{.Script}} /dest/{{.Output}}"

type Config struct {
	Path        string
	Container   string
	Script      string
	Destination string
	Output      string
}

func cwd() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(dir)
	return dir
}

func dockerCmd(scriptPath, outFile string) string {
	destDir := cwd()
	fullPath, err := filepath.Abs(scriptPath)
	tmpl, err := template.New("test").Parse(TEMPL)
	if err != nil {
		panic(err)
	}
	containersMap := map[string]string{
		"cpp":    "glot/clang",
		"golang": "glot/golang",
	}

	scriptSrc := path.Join(scriptPath, SCRIPT)
	script, err := ioutil.ReadFile(scriptSrc)
	lines := strings.Split(string(script), "\n")
	var lang string
	//fmt.Println(lines[0])
	fmt.Sscanf(lines[0], "# Language: %s", &lang)
	//fmt.Println(lang)

	config := &Config{
		Path:        fullPath,
		Container:   containersMap[lang],
		Script:      SCRIPT,
		Destination: destDir,
		Output:      outFile,
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, config)
	if err != nil {
		panic(err)
	}
	return b.String()
}

func primarySoln(baseDir, problem string) string {
	s := "%s/solutions/annotated/%s/primary-solution"
	return fmt.Sprintf(s, baseDir, problem)
}

func primaryGen(baseDir, problem string) string {
	s := "%s/solutions/annotated/%s/primary-generator"
	return fmt.Sprintf(s, baseDir, problem)
}

func primaryRunner(baseDir, problem string) string {
	s := "%s/solutions/annotated/primary-runner"
	return fmt.Sprintf(s, baseDir)
}

func main() {
	fmt.Println(dockerCmd(primarySoln(os.Args[1], os.Args[2]), "primary-soln"))
	fmt.Println(dockerCmd(primaryGen(os.Args[1], os.Args[2]), "gen"))
	fmt.Println(dockerCmd(primaryRunner(os.Args[1], os.Args[2]), "runtest"))
	//dockerCmd(os.Args[1])
}
