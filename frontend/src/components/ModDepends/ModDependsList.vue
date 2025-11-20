<script setup lang="ts">
import { DeleteSavedMod } from "@wailsjs/go/functools/FuncService";
import { computed, onMounted, ref } from "vue";
import { filterDiskModDepends } from "@/api/utils";
import type { ModDependency, PairDepend } from "@/types";
import ModDependCard from "./ModDependCard.vue";

interface Props {
	depends: ModDependency[];
}

const props = defineProps<Props>();
const pairDepends = ref<PairDepend[]>([]);

const notDowloadDepends = computed(() => {
	return props.depends.filter(
		(dep) =>
			!pairDepends.value
				.flatMap((pd) => pd.configDepend.Name)
				.includes(dep.Name),
	);
});

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
        <span class="downloads-title"> Dependency </span>
        <div class="depend-list" v-if="pairDepends.length > 0" v-for="depend in pairDepends" :key="depend.configDepend.Name">
          <ModDependCard :depend="depend.configDepend" :filename="depend.fileDepend" v-on:delete="removeDepend(depend.configDepend, depend.fileDepend)"/>
        </div>
        <div class="depend-list" v-for="depend in notDowloadDepends">
          <ModDependCard :depend="depend" :filename="depend.Name"/>
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
  font-size: 45px;
  color: white;
  padding: 10px;
  font-family: "MinecraftV1";
}

</style>