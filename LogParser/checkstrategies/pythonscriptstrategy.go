package checkstrategies

import (
	"fmt"
	"os/exec"
)

// PythonScriptStrategy struct for Python script check strategy
type PythonScriptStrategy struct {
	ScriptPath string
}

func (p *PythonScriptStrategy) Execute(logPath string) (string, error) {
	cmd := exec.Command("python", p.ScriptPath, logPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run Python script: %v\nOutput: %s", err, string(output))
	}
	return string(output), nil
}
