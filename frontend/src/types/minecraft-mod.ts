export interface  ModDependency {
  ModPageLink: string;
  URL: string;
  Name: string;
  Version: string;
  Loader: string;
  Dependency: ModDependency[];
  Details: DownloadInfo[];

  convertValues: any
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
  Dependency:  ModDependency[];
  Details: DownloadInfo[];

  convertValues: any
}