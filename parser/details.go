package parser

import (
	"fmt"
	"slices"
	"strings"

	"github.com/gocolly/colly/v2"
)

func ScrapDetails(link string) (MinecraftMod, error) {
	var mod MinecraftMod

	c := newCollector()
	screenshots := setupScreenshotHandler(c)
	details := setupDetailsHandler(c)
	dependencies := setupDependenciesHandler(c)

	c.Visit(link)
	c.Wait()

	mod.Screenshots = processScreenshots(screenshots())
	mod.Details = details()
	mod.Dependency = dependencies()

	return mod, nil
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

