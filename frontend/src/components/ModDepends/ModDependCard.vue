<script setup lang="ts">
import { openLink } from "@/api/utils";
import CrossIcon from "@/assets/images/close.png";
import type { ModDependency } from "@/types";
import Image from "../Image/Image.vue";

interface Props {
	depend?: ModDependency;
	filename?: string;
	onDelete?: () => void;
}

const props = defineProps<Props>();

const openModPage = () => {
	if (props.depend !== undefined && props.depend !== null) {
		openLink(props.depend.ModPageLink);
	}
};
</script>

<template>
     <div class="depend-card">
        <div class="depend-name">
            <div class="d-name" title="Открыть в браузере" @click="openModPage">
                <div v-if="depend !== undefined">
                    <p v-if="onDelete"> {{ depend.Name }} </p>
                    <div v-else class="d-cl"> 
                        <p>{{ depend.Name }}</p>
                        <span class="no-d-text">Зависимость не скачана</span>
                    </div>
                </div> 
                <div v-if="filename && onDelete">
                    ({{ filename }})
                </div>
            </div>
        </div>
        <Image v-if="onDelete" :img="CrossIcon" width="25px" height="25px" border-raduis="50%" title="Удалить" @click="onDelete"/>            
    </div>
</template>

<style lang="css" scoped>
.depend-card {
  display: flex;
  flex-direction: row;

  gap: 10px;
  align-items: center;

  padding: 10px;

  color: white;
  border: 1px solid white;
  border-radius: 30px;
  margin-bottom: 10px;
}

.depend-name {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;

  gap: 5px;
}

.d-name {
    display: flex;
    flex-direction: column;
    align-items: center;
    cursor: pointer;
    color: white;
}

.d-cl {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
}

.no-d-text {
    color: rgb(165, 72, 72);
}

</style>