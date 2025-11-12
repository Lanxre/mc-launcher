<script setup lang="ts">

import ImageButton from '../Buttons/ImageButton.vue';
import InputSearch from '../InputSearch/InputSearch.vue';
import SearchFilter from '../SearchFilter/SearchFilter.vue';

import StarIcon from '@/assets/images/star.png';
import FileIcon from '@/assets/images/file.png';


interface Props {
    versionE?: string;
    loaderE?: string;
    searchE: string;
    goToFavourites: () => void;
    goToDownloads: () => void;
    applySearch:  () =>  Promise<void>;
    applyFilters: () => Promise<void>;
    getUniqueVersions: string[];
    getUniqueLoaders:  string[]; 
}

const emit = defineEmits<{
  "update:searchE": [value: string];
  "update:versionE": [value: string];
  "update:loaderE": [value: string];
}>();

const updateSearch = (value: string) => {
  emit("update:searchE", value);
};

defineProps<Props>();

</script>

<template>
    <div class="filters">
        <ImageButton :img="StarIcon" @click="goToFavourites" border-radius="50%" title="Список избранных"/>
        <ImageButton :img="FileIcon" @click="goToDownloads" border-radius="50%" title="Список установленных модов"/>
        <InputSearch :modelValue="searchE" @update:modelValue="updateSearch" @search="applySearch" />
        <SearchFilter
            :versions="getUniqueVersions"
            :loaders="getUniqueLoaders"
            @filter="(v,l) => { 
                emit('update:versionE', v); 
                emit('update:loaderE', l);
                applyFilters();
            }"
        />
    </div>
</template>

<style lang="css" scoped>
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