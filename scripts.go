package main

import (
	"bufio"
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

//go:embed scripts/*
var scripts embed.FS

func runAppInitScript() {
	if !confirmPrompt(Warn(
		`The script will override any statico related files, 
this includes templates and styles,
please do not use this in an already initiatlized project, 
Do you still want to continue?`,
	)) {
		return
	}
	bootScript, err := scripts.ReadFile("scripts/bootstrap.sh")

	bail(err)
	c := exec.Command("bash")
	c.Stdin = strings.NewReader((string(bootScript)))

	_, err = c.Output()
	bail(err)
	fmt.Println(Success("Statico Minimal Site Created"))
}

func confirmPrompt(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
