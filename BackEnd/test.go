package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func execute(code string) error {
	fmt.Println("started executing code")
	cmd := exec.Command("python", "-c", code)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Print("output : ", string(output))
		fmt.Println("error : ", err.Error())
		return err
	}
	fmt.Print("output : ", string(output))
	return nil
}

func runPythonCode(code string) (string, time.Duration, error) {
	start := time.Now()

	cmd := exec.Command("python", "-c", code)
	output, err := cmd.CombinedOutput()

	elapsed := time.Since(start)

	if err != nil {
		return "", elapsed, fmt.Errorf("failed to execute Python code: %v", err)
	}

	return string(output), elapsed, nil
}

func getMemoryUsage() (string, error) {
	cmd := exec.Command("ps", "u", "-p", fmt.Sprintf("%d", os.Getpid()))
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get memory usage: %v", err)
	}
	return string(out), nil
}
