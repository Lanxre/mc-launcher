<script lang="ts" setup>
import {
	SortByLoader,
	SortByVersions,
} from "@wailsjs/go/functools/FuncService";
import { GetMods, GetModsByPage } from "@wailsjs/go/parser/ScraperService";
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { uniqueBy } from "@/api/utils";
import ModLoader from "@/components/Mods/ModLoader.vue";
import ModsFilter from "@/components/Mods/ModsFilter.vue";
import ModsList from "@/components/Mods/ModsList.vue";
import { useScrollManager } from "@/composables/useScrollManager";
import { useModStore } from "@/stores/modStore";
import type { MinecraftMod } from "@/types";

const route = useRoute();
const router = useRouter();
const modStore = useModStore();
const { restoreMainPageScroll } = useScrollManager();

const mods = ref<MinecraftMod[]>([]);
const allMods = ref<MinecraftMod[]>([]);
const versionE = ref("");
const loaderE = ref("");
const searchE = ref("");

const currentPage = ref(1);
const searchPage = ref(1);
const loadingMore = ref(false);
const hasMore = ref(true);

const loaderTrigger = ref<any>(null);
const mainElement = ref<HTMLElement | null>(null);
let observer: IntersectionObserver | null = null;

const getUniqueVersions = computed(() =>
	[...new Set(allMods.value.flatMap((m) => m.Versions || []))].sort(),
);

const getUniqueLoaders = computed(() =>
	[
		...new Set(allMods.value.flatMap((m) => m.Loaders || []).filter(Boolean)),
	].sort(),
);

const resetPaging = () => {
	currentPage.value = 1;
	searchPage.value = 1;
	hasMore.value = true;
};

const updateMods = (list: MinecraftMod[]) => {
	mods.value = list;
	allMods.value = [...list];
	hasMore.value = list.length > 0;
};

const loadInitialMods = async () => {
	if (modStore.getAllMods.length > 0) {
		allMods.value = modStore.getAllMods;
		mods.value = [...allMods.value];
		versionE.value = modStore.getVersionFilter ?? "";
		loaderE.value = modStore.getLoaderFilter ?? "";
		searchE.value = modStore.getSearchFilter ?? "";

		if (versionE.value || loaderE.value) await applyFilters();
		if (searchE.value) await applySearch();
	} else {
		const result = await GetMods();
		updateMods(uniqueBy(result, (m) => m.Name));
		modStore.addMods(allMods.value);
		modStore.setCurrentParsePage(1);
	}
	await nextTick();
	checkAutoLoad();
};

const loadMoreMods = async () => {
	if (loadingMore.value || !hasMore.value) return;
	loadingMore.value = true;

	try {
		const nextPage = searchE.value
			? ++searchPage.value
			: (modStore.getParsePage ?? currentPage.value) + 1;

		const result = await GetModsByPage(nextPage, searchE.value || null);
		if (!result?.length) return (hasMore.value = false);

		const existingUrls = new Set(mods.value.map((m) => m.ModPageLink));
		const newMods = uniqueBy(
			result.filter((m) => !existingUrls.has(m.ModPageLink)),
			(m) => m.Name,
		);

		if (!newMods.length) return (hasMore.value = false);

		mods.value.push(...newMods);
		allMods.value.push(...newMods);
		modStore.setCurrentParsePage(nextPage);

		if (versionE.value || loaderE.value) await applyFilters();
		await nextTick();
		reinitObserver();
	} finally {
		loadingMore.value = false;
	}
};

const applyFilters = async () => {
	let filtered = [...allMods.value];
	if (versionE.value) {
		filtered = await SortByVersions(filtered, versionE.value);
		modStore.setVersionFilter(versionE.value);
	}
	if (loaderE.value) {
		filtered = await SortByLoader(filtered, loaderE.value);
		modStore.setLoaderFilter(loaderE.value);
	}
	mods.value = filtered;
};

const applySearch = async () => {
	resetPaging();
	if (!searchE.value) {
		updateMods(
			modStore.getAllMods.length ? modStore.getAllMods : allMods.value,
		);
		modStore.setSearchFilter("");
	} else {
		const searched = await GetModsByPage(1, searchE.value);
		updateMods(searched);
		modStore.setSearchFilter(searchE.value);
	}

	if (versionE.value || loaderE.value) await applyFilters();
	await nextTick();
	reinitObserver();
};

const initObserver = async (el: HTMLElement) => {
	observer?.disconnect();
	observer = new IntersectionObserver(
		([entry]) => {
			if (entry?.isIntersecting && hasMore.value && !loadingMore.value)
				loadMoreMods();
		},
		{ rootMargin: "100px", threshold: 0.1 },
	);
	observer.observe(el);
};

const reinitObserver = async () => {
	if (!loaderTrigger.value?.loaderTrigger) return;
	observer?.disconnect();
	await nextTick();
	observer?.observe(loaderTrigger.value.loaderTrigger);
};

const checkAutoLoad = () => {
	if (!mainElement.value || !hasMore.value || loadingMore.value) return;
	const containerHeight = mainElement.value.getBoundingClientRect().height;
	if (containerHeight < window.innerHeight - 100) loadMoreMods();
};

const goToDownloads = () => {
	modStore.setAllMods(mods.value);
	modStore.setCurrentParsePage(currentPage.value);
	router.push({ name: "mod-downloads" });
};

const goToFavourites = () => {
	router.push({ name: "mod-favourites" });
};

watch(
	loaderTrigger,
	async (comp) => {
		if (comp?.loaderTrigger && hasMore.value) {
			await nextTick();
			initObserver(comp.loaderTrigger);
		}
	},
	{ flush: "post" },
);

watch(searchE, () => {
	if (!searchE.value) searchPage.value = 1;
});

onMounted(async () => {
	if (route.path === "/") restoreMainPageScroll();
	await loadInitialMods();
	window.addEventListener("resize", checkAutoLoad);
});

onUnmounted(() => {
	observer?.disconnect();
	window.removeEventListener("resize", checkAutoLoad);
});
</script>

<template>
  <div class="main" ref="mainElement">
    <ModsFilter
		v-model:searchE="searchE"
		v-model:versionE="versionE"
		v-model:loaderE="loaderE"
		
		:get-unique-loaders="getUniqueLoaders"
		:get-unique-versions="getUniqueVersions"

		:go-to-downloads="goToDownloads"
		:go-to-favourites="goToFavourites"
		
		:apply-filters="applyFilters"
		:apply-search="applySearch"
	/>
    <ModsList :mods="mods" :loaderf="loaderE" :versionf="versionE" :searchf="searchE" />

    <ModLoader
      ref="loaderTrigger"
      :has-more="hasMore"
      :loading-more="loadingMore"
      :count="mods.length"
      loading-text="Загрузка модов..."
      end-message="Все моды загружены"
      @load-more="loadMoreMods"
    />
  </div>
</template>


<style lang="css" scoped>
.main {
  display: flex;
  flex-direction: column;
  width: 100%;
  flex: 1;
  padding: 20px;
  gap: 15px;

  background-color: var(--main-bg-color);
}
</style>