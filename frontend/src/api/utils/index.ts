import type { MinecraftMod, ModDependency, PairDepend } from "@/types";
import { OpenExternalLink } from "@wailsjs/go/main/App";
import {
	GetSavedMods,
	SaveYamlModConfig,
} from "@wailsjs/go/functools/FuncService";
import { GetModDetails } from "@wailsjs/go/parser/ScraperService";

export const uniqueBy = <T>(array: T[], keyFn: (item: T) => string): T[] => {
	const seen = new Set<string>();
	const result: T[] = [];

	for (const item of array) {
		const key = keyFn(item);
		if (!seen.has(key)) {
			seen.add(key);
			result.push(item);
		}
	}

	return result;
};

export const openLink = async (url: string) => {
	await OpenExternalLink(url);
};

export const saveModToYaml = async (mod: MinecraftMod, filename: string) => {
	try {
		if (mod !== null && mod !== undefined) {
			await SaveYamlModConfig(mod, filename);
		}
	} catch (err) {
		console.error("Ошибка сохранения:", err);
	}
};

export const getMinecraftDownloadFileName = (
	modName: string,
	versions: string[],
): string => {
	const normalizedModName = modName.toLowerCase().replaceAll(" ", "_");
	const normalizedVersions = versions.join("_").toLowerCase();

	return `${normalizedModName}_${normalizedVersions}.jar`;
};

export const filterDiskModDepends = async (depends: ModDependency[]) => {
	const result: PairDepend[] = [];
	const savedModsOnDisk = await GetSavedMods();

	for (const depend of depends) {
		if (!depend?.Name) continue;

		const expectedFileName = depend.Name.replaceAll(" ", "_").toLowerCase();
		const matchingFile = savedModsOnDisk.find((file) =>
			file.toLowerCase().startsWith(expectedFileName),
		);

		if (matchingFile) {
			result.push({
				configDepend: depend,
				fileDepend: matchingFile,
			});
		}
	}

	return result;
};

export const filterNoDiskModDepends = async (depends: ModDependency[]) => {
	const result: ModDependency[] = [];
	const savedModsOnDisk = await GetSavedMods();

	for (const depend of depends) {
		if (!depend?.Name) continue;

		const expectedFileName = depend.Name.replaceAll(" ", "_").toLowerCase();
		const matchingFile = savedModsOnDisk.find((file) =>
			file.toLowerCase().startsWith(expectedFileName),
		);

		if (!matchingFile) {
			result.push(depend);
		}
	}

	return result;
};

export async function enrichDependencies(deps: ModDependency[]) {
	for (const dep of deps) {
		if (dep.Details == null || dep.Details.length === 0) {
			const data = await GetModDetails(dep.ModPageLink);
			dep.Details = data.Details;
		}
	}
}
