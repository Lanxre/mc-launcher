<script setup lang="ts">
import { GetMinecraftVersions, OpenModsFolder } from "@wailsjs/go/functools/FuncService";
import { LaunchMinecraft, LaunchMinecraftD } from "@wailsjs/go/launcher/LauncherService";
import { ShowInfoMessage } from "@wailsjs/go/main/App";
import { onMounted, ref } from "vue";
import PlayIcon from "@/assets/images/play.png";
import ImageButton from "../Buttons/ImageButton.vue";
import List from "../List/List.vue";
import Modal from "../Modal/Modal.vue";

const isModalOpen = ref<boolean>(false);
const selectedVersion = ref<string>("");
const minecraftVersion = ref<string[]>([]);

const openModal = () => (isModalOpen.value = true);

const load = async () => {
	try {
		const mcVersions = await GetMinecraftVersions();

		if (mcVersions !== undefined && mcVersions !== null) {
			minecraftVersion.value = mcVersions;
		}
	} catch (err) {
		console.error("Ошибка загрузки конфигурации minecraft", err);
	}
};

const openMinecraft = async () => {
	if (selectedVersion.value) {
		// await LaunchMinecraft(selectedVersion.value);
    await LaunchMinecraftD()
		await ShowInfoMessage("Успех", "Игра запускается");
	} else {
		await ShowInfoMessage("Ошибка", "Не выбрана версия!");
	}
};

const openModFolder = async () => {
  await OpenModsFolder()
}

onMounted(load);
</script>

<template>
    <div class="minecraft-settings">
        <ImageButton :img="PlayIcon" @click="openModal" border-radius="50%" title="Настройка запуска игры"/>
        <Modal v-model="isModalOpen" title="Настройки запуска">
            <div class="settings">
                <div class="settings-item">
                    <p class="text text-shd">Доступные установки</p>
                    <List :items="minecraftVersion" v-model="selectedVersion" placeholder="Установки" />
                </div>
                <div class="settings-item">
                  <button class="button confirm-button" @click="openModFolder"> Открыть папку с модами </button>
                </div>
            </div>

            <template #footer>
                <button class="button confirm-button" @click="openMinecraft"> Запустить </button>
            </template>
        </Modal>
    </div>
</template>

<style lang="css" scoped>

.minecraft-settings {
    display: flex;
    align-items: center;
}

.button {
  display: flex;
  align-items: center;
  border-radius: 50%;
  background-color: lightgray;
  padding: 10px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.button:hover {
  background-color: lightgreen;
}

.settings {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 40px 0;
  height: 350px;
}

.settings-item {
  display: flex;
  gap: 10px;
  font-size: 12px;
  margin-bottom: 10px;
}

.settings-item p {
    display: flex;
    text-align: center;
}

.button {
  padding: 8px 16px;
  border-radius: 6px;
  color: white;
  font-weight: 500;
  margin-right: 8px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.reset-button {
  background-color: #a51717;
}

.confirm-button {
  background-color: #2d8a4d;
}

.button:hover {
  opacity: 0.9;
}


</style>