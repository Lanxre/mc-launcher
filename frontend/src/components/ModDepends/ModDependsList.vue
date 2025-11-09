<script setup lang="ts">
import { onMounted, ref } from 'vue';
import ModDependCard from './ModDependCard.vue';
import type { ModDependency } from '@/types';
import { GetSavedMods } from '../../../wailsjs/go/main/FuncService';

interface Props {
    depends: ModDependency[]
}

interface PairDepend {
    configDepend: ModDependency
    fileDepend: string
}

const pairDepends = ref<PairDepend[]>([])

const props = defineProps<Props>()

onMounted(async () => {
    try {
        const savedModsOnDisk = await GetSavedMods()

        for (const depend of props.depends) {
            if (!depend?.Name) continue;
            
            const expectedFileName = depend.Name.replaceAll(" ", "_").toLowerCase()
             const matchingFile = savedModsOnDisk.find(file => 
                file.toLowerCase().startsWith(expectedFileName)
            );
    
            if (matchingFile) {
                pairDepends.value.push({
                    configDepend: depend,
                    fileDepend: matchingFile,
                });
            }
        }

    } catch (err) {
        console.error(err)
    }
})

</script>

<template>
    <div>
        <span class="downloads-title"> Зависимости </span>
        <div class="depend-list" v-for="depend in pairDepends" :key="depend.configDepend.Name">
            <ModDependCard :depend="depend.configDepend" :filename="depend.fileDepend"/>
        </div>
    </div>
</template>

<style lang="css" scoped>

.depened-list {
  display: flex;
  flex-direction: column;
  gap: 5px;
  
}

.downloads-title {
  display: flex;
  justify-content: center;
  font-size: 30px;
  color: white;
  padding: 10px;
}

</style>