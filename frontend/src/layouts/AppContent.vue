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

const displayedMods = ref<MinecraftMod[]>([]); 
const cachedMods = ref<MinecraftMod[]>([]);

const selectedVersion = ref("");
const selectedLoader = ref("");
const searchQuery = ref("");

const currentPage = ref(1);
const searchPage = ref(1);
const loadingMore = ref(false);
const hasMore = ref(true);

const loaderTriggerRef = ref<InstanceType<typeof ModLoader> | null>(null);
const mainContainerRef = ref<HTMLElement | null>(null);
let observer: IntersectionObserver | null = null;

const uniqueVersions = computed(() =>
  [...new Set(cachedMods.value.flatMap((m) => m.Versions || []))].sort()
);

const uniqueLoaders = computed(() =>
  [...new Set(cachedMods.value.flatMap((m) => m.Loaders || []).filter(Boolean))].sort()
);

const resetPagination = () => {
  currentPage.value = 1;
  searchPage.value = 1;
  hasMore.value = true;
};

const applyLocalFilters = async () => {
  let result = [...cachedMods.value];

  if (selectedVersion.value) {
    result = await SortByVersions(result, selectedVersion.value);
    modStore.setVersionFilter(selectedVersion.value);
  }
  
  if (selectedLoader.value) {
    result = await SortByLoader(result, selectedLoader.value);
    modStore.setLoaderFilter(selectedLoader.value);
  }

  displayedMods.value = result;
};

const loadMoreMods = async () => {
  if (loadingMore.value || !hasMore.value) return;
  loadingMore.value = true;

  try {
    const isSearch = !!searchQuery.value;
    const nextPage = isSearch 
      ? searchPage.value + 1 
      : (modStore.getParsePage ?? currentPage.value) + 1;

    const fetched = await GetModsByPage(nextPage, searchQuery.value || null);

    if (!fetched || fetched.length === 0) {
      hasMore.value = false;
      return;
    }

    const existingLinks = new Set(cachedMods.value.map((m) => m.ModPageLink));
    const uniqueNewMods = uniqueBy(
      fetched.filter((m) => !existingLinks.has(m.ModPageLink)),
      (m) => m.Name
    );

    if (uniqueNewMods.length === 0) {
      hasMore.value = false;
      return;
    }

    cachedMods.value.push(...uniqueNewMods);
    if (!isSearch) {
      modStore.setCurrentParsePage(nextPage);
    }
    
    if (isSearch) searchPage.value = nextPage;
    else currentPage.value = nextPage;

    await applyLocalFilters();
  } finally {
    loadingMore.value = false;
  }
};


const initData = async () => {
  if (modStore.getAllMods.length > 0) {
    cachedMods.value = modStore.getAllMods;
    selectedVersion.value = modStore.getVersionFilter ?? "";
    selectedLoader.value = modStore.getLoaderFilter ?? "";
    searchQuery.value = modStore.getSearchFilter ?? "";
    
    if (searchQuery.value) {
      await handleSearch(searchQuery.value, false);
    } else {
      await applyLocalFilters();
    }
  } 
  else {
    const result = await GetMods();
    cachedMods.value = uniqueBy(result, (m) => m.Name);
    displayedMods.value = [...cachedMods.value];
    modStore.addMods(cachedMods.value);
    modStore.setCurrentParsePage(1);
  }

  await nextTick();

  if (route.path === "/") {
    restoreMainPageScroll();
  }

  initObserver();
  checkAutoLoad();
};

let searchTimeout: ReturnType<typeof setTimeout>;

const onSearchInput = (query: string) => {
  clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => handleSearch(query), 400);
};

const handleSearch = async (query: string, updateStore = true) => {
  resetPagination();
  searchQuery.value = query;
  
  if (updateStore) modStore.setSearchFilter(query);

  if (!query) {
    cachedMods.value = modStore.getAllMods.length ? modStore.getAllMods : cachedMods.value;
  } else {
	const searched = await GetModsByPage(1, query);
    cachedMods.value = searched;
  }

  await applyLocalFilters();
  await nextTick();
  checkAutoLoad();
};

const goToDownloads = () => {
  modStore.setAllMods(displayedMods.value);
  modStore.setCurrentParsePage(currentPage.value);
  router.push({ name: "mod-downloads" });
};

const goToFavourites = () => {
  router.push({ name: "mod-favourites" });
};

const initObserver = () => {
  observer?.disconnect();
  
  const target = loaderTriggerRef.value?.$el || loaderTriggerRef.value?.loaderTrigger; // Поддержка доступа к DOM элементу компонента
  
  if (!target) return;

  observer = new IntersectionObserver(
    ([entry]) => {
      if (entry && entry.isIntersecting && hasMore.value && !loadingMore.value) {
        loadMoreMods();
      }
    },
    { rootMargin: "200px", threshold: 0.1 }
  );
  
  observer.observe(target);
};

const checkAutoLoad = () => {
  if (!mainContainerRef.value || !hasMore.value || loadingMore.value) return;
  const containerHeight = mainContainerRef.value.getBoundingClientRect().height;
  if (containerHeight < window.innerHeight) {
    loadMoreMods();
  }
};

watch([selectedVersion, selectedLoader], async () => {
  await applyLocalFilters();
  await nextTick();
  checkAutoLoad();
});

watch(() => loaderTriggerRef.value, () => {
  initObserver();
});

onMounted(async () => {
  await initData();
  window.addEventListener("resize", checkAutoLoad);
});

onUnmounted(() => {
  observer?.disconnect();
  window.removeEventListener("resize", checkAutoLoad);
});
</script>

<template>
  <div class="main" ref="mainContainerRef">
    <ModsFilter
      :searchE="searchQuery"
      @update:searchE="onSearchInput"
      
      v-model:versionE="selectedVersion"
      v-model:loaderE="selectedLoader"
      
      :get-unique-loaders="uniqueLoaders"
      :get-unique-versions="uniqueVersions"

      :go-to-downloads="goToDownloads"
      :go-to-favourites="goToFavourites"
      
      :apply-filters="applyLocalFilters"
      :apply-search="() => handleSearch(searchQuery)"
    />
    
    <ModsList 
      :mods="displayedMods" 
      :loaderf="selectedLoader" 
      :versionf="selectedVersion" 
      :searchf="searchQuery" 
    />

    <ModLoader
      ref="loaderTriggerRef"
      :has-more="hasMore"
      :loading-more="loadingMore"
      :count="displayedMods.length"
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