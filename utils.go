package main

import (
	"cmp"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"slices"

	"github.com/lanxre/mc-launcher/parser"
)

type FuncService struct {}

func NewFuncService() *FuncService {
    return &FuncService{}
}

func (s *FuncService) SortByVersions(mods []parser.MinecraftMod, version string) []parser.MinecraftMod {
    var filteredMods []parser.MinecraftMod
    
    for _, mod := range mods {
        for _, v := range mod.Versions {
            if v == version {
                slices.SortFunc(mod.Versions, func(a, b string) int {
                    return cmp.Compare(a, b)
                })
                filteredMods = append(filteredMods, mod)
                break
            }
        }
    }
    
    return filteredMods
}


func (s *FuncService) SortByLoader(mods []parser.MinecraftMod, searchedLoader string) []parser.MinecraftMod {
    var filteredMods []parser.MinecraftMod
    
    for _, mod := range mods {
        if slices.Contains(mod.Loaders, searchedLoader) {
            filteredMods = append(filteredMods, mod)
        }
    }
    
    return filteredMods
}

func(s *FuncService) GetSavedMods() ([]string, error ){
	finalPath, err := GetMinecraftModsPath()

	if err != nil {
		return nil, nil
	}

	entries, err := os.ReadDir(finalPath)
    if err != nil {
        return nil, nil
    }

	filesNames := make([]string, 0, len(entries))
    for _, entry := range entries {
        if !entry.IsDir() {
            filesNames = append(
                filesNames,
                entry.Name(),
            )
        }
    }

	return filesNames, nil
}

func (s * FuncService) DeleteSavedMod(modName string) {
    minecraftDir, _ := GetMinecraftModsPath()
    finalPath := path.Join(minecraftDir, modName)

    os.Remove(finalPath)
}

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
        return  "", err
    }

    return fmt.Sprintf("%s/%s",finalPath, filename), nil
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