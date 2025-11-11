<script setup lang="ts">
interface Props {
	img: string;
	alt?: string;
	disabled?: boolean;
	title?: string;
	borderRadius?: string;
}

const props = withDefaults(defineProps<Props>(), {
	alt: "icon",
	disabled: false,
	borderRadius: "5px",
});

const emit = defineEmits<{
	click: [event: MouseEvent];
}>();

const handleClick = (event: MouseEvent) => {
	if (!props.disabled) {
		emit("click", event);
	}
};
</script>

<template>
  <div
    class="b-img"
    :class="{ 'b-img--disabled': disabled }"
    :style="{ borderRadius: props.borderRadius }"
    @click="handleClick"
    role="button"
    :aria-disabled="disabled"
    :tabindex="disabled ? -1 : 0"
    :title="title"
  >
    <img :src="img" :alt="alt" />
  </div>
</template>

<style scoped lang="css">
.b-img {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 10px;
  margin: 3px;
  background-color: rgb(231, 252, 252);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  user-select: none;
}

.b-img img {
  height: 25px;
  width: 25px;
  object-fit: contain;
}

.b-img:hover:not(.b-img--disabled) {
  background-color: rgba(231, 252, 252, 0.8);
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.b-img:active:not(.b-img--disabled) {
  transform: translateY(0);
}

.b-img--disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>