<script lang="ts" setup>
import {
	GetModDepends,
	GetModDetails,
} from "@wailsjs/go/parser/ScraperService";
import { onMounted, ref } from "vue";
import { enrichDependencies, saveModToYaml } from "@/api/utils";
import ModDescription from "@/components/ModDetails/ModDescription.vue";
import ModDownloads from "@/components/ModDetails/ModDownloads.vue";
import ModScreenshots from "@/components/ModDetails/ModScreenshots.vue";
import ModTitle from "@/components/ModDetails/ModTitle.vue";
import ModVersions from "@/components/ModDetails/ModVersions.vue";
import AppFooter from "@/layouts/AppFooter.vue";
import AppHeader from "@/layouts/AppHeader.vue";
import { useModStore } from "@/stores/modStore";
import type { MinecraftMod, ModDependency } from "@/types";

const modStore = useModStore();
const mod = ref<MinecraftMod | null>(null);
const depends = ref<ModDependency[]>([]);
const isLoading = ref(true);
const isError = ref(false);

async function fetchDependencies(
	initialDeps: ModDependency[],
): Promise<ModDependency[]> {
	const visited = new Set();
	const allDeps: ModDependency[] = [];

	async function processDeps(deps: ModDependency[]) {
		if (deps == null) return;
		const uniqueDeps = deps.filter(
			(d) => d.ModPageLink && !visited.has(d.ModPageLink),
		);
		if (uniqueDeps.length === 0) return;

		uniqueDeps.forEach((d) => visited.add(d.ModPageLink));

		const newDeps = await GetModDepends(uniqueDeps);
		allDeps.push(...newDeps);

		const nestedDeps = newDeps
			.map((d) => d.Dependency)
			.filter(Boolean)
			.flat()
			.filter((d) => d.ModPageLink && !visited.has(d.ModPageLink));

		if (nestedDeps.length > 0) {
			await processDeps(nestedDeps);
		}
	}

	await processDeps(initialDeps);
	return allDeps;
}

function applyFilters(modObj: MinecraftMod) {
	const loaderFilter = modStore.getLoaderFilter;
	const versionFilter = modStore.getVersionFilter;

	if (loaderFilter) {
		modObj.Details = modObj.Details.filter((d) => d.Loader === loaderFilter);
	}

	if (versionFilter) {
		modObj.Details = modObj.Details.filter((d) => d.Version === versionFilter);
	}
}

async function loadModDetails() {
	try {
		const currentMod = modStore.currentMod;
		if (!currentMod?.ModPageLink) return;

		const fullDetails = await GetModDetails(currentMod.ModPageLink);
		currentMod.Screenshots = fullDetails.Screenshots;
		currentMod.Details = fullDetails.Details;
		currentMod.Dependency = fullDetails.Dependency;

		applyFilters(currentMod);
		let allDeps = await fetchDependencies(currentMod.Dependency);
		allDeps = allDeps.filter(
			(dep, i, self) =>
				dep.ModPageLink &&
				i === self.findIndex((d) => d.ModPageLink === dep.ModPageLink),
		);

		await enrichDependencies(allDeps);
		depends.value = allDeps;
		mod.value = currentMod;
		mod.value.Dependency.forEach((dep) => {
			const as = depends.value.find((d) => d.Name === dep.Name);
			if (as != undefined) {
				dep.Details = as.Details;
			}
		});
		console.log(mod.value);
	} catch (err) {
		console.error("Ошибка при загрузке данных мода:", err);
		isError.value = true;
	} finally {
		isLoading.value = false;
	}
}

onMounted(loadModDetails);
</script>


<template>
  <AppHeader />

  <div class="mod" v-if="mod">
    <ModScreenshots v-if="mod && mod.Screenshots !== null" :screenshots="mod.Screenshots" />

    <div class="mod-content">
      <ModTitle :name="mod.Name" :link-page="mod.ModPageLink" @save-favourite="saveModToYaml(mod, 'favourite')"/>
      <ModVersions :versions="mod.Versions" />
      <ModDescription :description="mod.Description" />
      <ModDownloads :mod="mod" :depends="depends" />
    </div>
  </div>
  <div v-else-if="isLoading" class="loading-state">
    <p>Загрузка информации о моде...</p>
  </div>

  <div v-else-if="isError" class="error-state">
    <p>Не удалось загрузить данные. Попробуйте позже.</p>
  </div>
  <AppFooter />
</template>

<style scoped>
.mod {
  display: flex;
  flex-direction: column;
  gap: 24px;
  padding: 20px;
}

.mod-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.loading-state,
.error-state {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 60vh;
  font-size: 1.2rem;
  color: #888;
}
</style>
