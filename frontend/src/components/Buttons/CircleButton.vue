<script lang="ts" setup>
import { ref, onMounted, onUnmounted } from 'vue'
import ArrowUpIcon from '@/assets/images/upload.png'

const showButton = ref(false)

const checkScrollPosition = () => {
  const scrollPosition = window.scrollY
  const windowHeight = window.innerHeight
  const documentHeight = document.documentElement.scrollHeight
  
  const middleThreshold = documentHeight / 2 - windowHeight / 2
  showButton.value = scrollPosition > middleThreshold
}

const scrollToTop = () => {
  window.scrollTo({
    top: 0,
    behavior: 'smooth'
  })
}

onMounted(() => {
  window.addEventListener('scroll', checkScrollPosition)
  checkScrollPosition()
})

onUnmounted(() => {
  window.removeEventListener('scroll', checkScrollPosition)
})
</script>

<template>
    <div 
        v-if="showButton" 
        class="button-circle"
        @click="scrollToTop"
    >
        <img class="icon" style="width: 20px; height: 20px; padding: 3px;" :src="ArrowUpIcon"/>
    </div>
</template>

<style lang="css" scoped>
.button-circle {
    position: fixed;
    bottom: 70px;
    right: 1%;
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