<script setup lang="ts">
import {
	DeleteSavedMod,
	GetYamlConfig,
	RemoveFromYamlConfig,
} from "@wailsjs/go/functools/FuncService";
import { ShowInfoMessage } from "@wailsjs/go/main/App";
import { onMounted, ref } from "vue";
import { getMinecraftDownloadFileName, uniqueBy } from "@/api/utils";
import ModDependsList from "@/components/ModDepends/ModDependsList.vue";
import ModsList from "@/components/Mods/ModsList.vue";
import View from "@/components/View/View.vue";
import AppFooter from "@/layouts/AppFooter.vue";
import AppHeader from "@/layouts/AppHeader.vue";
import type { MinecraftMod, ModDependency } from "@/types";

const savedMods = ref<MinecraftMod[]>([]);
const depends = ref<ModDependency[]>([]);

const loadDownloadedMods = async () => {
	try {
		const mods = await GetYamlConfig("downloads");
		savedMods.value = mods ?? [];

		if (savedMods.value !== null) {
			const uniqeDeps = uniqueBy(
				savedMods.value.flatMap((m) => m.Dependency),
				(d) => d.Name,
			);
			depends.value = uniqeDeps;
			console.log(depends.value);
		}
	} catch (err) {
		console.error("Ошибка при получении скачанных модов:", err);
	}
};

const removeFromDownloads = async (mod: MinecraftMod) => {
	try {
		savedMods.value = savedMods.value.filter((m) => m.Name !== mod.Name);
		await RemoveFromYamlConfig(mod, "downloads");
		const filename = getMinecraftDownloadFileName(mod.Name, mod.Versions);
		await DeleteSavedMod(filename);
		await ShowInfoMessage("Удалён", `Мод "${mod.Name}" успешно удалён`);
	} catch (err) {
		console.error("Ошибка при удалении мода:", err);
	}
};

onMounted(loadDownloadedMods);
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
      
      <ModDependsList v-if="depends.length > 0" :depends="depends"/>
      <span v-else class="downloads-title"> Нет зависимостей </span>
    
    </div>
  </View>
  <AppFooter />
</template>

<style lang="css" scoped>
.no-mods {
  text-align: center;
  color: #ccc;
  font-size: 18px;
  margin-top: 40px;
}

.downloads-title {
  display: flex;
  justify-content: center;
  font-size: 30px;
  color: white;
  padding: 10px;
}

</style>
