package functools

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lanxre/mc-launcher/backend/parser"
	"gopkg.in/yaml.v3"
)

const FAVOURITES = "favourite.yaml"

func (s *FuncService) getYamlFilePath(filename string) (string, error) {
	mcPath, err := GetMinecraftPath()
	if err != nil {
		return "", fmt.Errorf("failed to get Minecraft path: %w", err)
	}

	if filepath.Ext(filename) == "" {
		filename += ".yaml"
	}
	return filepath.Join(mcPath, filename), nil
}

func (s *FuncService) writeYamlFile(filepath string, data any) error {
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	if err := os.WriteFile(filepath, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filepath, err)
	}
	return nil
}

func (s *FuncService) readModsFromFile(path string) ([]parser.MinecraftMod, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []parser.MinecraftMod{}, nil
		}
		return nil, fmt.Errorf("failed to read YAML file: %w", err)
	}

	var mods []parser.MinecraftMod
	if err := yaml.Unmarshal(data, &mods); err == nil {
		return mods, nil
	}

	var single parser.MinecraftMod
	if err := yaml.Unmarshal(data, &single); err == nil {
		return []parser.MinecraftMod{single}, nil
	}

	return nil, fmt.Errorf("invalid YAML format in %s", path)
}

func (s *FuncService) writeModsToFile(path string, mods []parser.MinecraftMod) error {
	return s.writeYamlFile(path, mods)
}

func (s *FuncService) SaveYamlModConfig(data parser.MinecraftMod, filename string) error {
	filePath, err := s.getYamlFilePath(filename)

	if err != nil {
		os.Create(filePath)
	} 
	
	mods, _ := s.readModsFromFile(filePath)
	mods = append(mods, data)
	return s.writeModsToFile(filePath, mods)
}

func (s *FuncService) GetYamlConfig(filename string) ([]parser.MinecraftMod, error) {
	filePath, err := s.getYamlFilePath(filename)
	if err != nil {
		return nil, err
	}
	return s.readModsFromFile(filePath)
}

func (s *FuncService) RemoveFromYamlConfig(mod parser.MinecraftMod, filename string) error {
	mods, err := s.GetYamlConfig(filename)
	if err != nil {
		return err
	}

	filtered := make([]parser.MinecraftMod, 0, len(mods))
	for _, m := range mods {
		if m.Name != mod.Name {
			filtered = append(filtered, m)
		}
	}

	filePath, err := s.getYamlFilePath(filename)
	if err != nil {
		return err
	}

	if err := s.writeModsToFile(filePath, filtered); err != nil {
		return fmt.Errorf("failed to update favourites: %w", err)
	}

	return nil
}
