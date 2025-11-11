<script lang="ts" setup>
import { ref } from "vue";

import SettingsIcon from "@/assets/images/setting.png";
import List from "@/components/List/List.vue";
import Modal from "@/components/Modal/Modal.vue";
import ImageButton from "../Buttons/ImageButton.vue";

interface Props {
	versions: string[];
	loaders: string[];
}

defineProps<Props>();

const emit = defineEmits<{
	filter: [version: string, loader: string];
}>();

const isModalOpen = ref(false);
const version = ref("");
const loader = ref("");

const openModal = () => (isModalOpen.value = true);

const confirmFilter = () => {
	isModalOpen.value = false;
	emit("filter", version.value, loader.value);
};

const resetFilter = () => {
	isModalOpen.value = false;
	version.value = "";
	loader.value = "";
	emit("filter", version.value, loader.value);
};
</script>

<template>
  <div class="search-filter">
    <ImageButton :img="SettingsIcon" @click="openModal" style="border-radius: 50%;" title="Настройки поиска"/>

    <Modal v-model="isModalOpen" title="Настройки поиска">
      <div class="search-settings">
        <div class="search-settings-item">
          <p class="text text-shd">Загрузчик</p>
          <List :items="loaders" v-model="loader" placeholder="Загрузчик" />
        </div>

        <div class="search-settings-item">
          <p class="text text-shd">Версия</p>
          <List :items="versions" v-model="version" placeholder="Версия" />
        </div>
      </div>

      <template #footer>
        <button class="button reset-button" @click="resetFilter">Сбросить</button>
        <button class="button confirm-button" @click="confirmFilter">Подтвердить</button>
      </template>
    </Modal>
  </div>
</template>

<style scoped>
.search-filter {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 12px;
}

.search-button {
  display: flex;
  align-items: center;
  border-radius: 50%;
  background-color: lightgray;
  padding: 10px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.search-button:hover {
  background-color: lightgreen;
}

.search-settings {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 40px 0;
  height: 350px;
}

.search-settings-item {
  display: flex;
  gap: 10px;
  font-size: 12px;
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

.icon {
  width: 24px;
  height: 24px;
}
</style>
