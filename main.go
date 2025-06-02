package main

import (
	"flag"
	"fmt"
	"io/fs" // Added for filepath.WalkDir
	"os"
	"os/exec"
	"path/filepath"
	"strings" // Added for string manipulation
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

	// Replace prabogo with the new project name in go.mod
	goModPath := filepath.Join(projectPath, "go.mod")
	modifyGoMod(goModPath, projectName)

	// Modify import paths in all .go files in the cloned project
	oldModuleName := "prabogo"
	err = modifyProjectImportPaths(projectPath, oldModuleName, projectName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occurred while modifying project import paths: %v\n", err)
		// Depending on severity, you might want to os.Exit(1) here
	}

	fmt.Printf("\nProject created successfully!\n")
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Printf("  cp .env.example .env\n")
	fmt.Printf("  go mod tidy\n")
	fmt.Printf("  docker-compose up -d\n")
	fmt.Printf("  make run\n")
}

func modifyGoMod(path string, projectName string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading go.mod: %v\n", err)
		return
	}

	// Replace module name
	// Note: This is a simple replacement and might not work for complex module names
	content := string(data)
	newContent := "module " + projectName + "\n"
	newContent += content[len("module prabogo\n"):]

	err = os.WriteFile(path, []byte(newContent), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing go.mod: %v\n", err)
		return
	}
}

// modifyProjectImportPaths iterates over .go files in the project and updates import paths.
func modifyProjectImportPaths(projectRootPath string, oldModuleName string, newModuleName string) error {
	fmt.Printf("Updating import paths from %s to %s in project %s\n", oldModuleName, newModuleName, projectRootPath)

	oldImportPathToken := "\"" + oldModuleName + "/"
	newImportPathToken := "\"" + newModuleName + "/"

	return filepath.WalkDir(projectRootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accessing path %s: %v\n", path, err)
			return err // Propagate WalkDir errors
		}

		// Skip .git directory
		if strings.Contains(path, ".git"+string(filepath.Separator)) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Process only .go files, skip directories, go.mod, and go.sum
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".go") || strings.HasSuffix(path, "go.mod") || strings.HasSuffix(path, "go.sum") {
			return nil
		}

		contentBytes, readErr := os.ReadFile(path)
		if readErr != nil {
			fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", path, readErr)
			return nil // Continue with other files
		}

		originalContent := string(contentBytes)
		lines := strings.Split(originalContent, "\n")
		modifiedLines := make([]string, len(lines))
		changed := false
		inImportBlock := false

		for i, line := range lines {
			trimmedLine := strings.TrimSpace(line)
			newLine := line

			if strings.HasPrefix(trimmedLine, "import (") {
				inImportBlock = true
			} else if trimmedLine == ")" && inImportBlock {
				inImportBlock = false
			}

			isSingleLineImport := strings.HasPrefix(trimmedLine, "import ")
			// Covers: "path", alias "path", . "path"
			isPathInImportBlock := inImportBlock && (strings.HasPrefix(trimmedLine, "\"") || (strings.Count(trimmedLine, "\"") >= 2 && (strings.Contains(trimmedLine, " ") || strings.HasPrefix(trimmedLine, "."))))

			if isSingleLineImport || isPathInImportBlock {
				if strings.Contains(newLine, oldImportPathToken) {
					newLine = strings.ReplaceAll(newLine, oldImportPathToken, newImportPathToken)
					if newLine != line {
						changed = true
					}
				}
			}
			modifiedLines[i] = newLine
		}

		if changed {
			modifiedContent := strings.Join(modifiedLines, "\n")

			fileInfo, statErr := os.Stat(path)
			perm := fs.FileMode(0644) // Default permission
			if statErr == nil {
				perm = fileInfo.Mode()
			}

			writeErr := os.WriteFile(path, []byte(modifiedContent), perm)
			if writeErr != nil {
				fmt.Fprintf(os.Stderr, "Error writing updated content to %s: %v\n", path, writeErr)
				// Continue with other files
			} else {
				fmt.Printf("  Updated import paths in: %s\n", path)
			}
		}
		return nil
	})
}
