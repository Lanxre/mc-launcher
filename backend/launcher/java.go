package launcher

import (
	"os/exec"
	"path/filepath"
)

func GetJava() (string, error) {
	if path, err := exec.LookPath("java"); err == nil {
		return path, nil
	}

	return filepath.Join("java", "bin", "java"), nil
}