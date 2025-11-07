<script setup lang="ts">
import { openLink } from '@/api/utils';
import LinkIcon from '@/assets/images/link.png'

interface Props {
    name: string, 
    linkPage: string
}

withDefaults(defineProps<Props>(), {
    name: 'Название мода',
})

</script>

<template>
    <div class="mod-name" @click="openLink(linkPage)">
        <h1>{{ name }} </h1>
        <img :src="LinkIcon" alt="modlink" class="icon"/>
    </div>
</template>

<style lang="css" scoped>
.mod-name {
  display: flex;
  flex-direction: row;

  gap: 5px;
  cursor: pointer;
}

.mod-name img {
  width: 15px;
  height: 15px;
  filter: brightness(0) invert(1);
}
</style>