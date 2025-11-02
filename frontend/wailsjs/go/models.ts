export namespace main {
	
	export class DownloadInfo {
	    URL: string;
	    Version: string;
	    Loader: string;
	    LoaderType: string;
	    Downloads: string;
	    Screenshots: string[];
	
	    static createFrom(source: any = {}) {
	        return new DownloadInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.URL = source["URL"];
	        this.Version = source["Version"];
	        this.Loader = source["Loader"];
	        this.LoaderType = source["LoaderType"];
	        this.Downloads = source["Downloads"];
	        this.Screenshots = source["Screenshots"];
	    }
	}
	export class MinecraftMod {
	    Name: string;
	    Icon: string;
	    ModPageLink: string;
	    Description: string;
	    Versions: string[];
	    Screenshots: string[];
	    Details: DownloadInfo[];
	
	    static createFrom(source: any = {}) {
	        return new MinecraftMod(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Icon = source["Icon"];
	        this.ModPageLink = source["ModPageLink"];
	        this.Description = source["Description"];
	        this.Versions = source["Versions"];
	        this.Screenshots = source["Screenshots"];
	        this.Details = this.convertValues(source["Details"], DownloadInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

