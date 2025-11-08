<script setup lang="ts">
import { ref } from 'vue'
import { DownloadFileToMinecraftMods, DownloadsMods } from '@/../wailsjs/go/main/FileService'
import { ShowInfoMessage } from '@/../wailsjs/go/main/App'
import type { MinecraftMod, DownloadInfo, ModDependency } from '@/types'
import { saveModToYaml } from '@/api/utils'

interface Props {
  mod: MinecraftMod
  depends: ModDependency[]
}

const props = defineProps<Props>()

const isDownloading = ref(false)

const showNotify = async (title: string, message: string) => {
  await ShowInfoMessage(title, message)
}

const downloadMod = async (mod: MinecraftMod, detail: DownloadInfo) => {
  if (isDownloading.value) return
  isDownloading.value = true

  try {
    if (props.depends.length > 0) {
      const names = props.depends.map(d => d.Name)
      console.log(props.depends)
      const depFiles = props.depends.flatMap(d =>
        d.Details.filter(dl => dl.Version.includes(detail.Version) && dl.Loader === detail.Loader)
      )

      if (depFiles.length) {
        await DownloadsMods(names, depFiles)
      }
    }

    await DownloadFileToMinecraftMods(detail.URL, mod.Name, detail.Version)
    await saveModToYaml(mod, 'downloads')
    await showNotify('Успех', `Мод "${mod.Name}" успешно загружен!`)
  } catch (err: any) {
    console.error('Ошибка при скачивании мода:', err)
    await showNotify('Ошибка', `Не удалось скачать мод "${mod.Name}".`)
  } finally {
    isDownloading.value = false
  }
}
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