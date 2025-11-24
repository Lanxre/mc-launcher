package parser

import (
	"fmt"
	"strings"
)

const (
	baseURL       = "https://minecraft-inside.ru"
	modsPath      = "mods"
)

type ScraperService struct{}

func NewScraperService() *ScraperService {
	return &ScraperService{}
}

func (s *ScraperService) GetMods() ([]MinecraftMod, error) {
	return s.GetModsByPage(1, nil)
}

func (s *ScraperService) GetModsByPage(page int, inputSearch *string) ([]MinecraftMod, error) {
	url := s.buildURL(modsPath, page, inputSearch)
	return ScrapeMinecraftInsideModsFull(url)
}

func (s *ScraperService) GetModDetails(link string, versions []string) (MinecraftMod, error) {
	return ScrapDetails(link, versions)
}

func (s *ScraperService) GetSearchMods(searchedValue string, page int) ([]MinecraftMod, error) {
	return s.GetModsByPage(page, &searchedValue)
}

func (s *ScraperService) GetModDepends(depends []ModDependency, versions []string) []ModDependency {
	return ScrapeDependency(depends, versions)
}

func (s *ScraperService) GetMinecraftModDetailsV1(modUrl string) []MinecraftModDetails {
	return ScrapeMinecraftModDetails(modUrl)
}

func (s *ScraperService) buildURL(category string, page int, search *string) string {
	base := fmt.Sprintf("%s/%s/page/%d/", baseURL, category, page)
	if search != nil && *search != "" {
		query := strings.ReplaceAll(*search, " ", "+")
		return fmt.Sprintf("%s?q=%s", base, query)
	}
	return base
}