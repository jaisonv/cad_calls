package cad

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// Client handles communication with the CAD API via Python script
type Client struct {
	PythonScriptPath string
	ConfigPath       string
}

// NewClient creates a new CAD API client that uses the Python script
func NewClient(pythonScriptPath, configPath string) *Client {
	return &Client{
		PythonScriptPath: pythonScriptPath,
		ConfigPath:       configPath,
	}
}

// GetActiveCalls fetches active CAD calls by executing the Python script
func (c *Client) GetActiveCalls(take int) (*CADResponse, error) {
	// Determine which Python interpreter to use
	pythonCmd := c.findPythonInterpreter()

	// Execute the Python script
	cmd := exec.Command(pythonCmd, c.PythonScriptPath,
		"--take", fmt.Sprintf("%d", take),
		"--open",
		"--quiet",
	)

	// Set the working directory to the script location
	scriptDir := filepath.Dir(c.PythonScriptPath)
	cmd.Dir = scriptDir

	// Run the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute Python script: %w\nOutput: %s", err, string(output))
	}

	// Find the most recent output file in cadcalls_results directory
	resultsDir := filepath.Join(scriptDir, "cadcalls_results")
	files, err := os.ReadDir(resultsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read results directory: %w", err)
	}

	// Filter JSON files and sort by modification time
	var jsonFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			jsonFiles = append(jsonFiles, filepath.Join(resultsDir, file.Name()))
		}
	}

	if len(jsonFiles) == 0 {
		return nil, fmt.Errorf("no JSON output files found in %s", resultsDir)
	}

	// Sort by modification time (newest first)
	sort.Slice(jsonFiles, func(i, j int) bool {
		infoI, _ := os.Stat(jsonFiles[i])
		infoJ, _ := os.Stat(jsonFiles[j])
		return infoI.ModTime().After(infoJ.ModTime())
	})

	// Read the most recent file
	newestFile := jsonFiles[0]
	data, err := os.ReadFile(newestFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read output file %s: %w", newestFile, err)
	}

	// Parse the JSON response
	var cadResp CADResponse
	if err := json.Unmarshal(data, &cadResp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON from %s: %w", newestFile, err)
	}

	return &cadResp, nil
}

// findPythonInterpreter finds the best Python interpreter to use
// Prefers venv Python if available, falls back to system python3
func (c *Client) findPythonInterpreter() string {
	// Get the directory of the Python script
	scriptDir := filepath.Dir(c.PythonScriptPath)

	// Check for venv in the script's directory
	venvPython := filepath.Join(scriptDir, "venv", "bin", "python3")
	if _, err := os.Stat(venvPython); err == nil {
		return venvPython
	}

	// Check for venv in parent directory (common setup)
	parentDir := filepath.Dir(scriptDir)
	venvPython = filepath.Join(parentDir, "venv", "bin", "python3")
	if _, err := os.Stat(venvPython); err == nil {
		return venvPython
	}

	// Fall back to system python3
	return "python3"
}

// NewClientFromConfig creates a client using config.py settings
// This is a convenience function that uses the Python script from the parent directory
func NewClientFromConfig(projectRoot string) *Client {
	pythonScript := filepath.Join(projectRoot, "direct_api_post.py")
	configPath := filepath.Join(projectRoot, "config.py")
	return NewClient(pythonScript, configPath)
}
