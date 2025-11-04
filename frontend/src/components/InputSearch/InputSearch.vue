<script setup lang="ts">
import { onMounted, ref } from "vue";
import SearchIcon from "@/assets/images/search.png"

interface Props<T> {
  placeholder?: string
}

const props = withDefaults(defineProps<Props<any>>(), {
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
        <div class="search-button">
            <img :src="SearchIcon"
                class="icon" 
                style="cursor: pointer; padding: 5px;"
                @click="handleSearch"
            />
        </div>
    </div>
</template>

<style lang="css" scoped>
.input-search {
    display: flex;
    justify-content: row;

    gap: 10px;
}

.search {
    display: flex;

    padding: 10px;

    border-radius: 10px;
    border-style: none;
    
    background-color: rgb(34, 33, 33);
    border: 1px solid whitesmoke;
    outline: none;
}

.search-button {
    display: flex;
    align-items: center;
    padding: 5px;
    border-radius: 30%;
    background-color: lightgray;
    transition: 1s;
}

.search-button:hover {
    background-color: lightgreen;
}

</style>