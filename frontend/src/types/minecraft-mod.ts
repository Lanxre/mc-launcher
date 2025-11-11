export interface ModDependency {
	ModPageLink: string;
	Name: string;
	Dependency: ModDependency[];
	Details: DownloadInfo[];

	convertValues: any;
}

export interface DownloadInfo {
	URL: string;
	Version: string;
	Loader: string;
	Downloads: string;
}

export interface MinecraftMod {
	Name: string;
	Icon: string;
	ModPageLink: string;
	Description: string;
	Versions: string[];
	Screenshots: string[];
	Loaders: string[];
	Dependency: ModDependency[];
	Details: DownloadInfo[];

	convertValues: any;
}
