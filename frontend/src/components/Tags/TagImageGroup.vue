<script lang="ts" setup>
import TagImage from "@/components/Tags/TagImage.vue";

import type { PropType } from "vue";

import type { Tag } from "@/types/tag";

defineProps({
	tags: {
		type: Array as PropType<Tag[]>,
		default: () => [] as Tag[],
	},
});

const emit = defineEmits<{
	"tag-click": [tag: Tag, event: Event];
	"tag-close": [tag: Tag];
}>();
</script>

<template>
  <div class="tag-image-group">
    <slot>
      <TagImage
        v-for="tag in tags"
        :key="tag.id"
        :name="tag.name"
        :image-url="tag.imageUrl"
        :alt-text="tag.altText"
        :clickable="tag.clickable"
        :closable="tag.closable"
        :size="tag.size"
        :variant="tag.variant"
        :rounded="tag.rounded"
        @click="(e: Event) => emit('tag-click', tag, e)"
        @close="() => emit('tag-close', tag)"
      />
    </slot>
  </div>
</template>

<style lang="css" scoped>
.tag-image-group {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  align-items: center;
}
</style>