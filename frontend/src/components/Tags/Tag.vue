<script lang="ts" setup>
const props = defineProps({
	name: {
		type: String,
		required: true,
	},
	color: {
		type: String,
		default: "grey",
		validator: (value: string) =>
			["grey", "blue", "red", "green", "yellow", "purple"].includes(value),
	},
	clickable: {
		type: Boolean,
		default: true,
	},
	closable: {
		type: Boolean,
		default: false,
	},
	size: {
		type: String,
		default: "medium",
		validator: (value: string) => ["small", "medium", "large"].includes(value),
	},
});

const emit = defineEmits(["click", "close"]);

const handleClick = (event: Event) => {
	if (props.clickable) {
		emit("click", event);
	}
};

const handleClose = () => {
	emit("close");
};
</script>


<template>
    <div 
        class="tag" 
        :class="[
            `tag--${size}`,
            `tag--${color}`,
            { 
                'tag--clickable': clickable,
                'tag--closable': closable
            }
        ]"
        role="button"
        :tabindex="clickable ? 0 : -1"
        @click="handleClick"
        @keydown.enter="handleClick"
        @keydown.space="handleClick"
    >
        <span class="tag__text">{{ name }}</span>
        <button 
            v-if="closable" 
            class="tag__close" 
            @click.stop="handleClose"
            aria-label="Удалить тег"
        >
            ×
        </button>
    </div>
</template>

<style lang="css" scoped>
.tag {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    
    min-width: fit-content;
    height: fit-content;
    
    background-color: var(--tag-bg, #6b7280);
    color: white;
    
    padding: 6px 12px;
    margin: 2px;
    
    border-radius: 16px;
    border: none;
    
    font-family: inherit;
    font-size: 0.875rem;
    font-weight: 500;
    line-height: 1;
    
    transition: all 0.2s ease-in-out;
    
    user-select: none;
}

.tag--clickable {
    cursor: pointer;
}

.tag--clickable:hover {
    background-color: var(--tag-bg-hover);
    transform: translateY(-1px);
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.tag--clickable:active {
    transform: translateY(0);
}

.tag--clickable:focus-visible {
    outline: 2px solid #3b82f6;
    outline-offset: 1px;
}

.tag__text {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 120px;
}

.tag__close {
    display: flex;
    align-items: center;
    justify-content: center;
    
    width: 16px;
    height: 16px;
    
    background: rgba(0, 0, 0, 0.2);
    color: white;
    border: none;
    border-radius: 50%;
    
    font-size: 14px;
    line-height: 1;
    
    cursor: pointer;
    transition: background-color 0.2s ease;
}

.tag__close:hover {
    background: rgba(0, 0, 0, 0.3);
}

/* Размеры */
.tag--small {
    padding: 4px 8px;
    font-size: 0.75rem;
}

.tag--large {
    padding: 8px 16px;
    font-size: 1rem;
}

/* Цветовые варианты */
.tag--grey {
    --tag-bg: #6b7280;
    --tag-bg-hover: #4b5563;
}

.tag--blue {
    --tag-bg: #3b82f6;
    --tag-bg-hover: #2563eb;
}

.tag--red {
    --tag-bg: #ef4444;
    --tag-bg-hover: #dc2626;
}

.tag--green {
    --tag-bg: #10b981;
    --tag-bg-hover: #059669;
}

.tag--yellow {
    --tag-bg: #f59e0b;
    --tag-bg-hover: #d97706;
}

.tag--purple {
    --tag-bg: #8b5cf6;
    --tag-bg-hover: #7c3aed;
}
</style>