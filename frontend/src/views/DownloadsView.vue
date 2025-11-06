<script setup lang="ts">
import AppFooter from '@/layouts/AppFooter.vue';
import AppHeader from '@/layouts/AppHeader.vue';
import View from '@/components/View/View.vue';
import ImageButton from '@/components/Buttons/ImageButton.vue';
import CrossIcon from '@/assets/images/close.png'


import { onMounted, ref } from 'vue';
import { GetSavedMods, DeleteSavedMod} from '../../wailsjs/go/main/FuncService' 

const savedNameMods = ref<string[]>([])

const handleDeleteMod = async (modName: string) => {
    await DeleteSavedMod(modName)
    savedNameMods.value = savedNameMods.value.filter((m) => m != modName)
}

onMounted(async () => {
    savedNameMods.value = await GetSavedMods() 
})

</script>

<template>
    <AppHeader/>
    <View>
        <div class="downloads-list">
            <div v-for="modName in savedNameMods" :key="modName" class="mod-item">
               <div class="mod-name"> {{ modName }} </div>
               <div class="mod-details"> 
                  <ImageButton :img="CrossIcon" @click="handleDeleteMod(modName)" title="Удалить"/>  
               </div>
            </div>
        </div>
    </View>
    <AppFooter/>
</template>

<style lang="css" scoped>

.downloads-list {
    display: flex;
    flex-direction: column;

    height: 25px;
    padding: 10px;
    margin: 5px;

    height: 100%;
    min-height: 100vh;
}

.mod-item {
    display: flex;
    flex-direction: row;
}

.mod-name {
    display: flex;
    align-items: center;
    justify-content: center;

    width: 100%;
}

.mod-details {
    display: flex;
    flex-direction: row;
    justify-content: flex-end;

    width: 300px;
    gap: 5px;
}

</style>