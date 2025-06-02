package main

import (
	"flag"
	"fmt" // Added for filepath.WalkDir
	"os"
	"os/exec"
	"path/filepath"
	// Added for string manipulation
)

func main() {
	// Parse command-line arguments
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <project-name>\n", os.Args[0])
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	projectName := flag.Arg(0)
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	projectPath := filepath.Join(currentDir, projectName)

	// Check if the directory already exists
	if _, err := os.Stat(projectPath); err == nil {
		fmt.Fprintf(os.Stderr, "Error: Directory %s already exists\n", projectPath)
		os.Exit(1)
	}

	fmt.Printf("Creating new Prabogo project: %s\n", projectName)

	// Clone the repository
	gitCmd := exec.Command("git", "clone", "https://github.com/prabogo/prabogo.git", projectName)
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr
	if err := gitCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error cloning repository: %v\n", err)
		os.Exit(1)
	}

	// Remove the .git directory to start fresh
	err = os.RemoveAll(filepath.Join(projectPath, ".git"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error removing .git directory: %v\n", err)
		// Not exiting, as this is not always critical for the rest of the setup
	}

	fmt.Printf("\nProject created successfully!\n")
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Printf("  cp .env.example .env\n")
	fmt.Printf("  go mod tidy\n")
	fmt.Printf("  docker-compose up -d\n")
	fmt.Printf("  make run\n")
}
