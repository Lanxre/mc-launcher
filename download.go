package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"time"
	"strings"
)

type FileService struct {}

func NewFileService() *FileService {
    return &FileService{}
}

func (fs *FileService) DownloadFileToMinecraftMods(url, modName, version string) error {
    filename := fmt.Sprintf("%s_%s.jar", strings.ToLower(strings.ReplaceAll(modName, " ", "_", )), version)
    client := initHTTPClient()
    downloadFileFromURL(client, url, filename)

    return nil
}

func initHTTPClient() *http.Client {
	jar, _ := cookiejar.New(nil)
	
	return &http.Client{
		Timeout: 30 * time.Second,
		Jar:     jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func downloadFileFromURL(client *http.Client, url, filename string) error {
	fmt.Printf("Attempting to download from: %s\n", url)
	
	resp, err := makeRequest(client, url)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

    defer resp.Body.Close()

    switch resp.StatusCode {
	case http.StatusFound, http.StatusMovedPermanently:
		return handleRedirect(client, resp, filename)
	case http.StatusOK:
		return handleOKResponse(client, resp, filename)
	default:
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}

func makeRequest(client *http.Client, url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	setHeaders(req)
	return client.Do(req)
}

func setHeaders(req *http.Request) {
	headers := map[string]string{
		"authority":          "minecraft-inside.ru",
		"accept":             "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"accept-encoding":    "gzip, deflate, br, zstd",
		"accept-language":    "ru,en;q=0.9,ru-RU;q=0.8,en-US;q=0.7",
		"cookie":             "mobile=eaa1dd758d8c5d78da751fa16e1c5fa0e5a6fa4af01ccc4c15d2aa7a8fd6211ba%3A2%3A%7Bi%3A0%3Bs%3A6%3A%22mobile%22%3Bi%3A1%3Bs%3A2%3A%22no%22%3B%7D; adtech_uid=c2474981-54ad-49ef-bede-1f1b1a39b5e0%3Aminecraft-inside.ru; top100_id=t1.3039121.1529518743.1761071959165; tmr_lvid=e42a43b57362ce104f124ee371cc8799; tmr_lvidTS=1761071959185; _ym_uid=1761071959560828649; _ym_d=1761071959; domain_sid=SbXGA_cJpSsfai4e5ltdp%3A1762107732675; _ym_isad=1; cf_clearance=FlBtJCOylRx13LpDXd0k3KUatq60KvUmRdeENX3tHPo-1762168464-1.2.1.1-IerBotIhw91Up59Rz2z.SlzGa5TQeGve1gy1yONNMnfn.2UJy1iNkcnjAbXMhnprN0InoLbt.YXKHMTOeIXPIZVZoJ4i1m1zG8.IqbvvfFmH5hWSzqGrsnKrGxIsCq5dFr0tRJeqvPghOkbRzjFIoaaPLzz8A5SRfI.i7Rdr0yP3LUYR.Z9M2STL1ZHZBjrnkfKYjnfi.lz3NG9dR_pyE58PvFZwcYGdiq06Oeq295c; tmr_detect=1%7C1762176413410; t3_sid_3039121=s1.868198389.1762176411167.1762176458086.16.12.2.1..",
		"priority":           "u=0, i",
		"referer":            "https://minecraft-inside.ru/mods/188701-yori3os-grappling-hooks.html",
		"sec-ch-ua":          `"Google Chrome";v="141", "Not?A_Brand";v="8", "Chromium";v="141"`,
		"sec-ch-ua-arch":     `"x86"`,
		"sec-ch-ua-bitness":  `"64"`,
		"sec-ch-ua-full-version": `"141.0.7390.122"`,
		"sec-ch-ua-full-version-list": `"Google Chrome";v="141.0.7390.122", "Not?A_Brand";v="8.0.0.0", "Chromium";v="141.0.7390.122"`,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-model":    `""`,
		"sec-ch-ua-platform": `"Linux"`,
		"sec-ch-ua-platform-version": `"6.17.5"`,
		"sec-fetch-dest":     "document",
		"sec-fetch-mode":     "navigate",
		"sec-fetch-site":     "same-origin",
		"sec-fetch-user":     "?1",
		"upgrade-insecure-requests": "1",
		"user-agent":         "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36",
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}
}

func handleRedirect(client *http.Client, resp *http.Response, filename string) error {
	location := resp.Header.Get("Location")
	if location == "" {
		return fmt.Errorf("redirect without Location header")
	}
	
	return downloadDirectFile(client, location, filename)
}

func handleOKResponse(client *http.Client, resp *http.Response, filename string) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	
	bodyStr := string(body)
	
	if isCloudflarePage(bodyStr) {
		return fmt.Errorf("cloudflare protection detected")
	}
	
	downloadURL := extractDownloadURL(bodyStr)
	if downloadURL != "" {
		return downloadDirectFile(client, downloadURL, filename)
	}
	
	return fmt.Errorf("no download link found in response")
}

func isCloudflarePage(body string) bool {
	indicators := []string{
		"cloudflare",
		"challenge",
		"ray id",
		"Checking your browser",
		"DDoS protection",
	}
	
	lowerBody := strings.ToLower(body)
	for _, indicator := range indicators {
		if strings.Contains(lowerBody, strings.ToLower(indicator)) {
			return true
		}
	}
	return false
}

func extractDownloadURL(body string) string {
	patterns := []string{
		`href="([^"]*\.jar)"`,
		`location\.href=['"]([^'"]*\.jar)['"]`,
		`download-link.*?href="([^"]+)"`,
		`https://minecraft-inside\.ru/uploads/files/[^"']+\.jar`,
	}
	
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(body)
		if len(matches) > 1 {
			return matches[1]
		}
	}
	return ""
}

func downloadDirectFile(client *http.Client, fileURL, filename string) error {
	req, err := http.NewRequest("GET", fileURL, nil)
	if err != nil {
		return err
	}
	
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36")
	req.Header.Set("Accept", "*/*")
	
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

    var finalPath string
    if runtime.GOOS == "windows" {
        appData := os.Getenv("APPDATA")
        if appData == "" {
            return fmt.Errorf("APPDATA environment variable not found")
        }
        finalPath = filepath.Join(appData, ".minecraft", "mods", filename)
    } else {
        homeDir, err := os.UserHomeDir()
        if err != nil {
            return fmt.Errorf("failed to get user home directory: %v", err)
        }
        finalPath = filepath.Join(homeDir, ".minecraft", "mods", filename)
    }

    dir := filepath.Dir(finalPath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return fmt.Errorf("failed to create directory %s: %v", dir, err)
    }
	
	file, err := os.Create(finalPath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	_, err = io.Copy(file, resp.Body)
	return err
}