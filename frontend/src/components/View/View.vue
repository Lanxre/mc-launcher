<script lang="ts" setup>
const props = defineProps<{
	minHeight?: string;
}>();
</script>

<template>
    <div class="main-content" :style="{ minHeight: props.minHeight }">
        <slot></slot>
    </div>
</template>

<style lang="css" scoped>
.main-content {
  display: flex;
  flex-direction: column;
  width: 100%;
  min-height: 100vh;
  padding: 20px;
  gap: 15px;
  background-color: var(--main-bg-color);
}

</style>