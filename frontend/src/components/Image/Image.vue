<script setup lang="ts">
import { computed } from "vue";

interface Props {
	img: string;
	alt?: string;
	height?: string;
	width?: string;
	borderRaduis?: string;
	backgroundColor?: string;
}

const props = withDefaults(defineProps<Props>(), {
	alt: "no-img",
	height: "50px",
	width: "50px",
	backgroundColor: "grey",
});

const blockStyle = computed(() => ({
	borderRadius: props.borderRaduis,
	backgroundColor: props.backgroundColor,
}));
</script>

<template>
    <div class="image-block" :style="blockStyle">
        <img :src="img" :alt="alt" :style="{ height: height, width: width }"/>
    </div>
</template>

<style lang="css" scoped>
.image-block {
    display: flex;
    justify-content: center;
    padding: 10px;
    transition: cubic-bezier(0.165, 0.84, 0.44, 1);
}

.image-block:hover {
    background-color: rgb(235, 235, 235) !important;
    cursor: pointer;
}
</style>