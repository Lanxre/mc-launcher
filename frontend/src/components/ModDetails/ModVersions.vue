<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  versions: string[]
}

const props = defineProps<Props>()

const hasVersions = computed(() => props.versions && props.versions.length > 0)

</script>

<template>
  <div class="mod-versions">
    <p v-if="hasVersions">Доступные версии:</p>
    <p v-else>Версии не указаны</p>

    <div v-if="hasVersions" class="versions-list">
      <span
        v-for="v in props.versions"
        :key="v"
        class="version-chip"
      >
        {{ v }}
      </span>
    </div>
  </div>
</template>

<style scoped>
.mod-versions {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 8px;
}

.versions-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.version-chip {
  background-color: #e0f2f1;
  border: 1px solid #b2dfdb;
  border-radius: 6px;
  padding: 4px 8px;
  font-size: 0.9rem;
  color: #00695c;
  font-weight: 500;
  transition: all 0.2s ease;
}

.version-chip:hover {
  background-color: #b2dfdb;
  cursor: pointer;
}
</style>
