<script lang="ts" setup>
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { 
  SortByLoader, 
  SortByVersions 
} from "@wailsjs/go/functools/FuncService";
import { 
  GetMods, 
  GetModsByPage 
} from "@wailsjs/go/parser/ScraperService";

import { uniqueBy } from "@/api/utils";
import { useScrollManager } from "@/composables/useScrollManager";
import { useModStore } from "@/stores/modStore";
import type { MinecraftMod } from "@/types";

import ModLoader from "@/components/Mods/ModLoader.vue";
import ModsFilter from "@/components/Mods/ModsFilter.vue";
import ModsList from "@/components/Mods/ModsList.vue";

const route = useRoute();
const router = useRouter();
const modStore = useModStore();
const { restoreMainPageScroll } = useScrollManager();

const displayedMods = ref<MinecraftMod[]>([]); 
const cachedMods = ref<MinecraftMod[]>([]);

const selectedVersion = ref(modStore.getVersionFilter ?? "");
const selectedLoader = ref(modStore.getLoaderFilter ?? "");
const searchQuery = ref(modStore.getSearchFilter ?? "");

const page = ref(1);
const loading = ref(false);
const hasMore = ref(true);

const mainContainerRef = ref<HTMLElement | null>(null);
const loaderTriggerRef = ref<HTMLElement | null>(null);
let observer: IntersectionObserver | null = null;
let searchTimeout: ReturnType<typeof setTimeout>;

const uniqueVersions = computed(() =>
  [...new Set(cachedMods.value.flatMap(m => m.Versions || []))].sort()
);

const uniqueLoaders = computed(() =>
  [...new Set(cachedMods.value.flatMap(m => m.Loaders || []).filter(Boolean))].sort()
);

const applyFilters = async () => {
  let result = [...cachedMods.value];

  if (selectedVersion.value) {
    result = await SortByVersions(result, selectedVersion.value);
  }
  
  if (selectedLoader.value) {
    result = await SortByLoader(result, selectedLoader.value);
  }

  displayedMods.value = result;
  
  modStore.setVersionFilter(selectedVersion.value);
  modStore.setLoaderFilter(selectedLoader.value);
};

const fetchMods = async (isNextPage = false) => {
  if (loading.value || (!hasMore.value && isNextPage)) return;
  
  loading.value = true;
  const targetPage = isNextPage ? page.value + 1 : 1;

  try {
    const fetched = await GetModsByPage(targetPage, searchQuery.value || null);

    if (!fetched || fetched.length === 0) {
      hasMore.value = false;
      return;
    }

    const newMods = uniqueBy(fetched, (m) => m.Name);
    
    if (isNextPage) {
      const existingLinks = new Set(cachedMods.value.map(m => m.ModPageLink));
      const filteredNew = newMods.filter(m => !existingLinks.has(m.ModPageLink));
      cachedMods.value.push(...filteredNew);
    } else {
      cachedMods.value = newMods;
      hasMore.value = true;
    }

    page.value = targetPage;
    if (!searchQuery.value) modStore.setCurrentParsePage(targetPage);
    
    await applyFilters();
  } finally {
    loading.value = false;
    await nextTick();
    checkAutoLoad();
  }
};

const handleSearch = (query: string) => {
  clearTimeout(searchTimeout);
  searchQuery.value = query;
  modStore.setSearchFilter(query);
  searchTimeout = setTimeout(() => fetchMods(false), 400);
};

const initData = async () => {
  if (modStore.getAllMods.length > 0) {
    cachedMods.value = modStore.getAllMods;
    page.value = searchQuery.value ? 1 : (modStore.getParsePage || 1);
    await applyFilters();
  } else {
    const initial = await GetMods();
    cachedMods.value = uniqueBy(initial, m => m.Name);
    modStore.addMods(cachedMods.value);
    modStore.setCurrentParsePage(1);
    displayedMods.value = [...cachedMods.value];
  }
  initObserver();
};

const initObserver = () => {
  observer?.disconnect();
  if (!loaderTriggerRef.value) return;

  observer = new IntersectionObserver(([entry]) => {
    if (entry?.isIntersecting && hasMore.value && !loading.value) {
      fetchMods(true);
    }
  }, { rootMargin: "200px" });
  
  observer.observe(loaderTriggerRef.value);
};

const checkAutoLoad = () => {
  if (!mainContainerRef.value || !hasMore.value || loading.value) return;
  if (mainContainerRef.value.getBoundingClientRect().height < window.innerHeight) {
    fetchMods(true);
  }
};

const goToDownloads = () => {
  modStore.setAllMods(displayedMods.value);
  router.push({ name: "mod-downloads" });
};

watch([selectedVersion, selectedLoader], applyFilters);

onMounted(() => {
  initData();
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
      v-model:versionE="selectedVersion"
      v-model:loaderE="selectedLoader"
      :get-unique-loaders="uniqueLoaders"
      :get-unique-versions="uniqueVersions"
      :go-to-downloads="goToDownloads"
      :go-to-favourites="() => router.push({ name: 'mod-favourites' })"
      @update:searchE="handleSearch"
      @apply-filters="applyFilters"
    />
    
    <ModsList 
      :mods="displayedMods" 
      :loaderf="selectedLoader" 
      :versionf="selectedVersion" 
      :searchf="searchQuery" 
    />

    <div ref="loaderTriggerRef" class="loader-anchor">
      <ModLoader
        :has-more="hasMore"
        :loading-more="loading"
        :count="displayedMods.length"
        @load-more="fetchMods(true)"
      />
    </div>
  </div>
</template>

<style scoped>
.main {
  display: flex;
  flex-direction: column;
  width: 100%;
  flex: 1;
  padding: 20px;
  gap: 15px;
  background-color: var(--main-bg-color);
}
.loader-anchor {
  min-height: 20px;
}
</style>