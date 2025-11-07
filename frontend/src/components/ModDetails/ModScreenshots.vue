<script setup lang="ts">
import { computed } from 'vue'
import Slider from '@/components/Slider/Slider.vue'

interface Props {
  screenshots: string[]
}

const props = defineProps<Props>()

const slides = computed(() =>
  props.screenshots?.length
    ? props.screenshots.map((s, i) => ({
        image: s,
        alt: `Скриншот ${i + 1}`,
        id: i,
      }))
    : []
)
</script>

<template>
  <div class="screenshots">
    <Slider v-if="slides.length" :slides="slides" />
    <div v-else class="no-screenshots">
      <p class="text text-shd">Скриншоты не найдены</p>
    </div>
  </div>
</template>

<style lang="css" scoped>
.no-screenshots {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
</style>