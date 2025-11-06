package parser

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
