<script setup lang="ts">
import { onMounted, ref } from "vue";
import ModDependCard from "./ModDependCard.vue";
import type { ModDependency, PairDepend } from "@/types";
import { filterDiskModDepends } from "@/api/utils";
import { DeleteSavedMod } from "@wailsjs/go/functools/FuncService";

interface Props {
	depends: ModDependency[];
}

const props = defineProps<Props>();
const pairDepends = ref<PairDepend[]>([]);

const removeDepend = async (depend: ModDependency, filename: string) => {
	pairDepends.value = pairDepends.value.filter(
		(pd) => pd.configDepend.Name !== depend.Name,
	);
	await DeleteSavedMod(filename);
};

onMounted(async () => {
	try {
		const filtred = await filterDiskModDepends(props.depends);
		pairDepends.value = filtred;
	} catch (err) {
		console.error(err);
	}
});
</script>

<template>
    <div>
        <span class="downloads-title"> Зависимости </span>
        <div class="depend-list" v-for="depend in pairDepends" :key="depend.configDepend.Name">
            <ModDependCard :depend="depend.configDepend" :filename="depend.fileDepend" v-on:delete="removeDepend(depend.configDepend, depend.fileDepend)"/>
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