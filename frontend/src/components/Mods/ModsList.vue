<script setup lang="ts">
import { useRouter } from 'vue-router';
import { useModStore } from '@/stores/modStore';

import type { MinecraftMod } from '@/types';
import Tag from '@/components/Tags/Tag.vue';
import TagImage from '@/components/Tags/TagImage.vue';
import TagImageGroup from '@/components/Tags/TagImageGroup.vue';

import ForgeMinecraftIcon from '@/assets/images/forge_minecraft.jpeg'
import FabricMinecraftIcon from '@/assets/images/fabric_minecraft.png'
import DeleteIcon from '@/assets/images/trash.png'

import { uniqueBy } from '@/api/utils';
import { useScrollManager } from '@/composables/useScrollManager';

const router = useRouter()
const modStore = useModStore()

const { saveScrollBeforeLeave } = useScrollManager()

const props = defineProps<{
    mods: MinecraftMod[],
    loaderf?: string,
    versionf?: string,
    onDelete?: (mod: MinecraftMod) => void
}>()

const loaderToTag = (loader: string) => {
  const loaderMap: Record<string, { name: string; imageUrl: string }> = {
    'forge': { 
      name: 'Forge', 
      imageUrl: ForgeMinecraftIcon 
    },
    'fabric': { 
      name: 'Fabric', 
      imageUrl: FabricMinecraftIcon 
    }
  };
  
  return loaderMap[loader] || { name: loader, imageUrl: '' };
}

const getLoaderTags = (mod: MinecraftMod) => {
  return mod.Loaders.map(loader => {
    const tag = loaderToTag(loader);
    return {
      id: loader,
      name: tag.name,
      imageUrl: tag.imageUrl,
      altText: `${tag.name} loader`
    };
  }).filter(l => l.name !== 'quilt');
}

const redirectToMod = (mod: MinecraftMod): void => {
  saveScrollBeforeLeave()
  modStore.setCurrentMod(mod)

  if (props.loaderf !== undefined && props.loaderf !== "") {
    modStore.setLoaderFilter(props.loaderf)
  }

  if (props.versionf !== undefined && props.versionf !== "") {
    modStore.setVersionFilter(props.versionf)
  }

  router.push({
    name: 'mod',
    params: {
      name: mod.Name
    }
  })
}

const handleDelete = (mod: MinecraftMod, event: MouseEvent) => {
  event.stopPropagation()
  props.onDelete?.(mod)
}

</script>

<template>
    <div class="mods-grid">
        <div 
            v-for="item in uniqueBy(mods, mod => mod.Name)" 
            :key="item.ModPageLink"
            class="mod-card"
            @click="redirectToMod(item)"
        >
              <button 
                v-if="onDelete"
                class="delete-btn"
                @click="handleDelete(item, $event)"
                title="Удалить мод"
              >
                <img :src="DeleteIcon" alt="Удалить" />
              </button>
            <TagImageGroup 
                v-if="getLoaderTags(item).length > 0"
                :tags="getLoaderTags(item)"
                class="loader-tags"
            />
            
            <TagImage
                v-else-if="getLoaderTags(item).length === 1"
                v-bind="getLoaderTags(item)[0]"
                class="loader-tag"
                :name="item.Name"
            />
            
            <img 
                :src="item.Icon" 
                :alt="item.Name"
                class="icon"
            />
            <p class="mod-name">{{ item.Name }}</p>
            <div class="version-tags">
                <Tag v-for="version in item.Versions" :name="version"/>
            </div>
        </div>
    </div>
</template>

<style lang="css" scoped>
.mods-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
}

.mod-card {
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 15px;
  text-align: center;
  position: relative;
  cursor: pointer;
}

.mod-card:hover{
  background-color: rgb(126, 136, 124);
  transition: 0.9s;
}

.loader-tags {
  position: absolute;
  top: 8px;
  left: 8px;
  justify-content: flex-start;
}

.loader-tag {
  position: absolute;
  top: 8px;
  left: 8px;
}

.icon {
  width: 100px;
  height: 100px;
  object-fit: cover;
  border-radius: 8px;
  margin-top: 10px;
}

.mod-name {
  margin-top: 10px;
  font-weight: bold;
  word-break: break-word;
}

.version-tags {
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
  gap: 4px;
  margin-top: 8px;
}

.delete-btn {
  position: absolute;
  top: 6px;
  right: 6px;
  background: transparent;
  border: none;
  cursor: pointer;
  opacity: 0.6;
  transition: opacity 0.2s ease;
  z-index: 2;
}

.delete-btn:hover {
  opacity: 1;
}

.delete-btn img {
  width: 18px;
  height: 18px;
  filter: brightness(0) invert(1);
}

</style>