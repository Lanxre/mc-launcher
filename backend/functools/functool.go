package functools

import (
	"cmp"
	"os"
	"path"
	"slices"

	"github.com/lanxre/mc-launcher/backend/parser"
)

type FuncService struct{}

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

func (s *FuncService) GetSavedMods() ([]string, error) {
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

func (s *FuncService) DeleteSavedMod(modName string) {
	minecraftDir, _ := GetMinecraftModsPath()
	finalPath := path.Join(minecraftDir, modName)

	os.Remove(finalPath)
}

func (s *FuncService) IsModExist(modName string) bool {

	onFile := isExistModInDisk(modName)

	if !onFile {
		return !onFile
	}

	modsInConfig, err := s.GetYamlConfig("downloads")
	done := false

	if err != nil {
		return done
	}

	for _, mod := range modsInConfig {
		if mod.Name == modName {
			return !done
		}
	}

	return done
}
