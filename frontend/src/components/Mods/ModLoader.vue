<template>
  <div>
    <div 
      v-if="hasMore" 
      ref="loaderTrigger" 
      class="loader-trigger"
    >
      <div v-if="loadingMore" class="loading-text">
        {{ loadingText }}
      </div>
    </div>
    
    <div v-else class="end-message">
      {{ endMessage }} ({{ count }})
    </div>
  </div>
</template>

<script lang="ts" setup>
import { onMounted, onUnmounted, ref, watch } from "vue";

interface Props {
	hasMore: boolean;
	loadingMore: boolean;
	count: number;
	loadingText?: string;
	endMessage?: string;
	autoLoad?: boolean;
	rootMargin?: string;
	threshold?: number | number[];
}

const props = withDefaults(defineProps<Props>(), {
	loadingText: "Загрузка...",
	endMessage: "Все моды загружены",
	autoLoad: true,
	rootMargin: "0px 0px 30px 0px",
	threshold: 0,
});

const emit = defineEmits<{
	loadMore: [];
}>();

const loaderTrigger = ref<HTMLElement | null>(null);
let observer: IntersectionObserver | null = null;

const initObserver = () => {
	if (!props.autoLoad || !loaderTrigger.value) return;

	observer = new IntersectionObserver(
		(entries) => {
			entries.forEach((entry) => {
				if (entry.isIntersecting && props.hasMore && !props.loadingMore) {
					emit("loadMore");
				}
			});
		},
		{
			rootMargin: props.rootMargin,
			threshold: props.threshold,
		},
	);

	observer.observe(loaderTrigger.value);
};

const destroyObserver = () => {
	if (observer && loaderTrigger.value) {
		observer.unobserve(loaderTrigger.value);
		observer.disconnect();
		observer = null;
	}
};

defineExpose({ loaderTrigger });

watch(
	() => props.autoLoad,
	(newVal) => {
		destroyObserver();
		if (newVal) {
			initObserver();
		}
	},
);

watch(
	() => props.hasMore,
	(newVal) => {
		if (!newVal) {
			destroyObserver();
		}
	},
);

onMounted(() => {
	if (props.autoLoad) {
		initObserver();
	}
});

onUnmounted(() => {
	destroyObserver();
});
</script>

<style lang="css" scoped>
.loader-trigger {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 20px;
  min-height: 60px;
}

.loading-text {
  color: #666;
  font-size: 14px;
  font-weight: 500;
}

.end-message {
  text-align: center;
  padding: 20px;
  color: #888;
  font-size: 14px;
  border-top: 1px solid #e0e0e0;
  margin-top: 16px;
}
</style>