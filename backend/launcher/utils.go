package launcher

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func DownloadFile(url, path string) error {
	
	os.MkdirAll(filepath.Dir(path), 0755)
	out, _ := os.Create(path)

	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func GetJava() string {
	java, err := exec.LookPath("java")
	if err != nil {
		return "java"
	}
	return java
}