package parser

import (
	"fmt"
	"strings"
)

type ScraperService struct{}

func NewScraperService() *ScraperService {
	return &ScraperService{}
}

func (s *ScraperService) GetMods() ([]MinecraftMod, error) {
	url := "https://minecraft-inside.ru/mods/page/1/"
	return ScrapeMinecraftInsideModsFull(url)
}

func (s *ScraperService) GetModsByPage(page int, inputSearch *string) ([]MinecraftMod, error) {
	var url string
	if inputSearch != nil {
		searchedValue := strings.ReplaceAll(*inputSearch, " ", "+")
		url = fmt.Sprintf("https://minecraft-inside.ru/mods/page/%d/?q=%s", page, searchedValue)
	} else {
		url = fmt.Sprintf("https://minecraft-inside.ru/mods/page/%d/", page)
	}

	return ScrapeMinecraftInsideModsFull(url)
}

func (s *ScraperService) GetModDetails(link string, versions []string) (MinecraftMod, error) {
	return ScrapDetails(link, versions)
}

func (s *ScraperService) GetSearchMods(searchedValue string, page int) ([]MinecraftMod, error) {
	searchedValue = strings.ReplaceAll(searchedValue, " ", "+")
	url := fmt.Sprintf("https://minecraft-inside.ru/mods/page/%d/?q=%s", page, searchedValue)

	return ScrapeMinecraftInsideModsFull(url)
}

func (s *ScraperService) GetModDepends(depends []ModDependency, versions []string) []ModDependency {
	return ScrapeDependency(depends, versions)
}

func (s *ScraperService) GetMinecraftModDetailsV1(modUrl string) []MinecraftModDetails {
	modDetails := ScrapeMinecraftModDetails(modUrl)
	return modDetails
}
