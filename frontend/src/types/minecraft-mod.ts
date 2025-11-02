export interface DownloadInfo {
  URL: string;
  Version: string;
  Loader: string;
  LoaderType: string;
  Downloads: string;
  Screenshots: string[];
}

export interface MinecraftMod {
  Name: string;
  Icon: string;
  ModPageLink: string;
  Description: string;
  Versions: string[];
  Screenshots: string[];
  Details: DownloadInfo[];

  convertValues: any
}