<template>
  <Transition name="modal">
    <div v-if="isVisible" class="modal-overlay" @click.self="handleOverlayClick">
      <div class="modal" :class="modalClass" :style="modalStyle">
        <div class="modal__header" v-if="$slots.header || title">
          <slot name="header">
            <h2 class="modal__title">{{ title }}</h2>
          </slot>
          <button
            v-if="closable"
            class="modal__close"
            @click="close"
            aria-label="Закрыть"
          >
            ×
          </button>
        </div>

        <div class="modal__content">
          <slot></slot>
        </div>

        <div class="modal__footer" v-if="$slots.footer">
          <slot name="footer"></slot>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, useSlots } from "vue";

interface Props {
	modelValue: boolean;
	title?: string;
	width?: number | string;
	maxWidth?: number | string;
	closable?: boolean;
	closeOnOverlay?: boolean;
	closeOnEscape?: boolean;
	preventScroll?: boolean;
	transition?: string;
}

const props = withDefaults(defineProps<Props>(), {
	modelValue: false,
	closable: true,
	closeOnOverlay: true,
	closeOnEscape: true,
	preventScroll: false,
	width: "auto",
	maxWidth: "600px",
});

const emit = defineEmits<{
	"update:modelValue": [value: boolean];
	close: [];
	open: [];
}>();

const slots = useSlots();

const isVisible = computed({
	get: () => props.modelValue,
	set: (value) => emit("update:modelValue", value),
});

const modalStyle = computed(() => ({
	width: typeof props.width === "number" ? `${props.width}px` : props.width,
	maxWidth:
		typeof props.maxWidth === "number" ? `${props.maxWidth}px` : props.maxWidth,
}));

const modalClass = computed(() => ({
	"modal--with-header": !!props.title || !!slots.header,
	"modal--with-footer": !!slots.footer,
}));

const open = () => {
	isVisible.value = true;
	emit("open");
};

const close = () => {
	isVisible.value = false;
	emit("close");
};

const handleOverlayClick = () => {
	if (props.closeOnOverlay) {
		close();
	}
};

const handleEscapeKey = (event: KeyboardEvent) => {
	if (event.key === "Escape" && props.closeOnEscape && isVisible.value) {
		close();
	}
};

const preventScroll = (event: Event) => {
	if (props.preventScroll && isVisible.value) {
		event.preventDefault();
	}
};

onMounted(() => {
	document.addEventListener("keydown", handleEscapeKey);

	if (props.preventScroll) {
		document.addEventListener("wheel", preventScroll, { passive: false });
		document.addEventListener("touchmove", preventScroll, { passive: false });
	}
});

onUnmounted(() => {
	document.removeEventListener("keydown", handleEscapeKey);
	document.removeEventListener("wheel", preventScroll);
	document.removeEventListener("touchmove", preventScroll);
});
</script>

<style lang="css" scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  padding: 20px;
  cursor: default;
}

.modal {
  background: var(--main-bg-color);
  border-radius: 12px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.modal__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px;
  margin-bottom: 16px;
  gap: 16px
}

.modal__title {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 600;
  color: whitesmoke;
}

.modal__close {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  padding: 4px;
  color: #6b7280;
  border-radius: 4px;
  transition: all 0.2s;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal__close:hover {
  background: #f3f4f6;
  color: #374151;
}

.modal__content {
  padding: 0 24px;
  flex: 1;
  overflow-y: auto;
  max-height: calc(90vh - 120px);
}

.modal__footer {
  padding: 20px 24px;
  background: var(--main-bg-color);
  border-top: 1px solid #585c64;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 20px;
}

/* Анимации */
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active .modal,
.modal-leave-active .modal {
  transition: transform 0.3s ease, opacity 0.3s ease;
}

.modal-enter-from .modal {
  opacity: 0;
  transform: scale(0.9) translateY(-10px);
}

.modal-leave-to .modal {
  opacity: 0;
  transform: scale(0.9) translateY(10px);
}

/* Responsive */
@media (max-width: 640px) {
  .modal-overlay {
    padding: 10px;
  }
  
  .modal {
    width: 100% !important;
    max-width: 100% !important;
    border-radius: 8px;
  }
  
  .modal__header,
  .modal__content,
  .modal__footer {
    padding-left: 20px;
    padding-right: 20px;
  }
}
</style>1