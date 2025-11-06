<script lang="ts" setup>
import { onMounted, onUnmounted, ref, nextTick, computed, watch } from 'vue'
import { useModStore } from '@/stores/modStore'
import { GetMods, GetModsByPage } from '../../wailsjs/go/main/ScraperService'
import { SortByVersions, SortByLoader } from '../../wailsjs/go/main/FuncService'
import type { MinecraftMod } from '@/types'
import { uniqueBy } from '@/api/utils'

import ModsList from '@/components/Mods/ModsList.vue'
import ModLoader from '@/components/Mods/ModLoader.vue'
import SearchFilter from '@/components/SearchFilter/SearchFilter.vue'
import InputSearch from '@/components/InputSearch/InputSearch.vue'
import Button from '@/components/Buttons/Button.vue'
import { useScrollManager } from '@/composables/useScrollManager'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const modStore = useModStore()
const { restoreMainPageScroll } = useScrollManager()

const mods = ref<MinecraftMod[]>([])
const allMods = ref<MinecraftMod[]>([])

const loaderE = ref('')
const versionE = ref('')
const searchE = ref('')

const currentPage = ref(1)
const searchPage = ref(1)

const loadingMore = ref(false)
const hasMore = ref(true)

let observer: IntersectionObserver | null = null

const loaderTrigger = ref<any>(null)
const mainElement = ref<HTMLElement | null>(null)

const loadMods = async () => {
  if (modStore.getAllMods.length > 0) {
    mods.value = modStore.getAllMods
    allMods.value = modStore.getAllMods
    hasMore.value = true

    const storeVersionFilter = modStore.getVersionFilter
    const storeLoaderFilter = modStore.getLoaderFilter

    if (storeVersionFilter || storeLoaderFilter) {
      versionE.value = storeVersionFilter ?? ''
      loaderE.value = storeLoaderFilter ?? ''
      await handleFilter(versionE.value, loaderE.value)
    }

    const storeSearchFilter = modStore.getSearchFilter
    if (storeSearchFilter) {
      searchE.value = storeSearchFilter
      await handleSearchFilter(searchE.value)
    }
  } else {
    const result = await GetMods()
    mods.value = result as MinecraftMod[]
    allMods.value = [...uniqueBy(mods.value, m => m.Name) as MinecraftMod[]]
    hasMore.value = result.length > 0
    modStore.addMods(allMods.value)
    modStore.setCurrentParsePage(1)
  }

  await nextTick()
  checkIfNeedMoreContent()
}

const loadMoreMods = async () => {
  if (loadingMore.value || !hasMore.value) return

  loadingMore.value = true

  try {
    let result: MinecraftMod[] = []

    if (searchE.value) {
      const nextPage = searchPage.value + 1
      result = await GetModsByPage(nextPage, searchE.value)
      searchPage.value = nextPage
    } else {
      const nextPage = modStore.getParsePage !== null ? modStore.getParsePage + 1 : currentPage.value + 1
      result = await GetModsByPage(nextPage, null)
      currentPage.value = nextPage
      if (modStore.getParsePage !== null) {
        modStore.setCurrentParsePage(nextPage)
      }
    }

    if (!result || result.length === 0) {
      hasMore.value = false
      return
    }

    const existingUrls = new Set(mods.value.map(mod => mod.ModPageLink))
    const newMods = result.filter(mod => !existingUrls.has(mod.ModPageLink))

    if (newMods.length > 0) {
      const uniqueNewMods = uniqueBy(newMods, m => m.Name)
      mods.value.push(...uniqueNewMods)
      allMods.value.push(...uniqueNewMods)

      if (versionE.value || loaderE.value) {
        await handleFilter(versionE.value, loaderE.value)
      }

      await nextTick()
      reinitObserver()
    } else {
      hasMore.value = false
    }
  } catch {
    hasMore.value = false
  } finally {
    loadingMore.value = false
  }
}

const reinitObserver = async () => {
  if (!loaderTrigger.value?.loaderTrigger) return
  observer?.disconnect()
  await nextTick()
  observer?.observe(loaderTrigger.value.loaderTrigger)
}

const checkIfNeedMoreContent = () => {
  if (!mainElement.value || !hasMore.value || loadingMore.value) return

  const containerRect = mainElement.value.getBoundingClientRect()
  const availableHeight = window.innerHeight - 100

  if (containerRect.height < availableHeight) {
    loadMoreMods()
  }
}

const handleFilter = async (version: string, loader: string) => {
  let filteredMods = [...allMods.value]

  if (version) {
    versionE.value = version
    filteredMods = await SortByVersions(filteredMods, version)
    modStore.setVersionFilter(version)
  }

  if (loader) {
    loaderE.value = loader
    filteredMods = await SortByLoader(filteredMods, loader)
    modStore.setLoaderFilter(loader)
  }

  mods.value = filteredMods
}

const initObserver = async (element: HTMLElement) => {
  observer?.disconnect()
  observer = new IntersectionObserver(
    ([entry]) => {
      if (entry != undefined && entry.isIntersecting && hasMore.value && !loadingMore.value) {
        loadMoreMods()
      }
    },
    { root: null, rootMargin: '100px', threshold: 0.1 }
  )
  observer.observe(element)
}

const getUniqueVersions = computed(() => {
  const allVersions = allMods.value.flatMap(mod => mod.Versions || [])
  return [...new Set(allVersions)].sort()
})

const getUniqueLoaders = computed(() => {
  const allLoaders = allMods.value.flatMap(mod => mod.Loaders)
  return [...new Set(allLoaders)].filter(l => l !== '').sort()
})

const handleSearchFilter = async (inputSearch: string) => {
  searchE.value = inputSearch
  searchPage.value = 1

  if (!inputSearch) {
    mods.value = [...allMods.value]
    hasMore.value = true
    currentPage.value = 1
    modStore.setCurrentParsePage(1)
    modStore.setSearchFilter('')
  } else {
    const searchedMods = await GetModsByPage(1, inputSearch)
    modStore.setSearchFilter(inputSearch)
    mods.value = searchedMods as MinecraftMod[]
    allMods.value = [...searchedMods]
    hasMore.value = searchedMods.length > 0
    currentPage.value = 1
  }

  if (versionE.value || loaderE.value) {
    await handleFilter(versionE.value, loaderE.value)
  }

  await nextTick()
  reinitObserver()
}

const redirToDowloadMods = () => {
  modStore.setAllMods(mods.value)
  modStore.setCurrentParsePage(currentPage.value)
  router.push({ name: 'mod-downloads' })
}

watch(searchE, (newVal) => {
  if (!newVal) searchPage.value = 1
})

watch(loaderTrigger, async (comp) => {
  if (comp?.loaderTrigger && hasMore.value) {
    await nextTick()
    initObserver(comp.loaderTrigger)
  }
}, { flush: 'post' })

onMounted(async () => {
  if (route.path === '/') {
    restoreMainPageScroll()
  }

  await loadMods()
  await nextTick()

  window.addEventListener('resize', checkIfNeedMoreContent)
})

onUnmounted(() => {
  observer?.disconnect()
  window.removeEventListener('resize', checkIfNeedMoreContent)
})
</script>

<template>
  <div class="main" ref="mainElement">
    <div class="filters">
      <Button name="Список установленных модов" @click="redirToDowloadMods" />
      <InputSearch v-model="searchE" @search="handleSearchFilter" />
      <SearchFilter :versions="getUniqueVersions" :loaders="getUniqueLoaders" @filter="handleFilter" />
    </div>
    <ModsList :mods="mods" :loaderf="loaderE" :versionf="versionE" />
    <ModLoader
      ref="loaderTrigger"
      :has-more="hasMore"
      :loading-more="loadingMore"
      :count="mods.length"
      loading-text="Загрузка модов..."
      end-message="Все моды загружены"
      :auto-load="false"
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
}

.filters {
  display: flex;
  flex-direction: row;
  justify-content: flex-end;
  gap: 5px;
}
</style>