<script setup lang="ts">
import {
	DeleteSavedMod,
	GetSavedMods,
} from "@wailsjs/go/functools/FuncService";
import { ShowInfoMessage } from "@wailsjs/go/main/App";
import { onMounted, ref } from "vue";
import CrossIcon from "@/assets/images/close.png";
import Image from "@/components/Image/Image.vue";

const savedMods = ref<string[]>([]);

const onDelete = async (modName: string) => {
	savedMods.value = savedMods.value.filter((m) => m !== modName);
	await DeleteSavedMod(modName);
	await ShowInfoMessage("Успех", "Мод успешно удалён");
};

const load = async () => {
	try {
		savedMods.value = await GetSavedMods();
	} catch (err) {
		console.log("Ошибка загрузки с диска", err);
	}
};

onMounted(load);
</script>

<template>
    <div class="on-disk-mods">
        <span class="title" v-if="savedMods.length"> Mod on Disk </span>
        <div v-for="mod in savedMods" :key="mod" class="list-mod shadow">
            <div class="mod-name">{{ mod }}</div>
            <Image :img="CrossIcon" width="25px" height="25px" border-raduis="50%" title="Удалить" @click="onDelete(mod)"/>         
        </div>
    </div>
</template>

<style lang="css" scoped>
.on-disk-mods{
  display: flex;
  flex-direction: column;
  padding: 10px;
  gap: 5px;
  color: white;
}

.list-mod {
    display: flex;
    justify-content: center;
    text-align: center;
    align-items: center;

    border: 1px solid white;
    border-radius: 30px;
    padding: 10px;
    margin-bottom: 5px;
    background-color: var(--main-bg-color);
}

.title {
    display: flex;
    justify-content: center;
    text-align: center;
    font-size: 35px;
    font-family: "Minecraftv1";
    margin-bottom: 5px;
    color: white;
}

.mod-name {
    display: flex;
    justify-content: center;
    align-items: center;
    text-align: center;
    width: 100%;
    color: white;
}

</style>