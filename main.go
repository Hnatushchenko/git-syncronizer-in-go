package main

import (
	"log"
	"os/exec"
)

const sourcePath = "C:\\test4\\screens"
const destinationPath = "C:\\test4\\destination"

func main() {
	executeChangeDirectory()
	app := "git"
	arg0 := "status"
	cmd := exec.Command(app, arg0)
	stdout, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(stdout))
}

func executeChangeDirectory() {
	app := "cd"
	arg0 := "C:\\Users\\gnatu\\source\\repos\\L35\\github\\luminaryone-github-1"
	cmd := exec.Command(app, arg0)
	stdout, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(stdout))
	log.Println("Changed directory")
}
