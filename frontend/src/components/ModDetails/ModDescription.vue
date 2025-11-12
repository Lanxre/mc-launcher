<script setup lang="ts">
interface Props {
	description: string;
}

withDefaults(defineProps<Props>(), {
	description: "Описание мода",
});
</script>

<template>
    <div class="mod-description">
        <span> {{ description }} </span>
    </div>
</template>

<style lang="css" scoped>
.mod-description {
    margin-bottom: 3px;
    color: white;
    background-color: var(--main-bg-color);
}
</style>