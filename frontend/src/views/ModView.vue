<script lang="ts" setup>
import { onMounted, ref } from 'vue';
import { useModStore } from '@/stores/modStore';
import { GetModDetails } from "../../wailsjs/go/main/ScraperService"
import { DownloadFileToMinecraftMods } from "../../wailsjs/go/main/FileService"
import { ShowInfoMessage } from "../../wailsjs/go/main/App"


import AppHeader from '@/layouts/AppHeader.vue';
import AppFooter from '@/layouts/AppFooter.vue';
import Slider from '@/components/Slider/Slider.vue';

import type { MinecraftMod } from '@/types';
import { openLink } from '@/api/utils';
import LinkIcon from '@/assets/images/link.png'


const modStore = useModStore()

const mod = ref<MinecraftMod | null>(null)

const downloadMod = async (url: string, modName: string, modVersion: string) => {
  await DownloadFileToMinecraftMods(url, modName, modVersion);
  await showNotify()
}

const showNotify = async () => {
    await ShowInfoMessage("Успех", "Мод успешно сохранён!");
};

onMounted(async () => {
    mod.value = modStore.currentMod
    if (mod.value?.ModPageLink) {
      const customMod = await GetModDetails(mod.value.ModPageLink)
      mod.value.Screenshots = customMod.Screenshots
      mod.value.Details = customMod.Details
      mod.value.Dependency = customMod.Dependency
      console.log(mod.value)

      const loaderf = modStore.getLoaderFilter
      const versionf = modStore.getVersionFilter

      if (loaderf != null && loaderf != "") {
        mod.value.Details = mod.value.Details.filter((d) => d.Loader == loaderf)
      }

      if (versionf != null && versionf != "") {
        mod.value.Details = mod.value.Details.filter((d) => d.Version == versionf)
      }

    }
})

</script>

<template>
  <AppHeader/>

  <div class="mod">
        <div class="screenshots">
            <Slider v-if="mod?.Screenshots && mod.Screenshots.length > 0" :slides="mod?.Screenshots.map((s) => {
              return {
                image: s,
                alt: 'no screen',
              }
            })"
            />
            <div v-else class="no-screenshots">
              <p class="text text-shd">Скриншоты не найдены</p>
            </div>
          </div>
        <div style="height: 100%;">
          <div class="mod-content">
            <div class="mod-name" @click="openLink(mod?.ModPageLink!)">
              <h1>{{ mod?.Name }} </h1>
              <img :src="LinkIcon" alt="modlink" class="icon"/>
            </div>
            <p v-if="mod?.Versions && mod.Versions.length === 1">Версия: {{ mod?.Versions[0] }}</p>
            <p v-else>Версии: {{ mod?.Versions.join(', ') }}</p>
            <br/>
            <span style="margin-bottom: 3px;"> {{ mod?.Description }}</span>
          </div>
          <div class="mod-download">
              <button v-if="mod?.Details && mod.Details.length >= 1" v-for="(detail) in mod?.Details"  class="button" @click="downloadMod(detail.URL, mod.Name, detail.Version)" style="background-color: green; width: 50%; margin-bottom: 14px;">
                Скачать {{ detail.Version }} | {{ detail.Loader }} | Скачано {{ detail.Downloads ? detail.Downloads : '0' }} раз
              </button>
          </div>

        </div>
  </div>

  <AppFooter/>
</template>

<style lang="css">
.mod {
  display: flex;
  flex-direction: column;
  
  width: 100%;
  
  padding: 20px;
  
  gap: 15px;
}

.mod-content {
  display: flex;
  flex-direction: column;

}

.no-screenshots {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.mod-download{
  display: flex;
  flex-direction: column;

  justify-content: center;
  align-items: center;
  
}

.mod-name {
  display: flex;
  flex-direction: row;

  gap: 5px;
  cursor: pointer;
}

.mod-name img {
  width: 15px;
  height: 15px;
  filter: brightness(0) invert(1);
}

</style>