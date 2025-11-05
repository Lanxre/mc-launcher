<script lang="ts" setup>
interface Props {
  name?: string
  disabled?: boolean
  size?: 'small' | 'medium' | 'large'
}

const props = withDefaults(defineProps<Props>(), {
  name: 'Нажать',
  disabled: false,
  size: 'medium'
})

const emit = defineEmits<{
  click: [event: MouseEvent]
}>()

const handleClick = (event: MouseEvent) => {
  if (!props.disabled) {
    emit('click', event)
  }
}

const sizeClasses = {
  small: 'py-1 px-3 text-sm',
  medium: 'py-2 px-4 text-base',
  large: 'py-3 px-6 text-lg'
}
</script>

<template>
  <div class="b-main">
    <button
      :class="[
        'btn',
        sizeClasses[size],
        { 'btn-disabled': disabled }
      ]"
      :disabled="disabled"
      @click="handleClick"
    >
      {{ name }}
    </button>
  </div>
</template>

<style lang="css" scoped>
.b-main {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 10px;
  margin: 3px;
}

.btn {
  font-family: inherit;
  font-weight: 600;
  color: white;
  background-color: rgb(33, 180, 33);
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.py-1 { padding-top: 0.25rem; padding-bottom: 0.25rem; }
.px-3 { padding-left: 0.75rem; padding-right: 0.75rem; }
.text-sm { font-size: 0.875rem; line-height: 1.25rem; }

.py-2 { padding-top: 0.5rem; padding-bottom: 0.5rem; }
.px-4 { padding-left: 1rem; padding-right: 1rem; }
.text-base { font-size: 1rem; line-height: 1.5rem; }

.py-3 { padding-top: 0.75rem; padding-bottom: 0.75rem; }
.px-6 { padding-left: 1.5rem; padding-right: 1.5rem; }
.text-lg { font-size: 1.125rem; line-height: 1.75rem; }
</style>