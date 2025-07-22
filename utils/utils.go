package utils

import (
	"fmt"
	"os/exec"
)

func CheckDependencies() error {
	if _, err := exec.LookPath("fzf"); err != nil {
		return fmt.Errorf("fzf not found: %v", err)
	}
	if _, err := exec.LookPath("markitdown"); err != nil {
		return fmt.Errorf("markitdown not found: %v", err)
	}
	return nil
}
