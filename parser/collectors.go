package parser

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/gocolly/colly/v2"
)

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

func newCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
		colly.AllowedDomains("minecraft-inside.ru"),
		colly.Async(true),
	)
	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	})
	return c
}

func newDependencyCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.Async(true),
		colly.AllowedDomains("minecraft-inside.ru"),
	)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*minecraft-inside.*",
		Parallelism: 3,
	})
	return c
}