<script setup lang="ts">
import {
	GetYamlConfig,
	RemoveFromYamlConfig,
} from "@wailsjs/go/functools/FuncService";
import { ShowInfoMessage } from "@wailsjs/go/main/App";
import { onMounted, ref } from "vue";
import StarIcon from "@/assets/images/star.png";
import ModsList from "@/components/Mods/ModsList.vue";
import View from "@/components/View/View.vue";
import AppFooter from "@/layouts/AppFooter.vue";
import AppHeader from "@/layouts/AppHeader.vue";
import type { MinecraftMod } from "@/types";

const mods = ref<MinecraftMod[]>([]);
const isLoading = ref(true);
const isError = ref(false);

const fileName = "favourite.yaml";

const loadFavourites = async () => {
	try {
		const result = await GetYamlConfig(fileName);
		mods.value = Array.isArray(result) ? result : [];
	} catch (err) {
		console.error("Ошибка при получении списка избранных модов:", err);
		isError.value = true;
	} finally {
		isLoading.value = false;
	}
};

const removeFromFavourites = async (mod: MinecraftMod) => {
	try {
		mods.value = mods.value.filter((m) => m.Name !== mod.Name);
		await RemoveFromYamlConfig(mod, fileName);
		await ShowInfoMessage("Удалён", `Мод: ${mod.Name} убран из избранных`);
	} catch (err) {
		console.error("Ошибка при удаление мода", err);
	}
};

onMounted(loadFavourites);
</script>

<template>
  <AppHeader />

  <View min-height="100vh">
    <div class="mod-favourite">
      <div class="mod-favourite__title">
        <span>Список избранных</span>
        <img class="icon" :src="StarIcon" alt="Избранное"/>
      </div>

      <div v-if="isLoading" class="mod-favourite__state">
        <p>Загрузка избранных модов...</p>
      </div>

      <div v-else-if="isError" class="mod-favourite__state error">
        <p>Не удалось загрузить список. Попробуйте позже.</p>
      </div>

      <ModsList v-else-if="mods.length" :mods="mods" v-on:delete="removeFromFavourites" />
      <div v-else class="mod-favourite__state">
        <p>У вас пока нет избранных модов 😢</p>
      </div>
    </div>
  </View>

  <AppFooter />
</template>

<style scoped>
.mod-favourite {
  display: flex;
  flex-direction: column;
  color: white;
  padding: 20px 0;
}

.mod-favourite__title {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;

  font-size: 65px;
  font-family: 'Minecraftv1';
  color: white;

  filter: brightness(0) invert(1);
  margin-bottom: 30px;
}

.mod-favourite__state {
  text-align: center;
  color: #ccc;
  font-size: 18px;
}

.mod-favourite__state.error {
  color: #ff6b6b;
}

.mod-favourite__title span {
    display: flex;
    justify-content: center;
    align-items: center;
    text-align: center;
}

</style>
