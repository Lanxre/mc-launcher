package parser

const COUNT_RETRY = 3

func ScrapDetails(link string, versions []string) (MinecraftMod, error) {
	var mod MinecraftMod

	c := newCollector()
	screenshots := setupScreenshotHandler(c)
	details := setupDetailsHandler(c, versions)
	dependencies := setupDependenciesHandler(c)

	c.Visit(link)
	c.Wait()

	mod.Screenshots = processScreenshots(screenshots())
	mod.Details = details()
	mod.Dependency = dependencies()
	return mod, nil
}

func ScrapeMinecraftModDetails(modUrl string) []MinecraftModDetails {
	c := newCollectorWithRetry(COUNT_RETRY)

	minecraftModDeatails := setupMinecraftModDetails(c)

	c.Visit(modUrl)
	c.Wait()

	return minecraftModDeatails()
}