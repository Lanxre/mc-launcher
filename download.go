package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

type FileService struct {}

func NewFileService() *FileService {
    return &FileService{}
}

func(fs *FileService) DownloadFileToMinecraftMods(url, modName string) error {
	filename := fmt.Sprintf("%s.jar", modName)
    
    var finalPath string
    if runtime.GOOS == "windows" {
        appData := os.Getenv("APPDATA")
        if appData == "" {
            return fmt.Errorf("APPDATA environment variable not found")
        }
        finalPath = filepath.Join(appData, ".minecraft", "mods", filename)
    } else {
        finalPath = filename
    }

    dir := filepath.Dir(finalPath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return fmt.Errorf("failed to create directory %s: %v", dir, err)
    }

    out, err := os.Create(finalPath)
    if err != nil {
        return err
    }
    defer out.Close()

    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("bad status: %s", resp.Status)
    }

    _, err = io.Copy(out, resp.Body)
    if err != nil {
        return err
    }

    return nil
}