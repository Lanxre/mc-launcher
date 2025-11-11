package main

import (
	"embed"

	"github.com/lanxre/mc-launcher/backend/filetools"
	"github.com/lanxre/mc-launcher/backend/functools"
	"github.com/lanxre/mc-launcher/backend/parser"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	minecraftModsParser := parser.NewScraperService()
	funcService := functools.NewFuncService()
	fileService := filetools.NewFileService()

	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "MC-LAUNCHER",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.startup,
		Bind: []interface{}{
			app, minecraftModsParser, funcService, fileService,
		},
		Windows: &windows.Options{
			WebviewIsTransparent:              true,
			WindowIsTranslucent:               false,
			DisableWindowIcon:                 false,
			DisableFramelessWindowDecorations: false,
			WebviewUserDataPath:               "",
			WebviewBrowserPath:                "",
			Theme:                             windows.Dark,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
