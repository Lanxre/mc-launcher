package functools

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func GetMinecraftModsPath() (string, error) {
	var finalPath string
	if runtime.GOOS == "windows" {
		appData := os.Getenv("APPDATA")
		if appData == "" {
			return "", fmt.Errorf("APPDATA environment variable not found")
		}
		finalPath = filepath.Join(appData, ".minecraft", "mods")
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %v", err)
		}
		finalPath = filepath.Join(homeDir, ".minecraft", "mods")
	}

	return finalPath, nil
}

func GetMinecraftModPath(filename string) (string, error) {
	finalPath, err := GetMinecraftModsPath()

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", finalPath, filename), nil
}

func GetMinecraftPath() (string, error) {
	var finalPath string
	if runtime.GOOS == "windows" {
		appData := os.Getenv("APPDATA")
		if appData == "" {
			return "", fmt.Errorf("APPDATA environment variable not found")
		}
		finalPath = filepath.Join(appData, ".minecraft")
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %v", err)
		}
		finalPath = filepath.Join(homeDir, ".minecraft")
	}

	return finalPath, nil
}

func ConverModName(modName string) string {
	return strings.ToLower(strings.ReplaceAll(modName, " ", "_"))
}

func isExistModInDisk(modName string) bool {

	done := false

	minecraftPath, err := GetMinecraftModsPath()
	if err != nil {
		return done
	}

	entries, err := os.ReadDir(minecraftPath)
	if err != nil {
		return done
	}

	supossedName := ConverModName(modName)

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), supossedName) {
			return !done
		}
	}

	return done
}
