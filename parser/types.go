package parser

type ModDependency struct {
    ModPageLink string `yaml:"mod_page_link"`
    URL         string `yaml:"url"`
    Name        string `yaml:"name"`
    Version     string `yaml:"version"`
    Loader      string `yaml:"loader"`
    Dependency  []ModDependency `yaml:"dependencies"`
    Details     []DownloadInfo `yaml:"details"`
}

type DownloadInfo struct {
    URL       string `yaml:"url"`
    Version   string `yaml:"version"`
    Loader    string `yaml:"loader"`
    Downloads string `yaml:"downloads"`
}

type MinecraftMod struct {
    Name        string   `yaml:"name"`
    Icon        string   `yaml:"icon"`
    ModPageLink string   `yaml:"mod_page_link"`
    Description string   `yaml:"description"`
    Versions    []string `yaml:"versions"`
    Screenshots []string `yaml:"screenshots"`
    Loaders     []string `yaml:"loaders"`
    Dependency  []ModDependency `yaml:"dependencies"`
    Details     []DownloadInfo `yaml:"details"`
}