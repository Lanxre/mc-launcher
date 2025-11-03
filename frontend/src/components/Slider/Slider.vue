<template>
  <div class="slider">
    <div class="slider-container">
      <div 
        class="slides"
        :style="{ transform: `translateX(-${currentIndex * 100}%)` }"
      >
        <div
          v-for="(slide, index) in slides"
          :key="index"
          class="slide"
        >
          <img 
            :src="slide.image" 
            :alt="slide.alt || `Slide ${index + 1}`"
            class="slide-image"
          >
          <div v-if="slide.title" class="slide-title">
            {{ slide.title }}
          </div>
        </div>
      </div>
    </div>

    <button 
      v-if="showArrows"
      class="slider-arrow slider-arrow--prev"
      @click="prevSlide"
      :disabled="currentIndex === 0"
    >
      ‹
    </button>
    <button 
      v-if="showArrows"
      class="slider-arrow slider-arrow--next"
      @click="nextSlide"
      :disabled="currentIndex === slides.length - 1"
    >
      ›
    </button>

    <div v-if="showDots" class="slider-dots">
      <button
        v-for="(_, index) in slides"
        :key="index"
        class="slider-dot"
        :class="{ 'slider-dot--active': currentIndex === index }"
        @click="goToSlide(index)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue'

interface Slide {
  image: string
  alt?: string
  title?: string
}

interface Props {
  slides: Slide[]
  autoplay?: boolean
  autoplayInterval?: number
  showArrows?: boolean
  showDots?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  autoplay: false,
  autoplayInterval: 3000,
  showArrows: true,
  showDots: true
})

const currentIndex = ref(0)
let autoplayTimer: number | null = null

const nextSlide = () => {
  if (currentIndex.value < props.slides.length - 1) {
    currentIndex.value++
  } else {
    currentIndex.value = 0 
  }
}

const prevSlide = () => {
  if (currentIndex.value > 0) {
    currentIndex.value--
  } else {
    currentIndex.value = props.slides.length - 1 
  }
}

const goToSlide = (index: number) => {
  currentIndex.value = index
}

const startAutoplay = () => {
  if (props.autoplay) {
    autoplayTimer = window.setInterval(() => {
      nextSlide()
    }, props.autoplayInterval)
  }
}

const stopAutoplay = () => {
  if (autoplayTimer) {
    clearInterval(autoplayTimer)
    autoplayTimer = null
  }
}

watch(() => props.autoplay, (newValue) => {
  if (newValue) {
    startAutoplay()
  } else {
    stopAutoplay()
  }
})

onMounted(() => {
  if (props.autoplay) {
    startAutoplay()
  }
})

onUnmounted(() => {
  stopAutoplay()
})
</script>

<style lang="css" scoped>
.slider {
  position: relative;
  width: 100%;
  max-width: 800px;
  margin: 0 auto;
  overflow: hidden;
  border-radius: 8px;
}

.slider-container {
  overflow: hidden;
  border-radius: 8px;
}

.slides {
  display: flex;
  transition: transform 0.3s ease-in-out;
}

.slide {
  flex: 0 0 100%;
  position: relative;
}

.slide-image {
  width: 100%;
  height: 450px;
  object-fit: cover;
  display: block;
}

.slide-title {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  background: linear-gradient(transparent, rgba(0, 0, 0, 0.7));
  color: white;
  padding: 20px;
  text-align: center;
  font-size: 1.2rem;
}

.slider-arrow {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  background: rgba(255, 255, 255, 0.8);
  border: none;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  font-size: 20px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.3s;
}

.slider-arrow:hover:not(:disabled) {
  background: rgba(255, 255, 255, 1);
}

.slider-arrow:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.slider-arrow--prev {
  left: 10px;
}

.slider-arrow--next {
  right: 10px;
}

.slider-dots {
  position: absolute;
  bottom: 20px;
  left: 0;
  right: 0;
  display: flex;
  justify-content: center;
  gap: 8px;
}

.slider-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  border: none;
  background: rgba(255, 255, 255, 0.5);
  cursor: pointer;
  transition: background-color 0.3s;
}

.slider-dot--active {
  background: rgba(255, 255, 255, 1);
}

.slider-dot:hover {
  background: rgba(255, 255, 255, 0.8);
}
</style>