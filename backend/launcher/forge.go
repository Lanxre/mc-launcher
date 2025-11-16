package launcher

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lanxre/mc-launcher/backend/functools"
)

type ForgeInstaller struct {
	Processors []struct {
		Sides []string `json:"sides"`
		Jar   string   `json:"jar"`
		Class string   `json:"class"`
		Args  []string `json:"args"`
	} `json:"processors"`
}

func (l *LauncherService) LaunchMinecraft(version string) error {
	mcDir, err := functools.GetMinecraftPath()
	if err != nil {
		return err
	}

	if strings.Contains(version, "forge") {
		if err := installForgeIfNeeded(mcDir, version); err != nil {
			return err
		}
	}

	javaPath := GetJava()

	return startMinecraft(javaPath, mcDir, version)
}

func installForgeIfNeeded(mcDir, version string) error {
	base, build := parseForgeVersion(version)
	versionDir := filepath.Join(mcDir, "versions", version)
	if _, err := os.Stat(versionDir); err == nil {
		fmt.Println("Forge уже установлен.")
		return nil
	}

	fmt.Printf("📦 Устанавливаем Forge %s...\n", version)

	installerURL := fmt.Sprintf(
		"https://maven.minecraftforge.net/net/minecraftforge/forge/%s-%s/forge-%s-%s-installer.jar",
		base, build, base, build,
	)
	installerPath := filepath.Join(mcDir, "installers", fmt.Sprintf("forge-%s-installer.jar", version))
	os.MkdirAll(filepath.Dir(installerPath), 0755)

	if err := DownloadFile(installerURL, installerPath); err != nil {
		return fmt.Errorf("ошибка загрузки Forge: %w", err)
	}

	javaPath := GetJava()
	cmd := exec.Command(javaPath, "-jar", installerPath, "--installClient", mcDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ошибка установки Forge: %w", err)
	}

	fmt.Println("✅ Forge установлен успешно!")
	return nil
}

func parseForgeVersion(v string) (base, build string) {
	parts := strings.Split(v, "-forge-")
	if len(parts) < 2 {
		return v, ""
	}
	return parts[0], parts[1]
}

func startMinecraft(javaPath, mcDir, version string) error {
	versionJSONPath := filepath.Join(mcDir, "versions", version, version+".json")
	data, err := os.ReadFile(versionJSONPath)
	if err != nil {
		return err
	}

	var v struct {
		MainClass string `json:"mainClass"`
		Arguments struct {
			Game []interface{} `json:"game"`
			JVM  []interface{} `json:"jvm"`
		} `json:"arguments"`
		Libraries []struct {
			Name string `json:"name"`
		} `json:"libraries"`
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	args := []string{"-Xmx4G", "-Xms2G"}
	for _, a := range v.Arguments.JVM {
		args = append(args, fmt.Sprintf("%v", a))
	}

	var paths []string
	for _, lib := range v.Libraries {
		parts := strings.Split(lib.Name, ":")
		if len(parts) < 3 {
			continue
		}
		path := filepath.Join(mcDir, "libraries", strings.ReplaceAll(parts[0], ".", "/"), parts[1], parts[2], parts[1]+"-"+parts[2]+".jar")
		if _, err := os.Stat(path); err == nil {
			paths = append(paths, path)
		}
	}

	cp := strings.Join(paths, string(os.PathListSeparator))
	args = append(args, "-cp", cp, v.MainClass)

	gameArgs := []string{
		"--username", "SuperVitalij",
		"--version", version,
		"--gameDir", mcDir,
		"--assetsDir", filepath.Join(mcDir, "assets"),
		"--uuid", "00000000-0000-0000-0000-000000000000",
		"--accessToken", "0",
		"--userType", "legacy",
	}
	for _, a := range v.Arguments.Game {
		gameArgs = append(gameArgs, fmt.Sprintf("%v", a))
	}
	args = append(args, gameArgs...)

	cmd := exec.Command(javaPath, args...)
	cmd.Dir = mcDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Launching %s...\n", version)
	return cmd.Run()
}
