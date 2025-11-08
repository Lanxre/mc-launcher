<script setup lang="ts">
import { onMounted, ref } from "vue";
import SearchIcon from "@/assets/images/search.png"
import ImageButton from "../Buttons/ImageButton.vue";

interface Props {
  placeholder?: string
}

withDefaults(defineProps<Props>(), {
  placeholder: 'Название мода',
})

const searchValue = ref("")

const emit = defineEmits<{
  'update:modelValue': [value: string],
  'search': [inputSearch: string],
}>()

const handleEnterSearch = (e: KeyboardEvent) => {
    if(e.code === "Enter") {
        handleSearch()
    }
}

const handleSearch = () => {
    emit('update:modelValue', searchValue.value)
    emit('search', searchValue.value)
}

onMounted(() => {
    document.addEventListener('keyup', handleEnterSearch)
})

</script>

<template>
    <div class="input-search">
        <input v-model="searchValue" class="search" type="text" :placeholder="placeholder"/>
        <ImageButton :img="SearchIcon" @click="handleSearch" style="border-radius: 50%;"/>
    </div>
</template>

<style lang="css" scoped>
.input-search {
    display: flex;
    justify-content: row;
    
    align-items: center;
    gap: 10px;
    color: white;
}

.search {
    display: flex;

    padding: 10px;

    border-radius: 10px;
    border-style: none;
    
    background-color: rgb(34, 33, 33);
    border: 1px solid whitesmoke;
    outline: none;
    color: white;
}

.search-button {
    display: flex;
    align-items: center;
    border-radius: 50%;
    background-color: lightgray;
    cursor: pointer;
    transition: background-color 0.3s;
}

.search-button img {
    padding: 10px;
}

.search-button:hover {
    background-color: lightgreen;
}

</style>