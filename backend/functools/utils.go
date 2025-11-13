package functools

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func GetMinecraftModsPath() (string, error) {
	mcPath, err := GetMinecraftPath()
	if err != nil {
		return "", err
	}

	return filepath.Join(mcPath, "mods"), nil
}

func GetMinecraftModPath(filename string) (string, error) {
	finalPath, err := GetMinecraftModsPath()

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", finalPath, filename), nil
}

func GetMinecraftPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch runtime.GOOS {
	case "windows":
		return filepath.Join(home, "AppData", "Roaming", ".minecraft"), nil
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "minecraft"), nil
	default:
		return filepath.Join(home, ".minecraft"), nil
	}
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
