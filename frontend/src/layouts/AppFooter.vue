<script lang="ts" setup>
import GithubIcon from '@/assets/images/github_icon.png'
import TelegramIcon from '@/assets/images/telegram-4.png'

import CircleButton from '@/components/Buttons/CircleButton.vue'
import { OpenExternalLink } from '../../wailsjs/go/main/App'

const openLink = async (url: string) => {
    await OpenExternalLink(url)
}

</script>
       
<template>
    <div class="footer">
        <CircleButton bottom="70px" right="1%"/>
        <div class="link-group">
            <div class="group" @click="openLink('https://github.com/lanxre')">
                <img :src="GithubIcon" alt="github" class="icon"/>
                <p class="text" style="color: white;">GitHub</p>
            </div>

            <div class="group" @click="openLink('https://t.me/lanxreeazy')">
                <img :src="TelegramIcon" alt="github" class="icon"/>
                <p class="text" style="color: white;">Telegram</p>
                
            </div>
        </div>
    </div>
</template>

<style lang="css" scoped>
.footer {
    width: 100%;

    flex: 2;

    background-color: rgb(4, 88, 23);

    padding: 5px;
    margin:  0px;

    border-top-right-radius: 15%;
    border-top-left-radius:  25%;

    position: sticky;
    bottom: -1px;
}

.link-group {
    display: flex;
    flex-direction: row;

    justify-content: space-evenly;

}

.group {
    display: flex;
    flex-direction: column;

    align-items: center;
    cursor: pointer;
}

.group a {
    text-decoration: none;
    color: white;
}

.icon {
    width: 50px;
    height: 50px;

    display: flex;
    justify-content: center;
}

</style>