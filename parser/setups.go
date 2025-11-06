package parser

import (
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
)

func setupDependencyHandlers(c *colly.Collector, results map[string]*ModDependency, mu *sync.Mutex) {
	setupDependencyDetailsHandler(c, results, mu)
	setupSubDependenciesHandler(c, results, mu)
}

func setupDependencyDetailsHandler(c *colly.Collector, results map[string]*ModDependency, mu *sync.Mutex) {
	c.OnHTML("td.dl__info", func(e *colly.HTMLElement) {
		parentURL := e.Request.URL.String()

		download := DownloadInfo{
			URL:       e.Request.AbsoluteURL(e.ChildAttr("a", "href")),
			Version:   parseVersion(e.ChildText("span.dl__name")),
			Downloads: parseDownloadCount(e.ChildAttr("span.dl__link", "title")),
		}

		var loaders []string
		e.ForEach("span.dl__loader", func(_ int, el *colly.HTMLElement) {
			loaders = append(loaders, strings.TrimSpace(el.Text))
		})
		download.Loader = strings.Join(loaders, ", ")

		if download.URL != "" {
			mu.Lock()
			if mod, ok := results[parentURL]; ok {
				mod.Details = append(mod.Details, download)
			}
			mu.Unlock()
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

			if strings.HasPrefix(dep.ModPageLink, "https://minecraft-inside.ru/mods"){
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

func setupDetailsHandler(c *colly.Collector) func() []DownloadInfo {
	var details []DownloadInfo

	c.OnHTML("td.dl__info", func(e *colly.HTMLElement) {
		download := DownloadInfo{
			URL:       e.Request.AbsoluteURL(e.ChildAttr("a", "href")),
			Version:   parseVersion(e.ChildText("span.dl__name")),
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
				ModPageLink: fmt.Sprintf("https://minecraft-inside.ru%s", href),
				Name:        text,
			}

			if dep.Name != ".minecraft" && dep.Name != "Minecraft Forge" {
				depends = append(depends, dep)
			}
			isDependsFound = true
		})
	})

	return func() []ModDependency { return depends }
}