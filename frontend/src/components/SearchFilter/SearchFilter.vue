<script lang="ts" setup>
import { ref } from 'vue';

import SettingsIcon from '../../assets/images/setting.png'
import Modal from '../Modal/Modal.vue';
import List from '../List/List.vue';

interface Props {
    versions: string[]
    loaders: string[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'filter': [version: string, loader: string]
}>()

const isModalOpen = ref(false);
const version     = ref("");
const loader      = ref("");


const handleModalOpen = () => {
    isModalOpen.value = true
}

const handleModalConfirm = () => {
    isModalOpen.value = false
    emit('filter', version.value, loader.value)
}

const handleModalReset = () => {
    isModalOpen.value = false
    
    version.value = ""
    loader.value  = ""

    emit('filter', version.value, loader.value)
}

</script>

<template>
<div class="search-filter">
      <img :src="SettingsIcon" 
            class="icon" 
            style="cursor: pointer;"
            @click="handleModalOpen"
        />
      <Modal  v-model="isModalOpen"
              title="Фильтры поиска"
      >
      <div class="search-settings">
        <div class="search-settings-item">
            <p class="text text-shd">Загрузчик</p>
            <List :items="loaders" v-model="loader" placeholder="Загрузчик" style="width: 100px;"/>
        </div>
        
        <div class="search-settings-item">
            <p class="text text-shd">Версия</p>
            <List :items="versions" v-model="version" placeholder="Версия" style="width: 100px;"/>
        </div>
        
      </div>

      <template #footer>
        <button class="button"
                style="background-color: #a51717;"
                @click="handleModalReset">
            Сбросить
        </button>
        <button class="button"
                style="background-color: #123432;"
                @click="handleModalConfirm">
            Подтвердить
        </button>
      </template>
      </Modal>
    </div>
</template>

<style lang="css">

.search-filter {
    display: flex;
    flex-direction: row;
    justify-content: flex-end;
    gap: 12px;

    align-items: center;
}

.search-settings {
    display: flex;
    flex-direction: column;

    gap: 10px;
    padding: 40px 0 40px;
    height: 350px;
}

.search-settings-item {
    display: flex;
    flex-direction: row;

    justify-content: space-between;

}

.search-settings-item {
    font-size: 12px;
}

</style>