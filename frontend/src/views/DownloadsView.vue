<script setup lang="ts">
import { onMounted, ref } from 'vue'

import AppFooter from '@/layouts/AppFooter.vue'
import AppHeader from '@/layouts/AppHeader.vue'
import View from '@/components/View/View.vue'
import ModsList from '@/components/Mods/ModsList.vue'

import { GetYamlConfig, RemoveFromYamlConfig } from '../../wailsjs/go/main/FuncService'
import { ShowInfoMessage } from '../../wailsjs/go/main/App'
import type { MinecraftMod } from '@/types'

const savedMods = ref<MinecraftMod[]>([])

const loadDownloadedMods = async () => {
  try {
    const mods = await GetYamlConfig('downloads')
    savedMods.value = mods ?? []
  } catch (err) {
    console.error('Ошибка при получении скачанных модов:', err)
  }
}

const removeFromDownloads = async (mod: MinecraftMod) => {
  try {
    savedMods.value = savedMods.value.filter(m => m.Name !== mod.Name)
    await RemoveFromYamlConfig(mod, 'downloads')
    await ShowInfoMessage('Удалён', `Мод "${mod.Name}" успешно удалён`)
  } catch (err) {
    console.error('Ошибка при удалении мода:', err)
  }
}

onMounted(loadDownloadedMods)
</script>

<template>
  <AppHeader />
  <View min-height="100vh">
    <div class="downloads-list">
      <ModsList
        v-if="savedMods.length > 0"
        :mods="savedMods"
        :onDelete="removeFromDownloads"
      />
      <p v-else class="no-mods">Нет скачанных модов</p>
    </div>
  </View>
  <AppFooter />
</template>

<style scoped>
.downloads-list {
  display: flex;
  flex-direction: column;
  padding: 20px;
  min-height: 100vh;
  box-sizing: border-box;
}

.no-mods {
  text-align: center;
  color: #ccc;
  font-size: 18px;
  margin-top: 40px;
}
</style>
