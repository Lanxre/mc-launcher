<script setup lang="ts">
import {
	DeleteSavedMod,
	GetYamlConfig,
	RemoveFromYamlConfig,
	RemoveFromDownloads,
	GetSavedMods
} from "@wailsjs/go/functools/FuncService";
import { RemoveAllMods } from "@wailsjs/go/filetools/FileService"
import { ShowInfoMessage } from "@wailsjs/go/main/App";
import { onMounted, ref } from "vue";
import { getMinecraftDownloadFileName, uniqueBy } from "@/api/utils";
import ModDependsList from "@/components/ModDepends/ModDependsList.vue";
import ModDisk from "@/components/ModDepends/ModDisk.vue";
import ModsList from "@/components/Mods/ModsList.vue";
import View from "@/components/View/View.vue";
import AppHeader from "@/layouts/AppHeader.vue";
import TrashIcon from "@/assets/images/trash.png"
import ImageButton from "@/components/Buttons/ImageButton.vue";
import type { MinecraftMod, ModDependency } from "@/types";

const savedMods = ref<MinecraftMod[]>([]);
const depends = ref<ModDependency[]>([]);

const savedModsOnDisk = ref<string[]>([]);

const loadDownloadedMods = async () => {
	try {
		const mods = await GetYamlConfig("downloads");
		savedModsOnDisk.value = await GetSavedMods();
		savedMods.value = mods ?? [];

		if (savedMods.value !== null) {
			const uniqeDeps = uniqueBy(
				savedMods.value.flatMap((m) => m.Dependency),
				(d) => d.Name,
			);
			depends.value = uniqeDeps;
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

const removeAll = async () => {
	try {
		await RemoveFromDownloads();
		await RemoveAllMods();
		savedMods.value = [];
		depends.value = [];
		savedModsOnDisk.value = [];
		await ShowInfoMessage("Успех", "Моды удалены");
	} catch (err) {
		console.log("Ошибка при очистке", err);
	}
}

const onDeleteModOnDisk = async (modName: string) => {
	savedModsOnDisk.value = savedModsOnDisk.value.filter((m) => m !== modName);
	await DeleteSavedMod(modName);
	await ShowInfoMessage("Успех", "Мод успешно удалён");
};

onMounted(loadDownloadedMods);
</script>

<template>
  <AppHeader />
  <View min-height="100vh">
	<div class="filters">
		<ImageButton :img="TrashIcon" border-radius="50%" @click="removeAll" title="Удалить всё"/>
	</div>
    <div class="downloads-list">
      <ModsList
        v-if="savedMods.length > 0"
        :mods="savedMods"
        :onDelete="removeFromDownloads"
      />
      <ModDependsList v-if="depends.length > 0" :depends="depends"/>
	  <ModDisk :saved-mods="savedModsOnDisk" @delete="onDeleteModOnDisk"/>
    </div>
  </View>
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
  font-size: 18px;
  color: white;
  padding: 10px;
}

.filters {
  display: flex;
  flex-direction: row;
  justify-content: flex-end;
  gap: 5px;

  background-color: rgb(41, 48, 47);
  border-radius: 30px;
  padding: 10px;
  -webkit-box-shadow: 5px 5px 5px -5px rgba(151, 223, 163, 0.6) inset;
  -moz-box-shadow: 5px 5px 5px -5px rgba(151, 223, 163, 0.6) inset;
  box-shadow: 5px 5px 5px -5px rgba(151, 223, 163, 0.6) inset;
}

</style>
