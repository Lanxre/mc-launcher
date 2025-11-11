package parser

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
)

func ScrapeDependency(depends []ModDependency) []ModDependency {
	c := newDependencyCollector()

	results := make(map[string]*ModDependency)
	mu := &sync.Mutex{}

	setupDependencyHandlers(c, results, mu)

	for i := range depends {
		results[depends[i].ModPageLink] = &depends[i]
		c.Visit(depends[i].ModPageLink)
	}

	c.Wait()
	return depends
}

func ScrapeMinecraftInsideModsFull(url string) ([]MinecraftMod, error) {

	log.Printf("🔍 Scraping mods list: %s", url)

	c := newCollectorWithRetry(3)

	var mods []MinecraftMod
	c.OnHTML("div.box.box_grass.post", func(e *colly.HTMLElement) {
		mod := ParseModBlock(e)
		if mod.Name != "" {
			mods = append(mods, mod)
		}
	})

	if err := c.Visit(url); err != nil {
		return nil, fmt.Errorf("failed to visit page: %w", err)
	}

	c.Wait()

	return mods, nil
}

func ScrapeMinecraftPageMod(mod *MinecraftMod) error {
	if mod.ModPageLink == "" {
		return fmt.Errorf("mod '%s' has empty link", mod.Name)
	}

	c := newCollectorWithRetry(3)

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
			mod.Details = append(mod.Details, download)
		}
	})

	c.OnHTML("div.entry-content img[src*='/uploads/files/']:not([src*='/mini/'])", func(e *colly.HTMLElement) {
		src := e.Request.AbsoluteURL(e.Attr("src"))
		mod.Screenshots = append(mod.Screenshots, src)
	})

	if err := c.Visit(mod.ModPageLink); err != nil {
		return fmt.Errorf("failed to visit mod page: %w", err)
	}

	c.Wait()
	return nil
}

func ParseModBlock(e *colly.HTMLElement) MinecraftMod {
	name, versions := nameParser(strings.TrimSpace(e.ChildText("h2.box__title a")), "[")
	return MinecraftMod{
		Name:        name,
		Versions:    nameVersionParse(versions),
		Loaders:     e.ChildAttrs("i.icon", "title"),
		ModPageLink: e.Request.AbsoluteURL(e.ChildAttr("h2.box__title a", "href")),
		Icon:        e.Request.AbsoluteURL(e.ChildAttr("a.post__cover img", "src")),
		Description: cleanDescription(e.ChildText("div.box__body > div:first-child"), name),
	}
}
