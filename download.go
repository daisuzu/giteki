package main

import (
	"bufio"
	"fmt"
	"os/exec"
)

func executeDownloadScript(script string, opts []string) error {
	cmd := exec.Command(
		"python",
		append([]string{script}, opts...)...,
	)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	cmd.Start()

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	cmd.Wait()

	return nil
}

func Download(dst string, all bool, update bool) error {
	script, err := GetDownloadScript()
	if err != nil {
		return err
	}

	opts := []string{"--dst", dst}
	if all {
		opts = append(opts, "--all")
	}
	if update {
		opts = append(opts, "--update")
	}

	return executeDownloadScript(script, opts)
}
