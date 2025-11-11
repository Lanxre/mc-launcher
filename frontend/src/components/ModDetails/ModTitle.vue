<script setup lang="ts">
import { ShowInfoMessage } from '@wailsjs/go/main/App'
import { openLink } from '@/api/utils';

import LinkIcon from '@/assets/images/link.png'
import StarIcon from '@/assets/images/star.png'

interface Props {
    name: string, 
    linkPage: string
}

const props = withDefaults(defineProps<Props>(), {
    name: 'Название мода',
})

const emit = defineEmits<{
  saveFavourite: [event: MouseEvent]
}>()

const handleClick = (event: MouseEvent) => {
  ShowInfoMessage(`${props.name}`,"Мод добавлен в избранные")
  emit('saveFavourite', event)
}

</script>

<template>
    <div class="mod-name">
        <h1>{{ name }} </h1>
        <img :src="LinkIcon" alt="modlink" class="icon" title="Открыть в браузере" @click="openLink(linkPage)"/>
        <img :src="StarIcon" alt="modlink" class="icon" title="Добавить в избранное" @click="handleClick"/>
    </div>
</template>

<style lang="css" scoped>
.mod-name {
  display: flex;
  flex-direction: row;

  gap: 5px;
  cursor: pointer;
  color: white;
}

.mod-name img {
  width: 15px;
  height: 15px;
  filter: brightness(0) invert(1);
}
</style>