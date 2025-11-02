<script lang="ts" setup>
import MinecraftIcon from '../assets/images/minecraft_icon.jpg'
import { RouterLink } from 'vue-router';
</script>
       
<template>
    <div class="header">
        
        <RouterLink to="/" class="link">
            <div class="title">
                <img :src="MinecraftIcon" class="mc-logo"/>
                <p class="name text-shd">Mc Launcher</p>
                <p class="name-v2 text-shd">By Lanxre</p>
            </div>
        </RouterLink>
    </div>
</template>

<style lang="css" scoped>
.header {
    background-color: rgb(4, 88, 23);

    padding: 10px;
    margin: 0px;

    border-bottom-left-radius: 15%;
    border-bottom-right-radius: 25%;
}

.mc-logo {
    width: 55px;
    height: 55px;

    transform: rotate(-20deg);
}

.title {
    display: flex;

    flex-direction: row;
    justify-content: center;
    align-items: center;

    gap: 35px;

    cursor: pointer;
}

.name {
    display: flex;

    font-family: 'Minecraftv1';
    font-size: 50px;
}

.name-v2 {
    display: flex;

    font-family: 'Minecraftv1';
    font-size: 10px;

    color: yellow;

    margin-top: 10px;
    margin-left: -35px;

    transform: rotate(-40deg);
}

</style>