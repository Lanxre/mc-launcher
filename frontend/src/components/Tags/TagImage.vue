<script lang="ts" setup>
import { computed } from 'vue'

const props = defineProps({
    name: {
        type: String,
        required: true
    },
    imageUrl: {
        type: String,
        default: ''
    },
    altText: {
        type: String,
        default: ''
    },
    clickable: {
        type: Boolean,
        default: true
    },
    closable: {
        type: Boolean,
        default: false
    },
    size: {
        type: String,
        default: 'medium',
        validator: (value: string) => ['small', 'medium', 'large'].includes(value)
    },
    variant: {
        type: String,
        default: 'default',
        validator: (value: string) => ['default', 'outlined', 'ghost'].includes(value)
    },
    rounded: {
        type: Boolean,
        default: true
    }
})

const emit = defineEmits(['click', 'close', 'image-error'])

const fallbackText = computed(() => {
    return props.name.charAt(0).toUpperCase()
})

const handleClick = (event: Event) => {
    if (props.clickable) {
        emit('click', event)
    }
}

const handleClose = () => {
    emit('close')
}

const handleImageError = (error: Event) => {
    emit('image-error', error)
}
</script>

<template>
    <div 
        class="tag-image"
        :class="[
            `tag-image--${size}`,
            `tag-image--${variant}`,
            { 
                'tag-image--clickable': clickable,
                'tag-image--rounded': rounded
            }
        ]"
        role="button"
        :tabindex="clickable ? 0 : -1"
        @click="handleClick"
        @keydown.enter="handleClick"
        @keydown.space="handleClick"
    >
        <div class="tag-image__avatar">
            <img 
                v-if="imageUrl"
                :src="imageUrl" 
                :alt="altText"
                class="tag-image__img"
                @error="handleImageError"
            >
            <div v-else class="tag-image__fallback">
                {{ fallbackText }}
            </div>
        </div>
        
        <span class="tag-image__text">{{ name }}</span>
        
        <button 
            v-if="closable" 
            class="tag-image__close" 
            @click.stop="handleClose"
            aria-label="Удалить тег"
        >
            ×
        </button>
    </div>
</template>

<style lang="css" scoped>
.tag-image {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    
    min-width: fit-content;
    height: fit-content;
    
    background-color: #3e454d;
    color: #e9e9e9;
    border: 1px solid #e2e8f0;
    
    padding: 3px 8px 3px 3px;
    margin: 2px;
    
    border-radius: 20px;
    
    font-family: inherit;
    font-size: 0.875rem;
    font-weight: 500;
    line-height: 1;
    
    transition: all 0.2s ease-in-out;
    
    user-select: none;
}

.tag-image--clickable {
    cursor: pointer;
}

.tag-image--clickable:hover {
    background-color: #56626e;
    border-color: #cbd5e1;
    transform: translateY(-1px);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.tag-image--clickable:active {
    transform: translateY(0);
}

.tag-image--clickable:focus-visible {
    outline: 2px solid #3b82f6;
    outline-offset: 1px;
}

/* Avatar */
.tag-image__avatar {
    display: flex;
    align-items: center;
    justify-content: center;
    
    width: 20px;
    height: 20px;
    
    border-radius: 50%;
    overflow: hidden;
    flex-shrink: 0;
}

.tag-image--rounded .tag-image__avatar {
    border-radius: 6px;
}

.tag-image__img {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.tag-image__fallback {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    
    background-color: #3b82f6;
    color: white;
    font-weight: 600;
    font-size: 0.75rem;
}

/* Text */
.tag-image__text {
    white-space: nowrap;
    text-overflow: ellipsis;
    max-width: 120px;
}

/* Close button */
.tag-image__close {
    display: flex;
    align-items: center;
    justify-content: center;
    
    width: 16px;
    height: 16px;
    
    background: rgba(0, 0, 0, 0.1);
    color: #64748b;
    border: none;
    border-radius: 50%;
    
    font-size: 14px;
    line-height: 1;
    
    cursor: pointer;
    transition: all 0.2s ease;
    
    margin-left: 4px;
}

.tag-image__close:hover {
    background: rgba(0, 0, 0, 0.2);
    color: #475569;
}

/* Sizes */
.tag-image--small {
    padding: 4px 8px 4px 4px;
    font-size: 0.75rem;
    
    .tag-image__avatar {
        width: 20px;
        height: 20px;
    }
    
    .tag-image__fallback {
        font-size: 0.625rem;
    }
}

.tag-image--large {
    padding: 8px 16px 8px 8px;
    font-size: 1rem;
    
    .tag-image__avatar {
        width: 28px;
        height: 28px;
    }
    
    .tag-image__fallback {
        font-size: 0.875rem;
    }
}

/* Variants */
.tag-image--outlined {
    background-color: transparent;
    border: 1px solid #3b82f6;
    color: #3b82f6;
}

.tag-image--ghost {
    background-color: transparent;
    border: 1px solid transparent;
    
    &:hover {
        background-color: #f1f5f9;
    }
}

/* Rounded variant */
.tag-image--rounded {
    border-radius: 12px;
    
    .tag-image__avatar {
        border-radius: 6px;
    }
}
</style>