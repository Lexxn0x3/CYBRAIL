package checkstrategies

import (
	"fmt"
	"os/exec"
)

// BinaryExecutableStrategy struct for binary executable check strategy
type BinaryExecutableStrategy struct {
	ExecutablePath string
}

func (b *BinaryExecutableStrategy) Execute(logPath string) (string, error) {
	cmd := exec.Command(b.ExecutablePath, logPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run binary executable: %v\nOutput: %s", err, string(output))
	}
	return string(output), nil
}
