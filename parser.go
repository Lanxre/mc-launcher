package main

import (
	"slices"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
)

type ScraperService struct{}

type DownloadInfo struct {
	URL        string
	Version    string
	Loader     string
	LoaderType string
	Downloads  string
	Screenshots []string
}

type MinecraftMod struct {
	Name        string
	Icon        string
	ModPageLink string
	Description string
	Versions    []string
	Screenshots []string
	Details     []DownloadInfo
}

func NewScraperService() *ScraperService {
	return &ScraperService{}
}

func (s *ScraperService) GetMods() ([]MinecraftMod, error) {
	return ScrapeMinecraftInsideModsFull(1)
}

func (s *ScraperService) GetModsByPage(page int) ([]MinecraftMod, error) {
	return ScrapeMinecraftInsideModsFull(page)
}

func (s *ScraperService) GetModDetails(link string) (MinecraftMod, error) {
	return ScrapDetails(link)
}

func ScrapeMinecraftInsideModsFull(page int) ([]MinecraftMod, error) {
	url := fmt.Sprintf("https://minecraft-inside.ru/mods/page/%d/", page)
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
		return nil, fmt.Errorf("failed to visit page %d: %w", page, err)
	}

	c.Wait()
	log.Printf("✅ Found %d mods on page %d", len(mods), page)

	const maxWorkers = 5
	sem := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup

	for i := range mods {
		wg.Add(1)
		sem <- struct{}{}

		go func(mod *MinecraftMod) {
			defer wg.Done()
			defer func() { <-sem }()

			time.Sleep(300 * time.Millisecond)
			if err := ScrapeMinecraftPageMod(mod); err != nil {
				log.Printf("⚠️ Error scraping mod '%s': %v", mod.Name, err)
			} else {
				log.Printf("✅ Scraped mod: %s", mod.Name)
			}
		}(&mods[i])
	}

	wg.Wait()
	log.Printf("✅ Page %d: all mods scraped successfully", page)
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

func newCollectorWithRetry(maxRetries int) *colly.Collector {
	c := colly.NewCollector(
		colly.AllowedDomains("minecraft-inside.ru"),
		colly.Async(true),
	)

	c.OnError(func(r *colly.Response, err error) {
		if r.StatusCode == 429 {
			retryAfter := 5 + (r.Request.Ctx.GetAny("retryCount").(int) * 3)
			log.Printf("⏳ Got 429 Too Many Requests. Retrying %s after %d sec...", r.Request.URL, retryAfter)
			time.Sleep(time.Duration(retryAfter) * time.Second)

			retryCount := r.Request.Ctx.GetAny("retryCount").(int)
			if retryCount < maxRetries {
				r.Request.Ctx.Put("retryCount", retryCount+1)
				if err := r.Request.Retry(); err != nil {
					log.Printf("❌ Retry failed for %s: %v", r.Request.URL, err)
				}
			} else {
				log.Printf("❌ Max retries reached for %s", r.Request.URL)
			}
		} else {
			log.Printf("⚠️ Error %d on %s: %v", r.StatusCode, r.Request.URL, err)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		if _, ok := r.Ctx.GetAny("retryCount").(int); !ok {
			r.Ctx.Put("retryCount", 0)
		}
	})

	return c
}

func ScrapDetails(link string) (MinecraftMod, error) {
    var mod MinecraftMod

    c := colly.NewCollector(
        colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
        colly.AllowedDomains("minecraft-inside.ru"),
        colly.Async(true),
    )

    c.WithTransport(&http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
    })

	var screenshots []string
	c.OnHTML("img[src*='/uploads/files/']", func(e *colly.HTMLElement) {
		imgSrc := e.Attr("src")
		fmt.Printf("Found image with uploads path: %s\n", imgSrc)
		
		if strings.Contains(imgSrc, "/mini/") {
			absoluteURL := e.Request.AbsoluteURL(imgSrc)
			found := slices.Contains(screenshots, absoluteURL)
			
			if !found {
				screenshots = append(screenshots, absoluteURL)
				fmt.Printf("✅ Added image from uploads selector: %s\n", absoluteURL)
				
			}
		}
	})


	c.OnHTML("td.dl__info", func(h *colly.HTMLElement) {
		var detail DownloadInfo
		detail.URL = h.ChildAttr("a","href")

		h.ForEach("span.dl__name", func(_ int, el *colly.HTMLElement) {
			vers := strings.TrimSpace(el.Text)
			data := strings.Split(vers, " ")
			
			detail.Version = data[1]

			if len(data) > 2 {
				detail.Loader = data[2]
			}

			
			mod.Details = append(mod.Details, detail)
		
		})
	})

    c.Visit(link)
	c.Wait()

	screenshots = processScreenshots(screenshots)
    mod.Screenshots = screenshots
    
	return mod, nil
}

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

func ParseModBlock(e *colly.HTMLElement) MinecraftMod {
	name, versions := nameParser(strings.TrimSpace(e.ChildText("h2.box__title a")), "[")

	return MinecraftMod{
		Name:        name,
		Versions:    nameVersionParse(versions),
		ModPageLink: e.Request.AbsoluteURL(e.ChildAttr("h2.box__title a", "href")),
		Icon:        e.Request.AbsoluteURL(e.ChildAttr("a.post__cover img", "src")),
		Description: cleanDescription(e.ChildText("div.box__body > div:first-child"), name),
	}
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
	return strings.TrimPrefix(text, "Для ")
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
