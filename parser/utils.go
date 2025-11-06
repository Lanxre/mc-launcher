package parser

import (
	"fmt"
	"strings"
)

func processScreenshots(urls []string) []string {
    urls = removeDuplicates(urls)
    for i, url := range urls {
        parts := strings.Split(url, ".")
        if len(parts) > 0 {
            ext := parts[len(parts)-1]
            switch ext {
            case "png", "jpg":
                urls[i] = strings.Replace(url, "mini", "thumb", 1)
            case "gif":
                urls[i] = strings.Replace(url, "/mini", "", 1)
            }
        }
    }

    if len(urls) > 5 {
        urls = urls[:len(urls)-5]
    }

    return urls
}

func removeDuplicates(urls []string) []string {
    keys := make(map[string]bool)
    var result []string
    for _, url := range urls {
        if _, value := keys[url]; !value {
            keys[url] = true
            result = append(result, url)
        }
    }
    return result
}

func nameParser(fullName, symbol string) (string, string) {
	if idx := strings.Index(fullName, symbol); idx != -1 {
		return strings.TrimSpace(fullName[:idx]), fullName[idx:]
	}
	return fullName, ""
}

func nameVersionParse(versions string) []string {
	var result []string
	for {
		start := strings.IndexByte(versions, '[')
		if start == -1 {
			break
		}
		end := strings.IndexByte(versions[start:], ']')
		if end == -1 {
			break
		}
		version := strings.TrimSpace(versions[start+1 : start+end])
		if version != "" {
			result = append(result, version)
		}
		versions = versions[start+end+1:]
	}
	return result
}

func parseVersion(text string) string {
	text = strings.TrimSpace(text)
	if idx := strings.Index(text, "<span"); idx != -1 {
		text = strings.TrimSpace(text[:idx])
	}

	text = strings.TrimPrefix(text, "Для ")
	
	versText := strings.Split(text, " ")

	if len(versText) > 2 {
		text = fmt.Sprintf("%s, %s", versText[0], versText[2])
	}else{
		text = versText[0]
	}

	return text
}

func parseDownloadCount(title string) string {
	if !strings.Contains(title, "Скачиваний:") {
		return ""
	}
	parts := strings.SplitN(title, ":", 2)
	if len(parts) < 2 {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

func cleanDescription(desc, name string) string {
	desc = strings.TrimSpace(desc)
	if strings.Contains(desc, name) {
		desc = strings.Replace(desc, name+" добавляет", "", 1)
	}
	return strings.TrimSpace(desc)
}
