<script lang="ts" setup>
import { onMounted, onUnmounted, ref } from "vue";
import ArrowUpIcon from "@/assets/images/upload.png";

interface Props {
	bottom?: string | number;
	right?: string | number;
	left?: string | number;
	top?: string | number;
	image?: string;
}

const props = withDefaults(defineProps<Props>(), {
	bottom: "70px",
	right: "1%",
	left: undefined,
	top: undefined,
});

const showButton = ref(false);

const checkScrollPosition = () => {
	const scrollPosition = window.scrollY;
	const windowHeight = window.innerHeight;
	const documentHeight = document.documentElement.scrollHeight;

	const middleThreshold = documentHeight / 2 - windowHeight / 2;
	showButton.value = scrollPosition > middleThreshold;
};

const scrollToTop = () => {
	window.scrollTo({
		top: 0,
		behavior: "smooth",
	});
};

onMounted(() => {
	window.addEventListener("scroll", checkScrollPosition);
	checkScrollPosition();
});

onUnmounted(() => {
	window.removeEventListener("scroll", checkScrollPosition);
});
</script>

<template>
    <div 
        v-if="showButton" 
        class="button-circle"
        :style="{
          bottom: typeof bottom === 'number' ? `${bottom}px` : bottom,
          right: typeof right === 'number' ? `${right}px` : right,
          left: typeof left === 'number' ? `${left}px` : left,
          top: typeof top === 'number' ? `${top}px` : top,
        }"
        @click="scrollToTop"
    >
        <img v-if="image" class="icon" style="width: 20px; height: 20px; padding: 3px;" :src="image"/>
        <img v-else class="icon" style="width: 20px; height: 20px; padding: 3px;" :src="ArrowUpIcon"/>
    </div>
</template>

<style lang="css" scoped>
.button-circle {
    position: fixed;
    width: 50px;
    height: 50px;
    background-color: grey;
    display: flex;
    justify-content: center;
    align-items: center;
    border-radius: 100%;
    z-index: 9999;
    transition: opacity 0.3s ease, transform 0.3s ease;
}

.button-circle:hover {
    background-color: lightslategray;
    cursor: pointer;
    transform: scale(1.1);
}
</style>