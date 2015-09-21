package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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

func dockerCmd(scriptPath, outFile, destDir string) string {
	//fmt.Println(scriptPath)
	destDir = path.Join(cwd(), destDir)
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

func mySolution(baseDir, problem, mySoln string) string {
	s := "%s/my-solutions/%s/%s"
	return fmt.Sprintf(s, baseDir, problem, mySoln)
}

func createExe(problemsRepo, problem, mySolnRepo, mySolnDir string) {
	destDir := "work_dir"
	binDir, _ := filepath.Abs(destDir)
	runCmd := "docker run --rm -v %s:/app -w /app ubuntu ./runtest ./gen ./my-soln ./primary-soln"
	os.Mkdir(destDir, 0777)
	c1 := dockerCmd(primarySoln(problemsRepo, problem), "primary-soln", destDir)
	c2 := dockerCmd(primaryGen(problemsRepo, problem), "gen", destDir)
	c3 := dockerCmd(primaryRunner(problemsRepo, problem), "runtest", destDir)
	c4 := dockerCmd(mySolution(mySolnRepo, problem, mySolnDir), "my-soln", destDir)
	c5 := fmt.Sprintf(runCmd, binDir)

	script := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n", c1, c2, c3, c4, c5)
	ioutil.WriteFile("generated_script.sh", []byte(script), 0777)
	fmt.Println(script)

	exec.Command("bash", "-c", c1).Run()
	exec.Command("bash", "-c", c2).Run()
	exec.Command("bash", "-c", c3).Run()
	exec.Command("bash", "-c", c4).Run()
	finalOutput, _ := exec.Command("bash", "-c", c5).CombinedOutput()
	fmt.Println(string(finalOutput))
	fmt.Printf("Btw, we are in %s\n", cwd())

}
