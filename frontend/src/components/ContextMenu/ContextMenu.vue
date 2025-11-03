<script setup lang="ts">
import type { ContextMenuItem } from '@/types';
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'


interface Props {
  items: ContextMenuItem[]
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false
})

const emit = defineEmits<{
  'item-click': [item: ContextMenuItem]
  'menu-show': []
  'menu-hide': []
}>()

const isVisible = ref(false)
const menuPosition = ref({ x: 0, y: 0 })
const menuRef = ref<HTMLElement>()
const activeSubmenu = ref<string | null>(null)
const submenuTimer = ref<number>()

const menuStyle = computed(() => ({
  left: `${menuPosition.value.x}px`,
  top: `${menuPosition.value.y}px`,
  display: isVisible.value ? 'block' : 'none'
}))

const submenuStyle = computed(() => ({
  left: '100%',
  top: '0'
}))

const handleContextMenu = (event: MouseEvent) => {
  if (props.disabled) return

  event.preventDefault()
  isVisible.value = true
  activeSubmenu.value = null

  nextTick(() => {
    positionMenu(event.clientX, event.clientY)
    emit('menu-show')
  })
}

const positionMenu = (x: number, y: number) => {
  if (!menuRef.value) return

  const menuWidth = menuRef.value.offsetWidth
  const menuHeight = menuRef.value.offsetHeight
  const windowWidth = window.innerWidth
  const windowHeight = window.innerHeight

  // Корректируем позицию чтобы меню не выходило за границы экрана
  let adjustedX = x
  let adjustedY = y

  if (x + menuWidth > windowWidth) {
    adjustedX = windowWidth - menuWidth - 10
  }

  if (y + menuHeight > windowHeight) {
    adjustedY = windowHeight - menuHeight - 10
  }

  menuPosition.value = {
    x: adjustedX,
    y: adjustedY
  }
}

const handleItemClick = (item: ContextMenuItem) => {
  if (item.handler) {
    item.handler()
  }
  emit('item-click', item)
  closeMenu()
}

const openSubmenu = (item: ContextMenuItem) => {
  if (submenuTimer.value) {
    clearTimeout(submenuTimer.value)
  }
  activeSubmenu.value = item.id
}

const closeSubmenu = () => {
  submenuTimer.value = window.setTimeout(() => {
    activeSubmenu.value = null
  }, 150)
}

const closeMenu = () => {
  isVisible.value = false
  activeSubmenu.value = null
  emit('menu-hide')
}

const handleClickOutside = (event: MouseEvent) => {
  if (menuRef.value && !menuRef.value.contains(event.target as Node)) {
    closeMenu()
  }
}

const handleEscapeKey = (event: KeyboardEvent) => {
  if (event.key === 'Escape' && isVisible.value) {
    closeMenu()
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
  document.addEventListener('keydown', handleEscapeKey)
  document.addEventListener('contextmenu', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
  document.removeEventListener('keydown', handleEscapeKey)
  document.removeEventListener('contextmenu', handleClickOutside)
  
  if (submenuTimer.value) {
    clearTimeout(submenuTimer.value)
  }
})
</script>


<template>
  <div class="context-menu-container" @contextmenu.prevent="handleContextMenu">
    <slot></slot>

    <div
      v-if="isVisible"
      ref="menuRef"
      class="context-menu"
      :style="menuStyle"
      @click.stop
    >
      <div
        v-for="item in items"
        :key="item.id"
        class="context-menu__item"
        :class="{
          'context-menu__item--disabled': item.disabled,
          'context-menu__item--danger': item.danger
        }"
        @click="!item.disabled && handleItemClick(item)"
      >
        <span class="context-menu__item-icon" v-if="item.icon">
          {{ item.icon }}
        </span>
        <span class="context-menu__item-label">{{ item.label }}</span>
        <span class="context-menu__item-shortcut" v-if="item.shortcut">
          {{ item.shortcut }}
        </span>
      </div>

      <div
        v-if="items.some(item => item.children)"
        class="context-menu__divider"
      ></div>

      <div
        v-for="item in items.filter(item => item.children)"
        :key="`sub-${item.id}`"
        class="context-menu__item context-menu__item--has-children"
        @mouseenter="openSubmenu(item)"
        @mouseleave="closeSubmenu"
      >
        <span class="context-menu__item-icon" v-if="item.icon">
          {{ item.icon }}
        </span>
        <span class="context-menu__item-label">{{ item.label }}</span>
        <span class="context-menu__item-arrow">▶</span>

        <div
          v-if="activeSubmenu === item.id"
          class="context-menu__submenu"
          :style="submenuStyle"
        >
          <div
            v-for="child in item.children"
            :key="child.id"
            class="context-menu__item"
            :class="{
              'context-menu__item--disabled': child.disabled,
              'context-menu__item--danger': child.danger
            }"
            @click="!child.disabled && handleItemClick(child)"
          >
            <span class="context-menu__item-icon" v-if="child.icon">
              {{ child.icon }}
            </span>
            <span class="context-menu__item-label">{{ child.label }}</span>
            <span class="context-menu__item-shortcut" v-if="child.shortcut">
              {{ child.shortcut }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>


<style lang="css" scoped>
.context-menu-container {
  width: 100%;
  height: 100%;
}

.context-menu {
  position: fixed;
  z-index: 1000;
  min-width: 200px;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
  padding: 4px;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  font-size: 14px;
  animation: contextMenuAppear 0.15s ease-out;
}

.context-menu__item {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.15s ease;
  position: relative;
  user-select: none;
}

.context-menu__item:hover {
  background: #f1f5f9;
}

.context-menu__item--disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.context-menu__item--disabled:hover {
  background: transparent;
}

.context-menu__item--danger {
  color: #dc2626;
}

.context-menu__item--danger:hover {
  background: #fef2f2;
}

.context-menu__item-icon {
  margin-right: 8px;
  width: 16px;
  text-align: center;
  font-size: 12px;
}

.context-menu__item-label {
  flex: 1;
}

.context-menu__item-shortcut {
  margin-left: 12px;
  font-size: 12px;
  color: #64748b;
}

.context-menu__item-arrow {
  margin-left: 8px;
  font-size: 10px;
  color: #64748b;
}

.context-menu__item--has-children {
  position: relative;
}

.context-menu__submenu {
  position: absolute;
  min-width: 160px;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
  padding: 4px;
  animation: submenuAppear 0.15s ease-out;
}

.context-menu__divider {
  height: 1px;
  background: #e2e8f0;
  margin: 4px 0;
}

@keyframes contextMenuAppear {
  from {
    opacity: 0;
    transform: scale(0.95) translateY(-5px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

@keyframes submenuAppear {
  from {
    opacity: 0;
    transform: translateX(-5px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}
</style>