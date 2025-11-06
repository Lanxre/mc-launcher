package parser

type ModDependency struct {
	ModPageLink string
	URL     	string
	Name   	 	string
	Version 	string
	Loader 		string
	Dependency  []ModDependency
	Details     []DownloadInfo
}

type DownloadInfo struct {
	URL        string
	Version    string
	Loader     string
	Downloads  string
}

type MinecraftMod struct {
	Name        string
	Icon        string
	ModPageLink string
	Description string
	Versions    []string
	Screenshots []string
	Loaders     []string
	Dependency  []ModDependency
	Details     []DownloadInfo
}