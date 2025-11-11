<script setup lang="ts">
import { ref } from "vue";
import {
	DownloadFileToMinecraftMods,
	DownloadsMods,
} from "@wailsjs/go/filetools/FileService";
import { ShowInfoMessage } from "@wailsjs/go/main/App";
import { IsModExist } from "@wailsjs/go/functools/FuncService";
import {
	saveModToYaml,
	getMinecraftDownloadFileName,
	filterNoDiskModDepends,
} from "@/api/utils";
import type { MinecraftMod, DownloadInfo, ModDependency } from "@/types";

interface Props {
	mod: MinecraftMod;
	depends: ModDependency[];
}

const props = defineProps<Props>();
const isDownloading = ref(false);

const showNotify = async (title: string, message: string): Promise<void> => {
	try {
		await ShowInfoMessage(title, message);
	} catch {
		console.warn("Не удалось показать уведомление");
	}
};

const compareVersions = (v1: string, v2: string): number => {
	const parse = (v: string): number[] =>
		v
			.split(".")
			.map((n) => Number(n))
			.filter((n) => !isNaN(n));

	const [a1 = 0, a2 = 0, a3 = 0] = parse(v1);
	const [b1 = 0, b2 = 0, b3 = 0] = parse(v2);

	if (a1 !== b1) return a1 - b1;
	if (a2 !== b2) return a2 - b2;
	return a3 - b3;
};

const getFirstVersion = (dl?: DownloadInfo): string => {
	if (!dl || typeof dl.Version !== "string") return "0.0.0";
	const first = dl.Version.split(",")[0]?.trim();
	return first && first.length ? first : "0.0.0";
};

const downloadMod = async (
	mod: MinecraftMod,
	detail: DownloadInfo,
): Promise<void> => {
	if (isDownloading.value) return;
	isDownloading.value = true;

	const isExist = await IsModExist(mod.Name);

	if (isExist) {
		isDownloading.value = false;
		await showNotify("Предупреждение", "Мод уже скачан");
		return;
	}

	try {
		const filtred = await filterNoDiskModDepends(props.depends);
		const depFiles: DownloadInfo[] = filtred
			.flatMap((dep: ModDependency): DownloadInfo[] => {
				if (!dep?.Details || !Array.isArray(dep.Details)) return [];

				const filtered = dep.Details.filter(
					(dl: DownloadInfo): dl is DownloadInfo => {
						if (
							!dl ||
							typeof dl.Version !== "string" ||
							typeof dl.Loader !== "string"
						)
							return false;

						const versions = dl.Version.split(",")
							.map((v) => v.trim())
							.filter(Boolean);
						const loaders = dl.Loader.split(",")
							.map((l) => l.trim())
							.filter(Boolean);
						if (!loaders.includes(detail.Loader)) return false;

						const valid = versions
							.filter((v) => compareVersions(v, detail.Version) <= 0)
							.sort((a, b) => compareVersions(b, a))[0];

						return Boolean(valid);
					},
				).sort((a, b) =>
					compareVersions(getFirstVersion(b), getFirstVersion(a)),
				);

				const best = filtered[0];
				return best ? [best] : [];
			})
			.filter((d): d is DownloadInfo => Boolean(d));

		if (depFiles.length > 0) {
			const depNames = props.depends.map((d) => d?.Name ?? "dependency");
			await DownloadsMods(depNames, depFiles);
		}

		const filename = getMinecraftDownloadFileName(mod.Name, mod.Versions);
		await DownloadFileToMinecraftMods(detail.URL, filename);
		await saveModToYaml(mod, "downloads");
		await showNotify("Успех", `Мод "${mod.Name}" успешно загружен!`);
	} catch (err) {
		console.error("Ошибка при скачивании мода:", err);
		await showNotify("Ошибка", `Не удалось скачать мод "${mod.Name}".`);
	} finally {
		isDownloading.value = false;
	}
};
</script>

<template>
  <div class="mod-download">
    <button
      v-for="detail in mod?.Details || []"
      :key="`${mod.Name}-${detail.Version}-${detail.Loader}`"
      class="button"
      :disabled="isDownloading"
      @click="downloadMod(mod, detail)"
      style="background-color: green; width: 50%; margin-bottom: 14px;"
    >
      <span v-if="!isDownloading">
        Скачать {{ detail.Version }} | {{ detail.Loader }} |
        Скачано {{ detail.Downloads ?? 0 }} раз
      </span>
      <span v-else>Загрузка...</span>
    </button>
  </div>
</template>

<style lang="css" scoped>
.mod-download{
  display: flex;
  flex-direction: column;

  justify-content: center;
  align-items: center;
  
}
</style>