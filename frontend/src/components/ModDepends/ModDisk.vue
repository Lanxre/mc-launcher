<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { GetSavedMods } from '@wailsjs/go/functools/FuncService';

const savedMods = ref<string[]>([])

const load = async () => {
    try {
        savedMods.value = await GetSavedMods();
    } catch (err) {
        console.log("Ошибка загрузки с диска", err)
    }
}

onMounted(load)

</script>

<template>
    <div class="on-disk-mods">
        <span class="title" v-if="savedMods.length"> Mod on Disk </span>
        <span class="title" v-else> No have mods </span>

        <div v-for="mod in savedMods" :key="mod" class="list-mod shadow">
            {{ mod }}
        </div>
    </div>
</template>

<style lang="css" scoped>
.on-disk-mods{
  display: flex;
  flex-direction: column;
  padding: 10px;
  gap: 5px;
}

.list-mod {
    display: flex;
    justify-content: center;
    text-align: center;

    border: 1px solid white;
    border-radius: 30px;
    padding: 10px;
    margin-bottom: 5px;
}

.title {
    display: flex;
    justify-content: center;
    text-align: center;
    font-size: 35px;
    font-family: "Minecraftv1";
    margin-bottom: 5px;
}

</style>