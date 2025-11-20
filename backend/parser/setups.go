package parser

import (
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
)

func setupDependencyHandlers(c *colly.Collector, results map[string]*ModDependency, mu *sync.Mutex, versions []string) {
	setupDependencyDetailsHandler(c, results, versions)
	setupSubDependenciesHandler(c, results, mu)
}

func setupDependencyDetailsHandler(c *colly.Collector, results map[string]*ModDependency, versions []string) {
	var details []DownloadInfo
	c.OnHTML("script", func(e *colly.HTMLElement) {
		parentURL := e.Request.URL.String()
		script := e.Text
		if !strings.Contains(script, "var dbox_data =") {
			return
		}

		re := regexp.MustCompile(`files\s*:\s*(\[\s*{[\s\S]*?\}\s*\])`)
		matches := re.FindStringSubmatch(script)
		if len(matches) < 2 {
			return
		}

		jsonRaw := matches[1]

		cleaner := regexp.MustCompile(`</?span[^>]*>`)
		jsonRaw = cleaner.ReplaceAllString(jsonRaw, "")
		jsonRaw = strings.ReplaceAll(jsonRaw, "\\\"", "\"")
		jsonRaw = strings.ReplaceAll(jsonRaw, "\\", "")

		var files []map[string]any
		if err := json.Unmarshal([]byte(jsonRaw), &files); err != nil {
			return
		}

		details = details[:0]

		for _, f := range files {
			name := f["name"].(string)
			fileID := f["id"].(string)
			downloads := getString(f, "downloads")

			loaders := f["loaders"].([]interface{})
			loader := loaders[0].(string)

			parsedVerson := parseVersion(name)
			parsedSplitVersion := strings.Split(parsedVerson, ",")
			isFound := false

			for _, pVersion := range parsedSplitVersion {
				for _, version := range versions {
					if version == pVersion {
						isFound = true
					}
				}
			}

			if mod, ok := results[parentURL]; ok && isFound {
				mod.Details = append(mod.Details, DownloadInfo{
					URL:       "https://minecraft-inside.ru/download/" + fileID + "/",
					Version:   parsedVerson,
					Loader:    loader,
					Downloads: downloads,
				})
			}
		}
		

	})
}

func setupSubDependenciesHandler(c *colly.Collector, results map[string]*ModDependency, mu *sync.Mutex) {
	c.OnHTML("div.box__body", func(e *colly.HTMLElement) {
		parentURL := e.Request.URL.String()

		var subDeps []ModDependency
		e.ForEach("ol li", func(_ int, li *colly.HTMLElement) {
			link := li.DOM.Find("a").First()
			if link.Length() == 0 {
				return
			}

			href, _ := link.Attr("href")
			text := strings.TrimSpace(link.Text())
			if text == ".minecraft" || text == "Minecraft Forge" {
				return
			}

			dep := ModDependency{
				ModPageLink: fmt.Sprintf("https://minecraft-inside.ru%s", href),
				Name:        text,
			}

			if strings.HasPrefix(dep.ModPageLink, "https://minecraft-inside.ru/mods") {
				subDeps = append(subDeps, dep)
			}

		})

		if len(subDeps) > 0 {
			mu.Lock()
			if mod, ok := results[parentURL]; ok {
				mod.Dependency = append(mod.Dependency, subDeps...)
			}
			mu.Unlock()
		}
	})
}

func setupScreenshotHandler(c *colly.Collector) func() []string {
	var screenshots []string

	c.OnHTML("div.box__body img[src*='/uploads/files/']", func(e *colly.HTMLElement) {
		imgSrc := e.Attr("src")
		absoluteURL := e.Request.AbsoluteURL(imgSrc)
		if !slices.Contains(screenshots, absoluteURL) {
			screenshots = append(screenshots, absoluteURL)
		}
	})

	return func() []string { return screenshots }
}

func setupDetailsHandler(c *colly.Collector, versions []string) func() []DownloadInfo {
	var details []DownloadInfo

	c.OnHTML("td.dl__info", func(e *colly.HTMLElement) {

		parsedVersion := parseVersion(e.ChildText("span.dl__name"))
		if !slices.Contains(versions, parsedVersion) {
			return
		}

		download := DownloadInfo{
			URL:       e.Request.AbsoluteURL(e.ChildAttr("a", "href")),
			Version:   parsedVersion,
			Downloads: parseDownloadCount(e.ChildAttr("span.dl__link", "title")),
		}

		var loaders []string
		e.ForEach("span.dl__loader", func(_ int, el *colly.HTMLElement) {
			loaders = append(loaders, strings.TrimSpace(el.Text))
		})
		download.Loader = strings.Join(loaders, ", ")

		if download.URL != "" {
			details = append(details, download)
		}
	})

	return func() []DownloadInfo { return details }
}

func setupMinecraftModDetails(c *colly.Collector) func() []MinecraftModDetails {
	var details []MinecraftModDetails

	c.OnHTML("script", func(e *colly.HTMLElement) {
		script := e.Text
		if !strings.Contains(script, "var dbox_data =") {
			return
		}

		re := regexp.MustCompile(`files\s*:\s*(\[\s*{[\s\S]*?\}\s*\])`)
		matches := re.FindStringSubmatch(script)
		if len(matches) < 2 {
			return
		}

		jsonRaw := matches[1]

		cleaner := regexp.MustCompile(`</?span[^>]*>`)
		jsonRaw = cleaner.ReplaceAllString(jsonRaw, "")
		jsonRaw = strings.ReplaceAll(jsonRaw, "\\\"", "\"")
		jsonRaw = strings.ReplaceAll(jsonRaw, "\\", "")

		var files []map[string]any
		if err := json.Unmarshal([]byte(jsonRaw), &files); err != nil {
			return
		}

		details = details[:0]

		for _, f := range files {
			name := f["name"].(string)
			fileID := f["id"].(string)
			size := getString(f, "size")
			date := getString(f, "created")
			downloads := getString(f, "downloads")

			loaders := f["loaders"].([]interface{})
			loader := loaders[0].(string)

			details = append(details, MinecraftModDetails{
				FileID:      fileID,
				Version:     parseVersion(name),
				Loader:      loader,
				Date:        date,
				Size:        size,
				Downloads:   downloads,
				DownloadURL: "https://minecraft-inside.ru/download/" + fileID + "/",
			})
		}

	})

	return func() []MinecraftModDetails { return details }
}

func setupDependenciesHandler(c *colly.Collector) func() []ModDependency {
	var depends []ModDependency
	isDependsFound := false

	c.OnHTML("div.box__body", func(e *colly.HTMLElement) {
		if isDependsFound {
			return
		}

		e.ForEach("ol li", func(_ int, li *colly.HTMLElement) {
			link := li.DOM.Find("a").First()
			if link.Length() == 0 {
				return
			}

			href, _ := link.Attr("href")
			text := link.Text()

			dep := ModDependency{
				ModPageLink: href,
				Name:        text,
			}

			if strings.HasPrefix(dep.ModPageLink, "/mods") {
				dep.ModPageLink = fmt.Sprintf("https://minecraft-inside.ru%s", dep.ModPageLink)
			}

			if dep.Name != ".minecraft" && dep.Name != "Minecraft Forge" && dep.Name != "Fabric" {
				depends = append(depends, dep)
			}
			isDependsFound = true
		})
	})

	return func() []ModDependency { return depends }
}
