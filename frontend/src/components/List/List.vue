<template>
  <div class="dropdown" ref="dropdownRef">
    <div
      class="dropdown__trigger"
      :class="{
        'dropdown__trigger--open': isOpen,
        'dropdown__trigger--disabled': disabled
      }"
      @click="toggleDropdown"
    >
      <slot name="trigger" :selected="selectedItems" :is-open="isOpen">
        <div class="dropdown__trigger-content">
          <span class="dropdown__placeholder">
            {{ getTriggerText() }}
          </span>
          <span class="dropdown__arrow" :class="{ 'dropdown__arrow--open': isOpen }">
            ▼
          </span>
        </div>
      </slot>
    </div>
    <Transition name="dropdown">
      <div v-if="isOpen" class="dropdown__menu">
        <div
          v-for="item in items"
          :key="getKey(item)"
          class="dropdown__item"
          :class="{
            'dropdown__item--selected': isSelected(item),
            'dropdown__item--disabled': isDisabled(item)
          }"
          @click="handleItemClick(item)"
        >
          <slot name="item" :item="item" :selected="isSelected(item)">
            <div class="dropdown__item-content">
              <span class="dropdown__item-text">{{ getItemLabel(item) }}</span>
            </div>
          </slot>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'

interface Props<T> {
  items: T[]
  modelValue?: T | T[]
  multiple?: boolean
  disabled?: boolean
  disabledItems?: T[]
  placeholder?: string
  itemKey?: keyof T | ((item: T) => string | number)
  itemLabel?: keyof T | ((item: T) => string)
}

const props = withDefaults(defineProps<Props<any>>(), {
  multiple: false,
  disabled: false,
  disabledItems: () => [],
  placeholder: 'Выберите элемент',
  itemKey: 'id',
  itemLabel: 'name'
})

const emit = defineEmits<{
  'update:modelValue': [value: any]
  'select': [item: any]
  'toggle': [isOpen: boolean]
}>()

const isOpen = ref(false)
const dropdownRef = ref<HTMLElement>()

const internalSelected = ref<any[]>([])

watch(() => props.modelValue, (newValue) => {
  if (newValue === undefined || newValue === null || newValue === '') {
    internalSelected.value = []
  } else {
    internalSelected.value = Array.isArray(newValue) ? [...newValue] : [newValue]
  }
}, { immediate: true })

const selectedItems = computed(() => internalSelected.value)

const getKey = (item: any): string | number => {
  if (typeof props.itemKey === 'function') {
    return props.itemKey(item)
  }
  return item[props.itemKey]
}

const getItemLabel = (item: any): string => {
  if (typeof props.itemLabel === 'function') {
    return props.itemLabel(item)
  }
  return item[props.itemLabel] || String(item)
}

const isSelected = (item: any): boolean => {
  return selectedItems.value.some(selected => getKey(selected) === getKey(item))
}

const isDisabled = (item: any): boolean => {
  return props.disabledItems.some(disabled => getKey(disabled) === getKey(item))
}

const getTriggerText = (): string => {
  if (selectedItems.value.length === 0 || !selectedItems.value[0]) {
    return props.placeholder
  }

  if (props.multiple) {
    if (selectedItems.value.length === 1) {
      return getItemLabel(selectedItems.value[0])
    }
    return `Выбрано: ${selectedItems.value.length}`
  }

  return getItemLabel(selectedItems.value[0])
}
const toggleDropdown = () => {
  if (props.disabled) return
  isOpen.value = !isOpen.value
  emit('toggle', isOpen.value)
}

const handleItemClick = (item: any) => {
  if (isDisabled(item)) return

  let newSelected: any[]

  if (props.multiple) {
    if (isSelected(item)) {
      newSelected = selectedItems.value.filter(
        selected => getKey(selected) !== getKey(item)
      )
    } else {
      newSelected = [...selectedItems.value, item]
    }
  } else {
    newSelected = [item]
    isOpen.value = false
    emit('toggle', false)
  }

  internalSelected.value = newSelected

  const emitValue = props.multiple ? newSelected : newSelected[0] || null
  emit('update:modelValue', emitValue)
  emit('select', item)
}

const handleClickOutside = (event: MouseEvent) => {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target as Node)) {
    isOpen.value = false
    emit('toggle', false)
  }
}

const handleEscape = (event: KeyboardEvent) => {
  if (event.key === 'Escape' && isOpen.value) {
    isOpen.value = false
    emit('toggle', false)
  }
}


onMounted(() => {
  document.addEventListener('click', handleClickOutside)
  document.addEventListener('keydown', handleEscape)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
  document.removeEventListener('keydown', handleEscape)
})
</script>

<style scoped>
.dropdown {
  position: relative;
  display: flex;
  width: 100%;
}

.dropdown__trigger {
  padding: 10px 16px;
  border: 1px solid #59595a;
  border-radius: 6px;
  background: var(--main-bg-color);
  cursor: pointer;
  transition: all 0.2s ease;
  user-select: none;
}

.dropdown__trigger:hover {
  border-color: #9ca3af;
}

.dropdown__trigger--open {
  border-color: #123432;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.dropdown__trigger--disabled {
  opacity: 0.5;
  cursor: not-allowed;
  background: #123432;
}

.dropdown__trigger--disabled:hover {
  border-color: #123432;
}

.dropdown__trigger-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
}

.dropdown__placeholder {
  color: #6b7280;
  width: 60px;
}

.dropdown__arrow {
  transition: transform 0.2s ease;
  color: #6b7280;
  font-size: 12px;
}

.dropdown__arrow--open {
  transform: rotate(180deg);
}

.dropdown__menu {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: var(--main-bg-color);
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
  max-height: 200px;
  overflow-y: auto; 
  z-index: 1000;
  margin-top: 4px;
}

.dropdown__item {
  padding: 10px 16px;
  cursor: pointer;
  transition: background-color 0.2s ease;
  border-bottom: 1px solid #f3f4f6;
  min-height: 20px;
}

.dropdown__item:last-child {
  border-bottom: none;
}

.dropdown__item:hover {
  background-color: #123432;
}

.dropdown__item--selected {
  background-color: var(--main-bg-color);
  color: #ffffff;
}

.dropdown__item--disabled {
  opacity: 0.5;
  cursor: not-allowed;
  background-color: #f9fafb;
}

.dropdown__item--disabled:hover {
  background-color: #f9fafb;
}

.dropdown__item-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.dropdown__item-text {
  flex: 1;
}

.dropdown__item-checkmark {
  color: #3b82f6;
  font-weight: bold;
  margin-left: 8px;
}

/* Анимации */
.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.2s ease;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>