<script setup lang="ts">
import { computed } from "vue";
import { useRouter } from "vue-router";
import { uniqueBy } from "@/api/utils";
import { useScrollManager } from "@/composables/useScrollManager";
import { useModStore } from "@/stores/modStore";
import type { MinecraftMod } from "@/types";

import FabricIcon from "@/assets/images/fabric_minecraft.png";
import ForgeIcon from "@/assets/images/forge_minecraft.jpeg";
import DeleteIcon from "@/assets/images/trash.png";

import Tag from "@/components/Tags/Tag.vue";
import TagImage from "@/components/Tags/TagImage.vue";
import TagImageGroup from "@/components/Tags/TagImageGroup.vue";

const props = defineProps<{
  mods: MinecraftMod[];
  onDelete?: (mod: MinecraftMod) => void;
}>();

const router = useRouter();
const modStore = useModStore();
const { saveScrollBeforeLeave } = useScrollManager();

const LOADER_CONFIG: Record<string, { name: string; imageUrl: string }> = {
  forge: { name: "Forge", imageUrl: ForgeIcon },
  fabric: { name: "Fabric", imageUrl: FabricIcon },
};

const processedMods = computed(() => {
  const uniqueMods = uniqueBy(props.mods, (mod) => mod.Name);
  
  return uniqueMods.map((mod) => {
    const loaderTags = (mod.Loaders || [])
      .filter((l) => l.toLowerCase() !== "quilt")
      .map((loader) => ({
        id: loader,
        name: LOADER_CONFIG[loader.toLowerCase()]?.name || loader,
        imageUrl: LOADER_CONFIG[loader.toLowerCase()]?.imageUrl || "",
        altText: `${loader} loader`,
      }));

    return { ...mod, loaderTags };
  });
});

const handleModClick = (mod: MinecraftMod) => {
  saveScrollBeforeLeave();
  modStore.setCurrentMod(mod);
  router.push({ name: "mod", params: { name: mod.Name } });
};

const handleDelete = (mod: MinecraftMod, event: MouseEvent) => {
  event.stopPropagation();
  props.onDelete?.(mod);
};
</script>

<template>
  <div class="mods-grid">
    <div
      v-for="item in processedMods"
      :key="item.ModPageLink"
      class="mod-card"
      @click="handleModClick(item)"
    >
      <button
        v-if="onDelete"
        class="delete-btn"
        title="Удалить мод"
        @click="handleDelete(item, $event)"
      >
        <img :src="DeleteIcon" alt="Удалить" />
      </button>

      <div class="loader-container">
        <TagImageGroup
          v-if="item.loaderTags.length > 1"
          :tags="item.loaderTags"
        />
        <TagImage
          v-else-if="item.loaderTags.length === 1"
          v-bind="item.loaderTags[0]"
        />
      </div>

      <img :src="item.Icon" :alt="item.Name" class="icon" loading="lazy" />
      
      <p class="mod-name">{{ item.Name }}</p>
      
      <div class="version-tags">
        <Tag 
          v-for="version in item.Versions" 
          :key="version" 
          :name="version" 
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
.mods-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
}

.mod-card {
  position: relative;
  border: 1px solid whitesmoke;
  border-radius: 8px;
  padding: 15px;
  text-align: center;
  cursor: pointer;
  transition: background-color 0.3s ease;
  box-shadow: inset 5px 5px 5px -5px rgba(255, 255, 255, 0.6);
}

.mod-card:hover {
  background-color: rgb(69, 71, 69);
}

.loader-container {
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
  color: white;
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