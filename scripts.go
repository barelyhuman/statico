package main

import (
	"embed"
	"fmt"
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
