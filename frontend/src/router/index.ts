import { createRouter, createWebHistory } from "vue-router";

import HomeView from "@/views/HomeView.vue";
import ModView from "@/views/ModView.vue";
import DownloadsView from "@/views/DownloadsView.vue";
import ModFavourites from "@/views/ModFavourites.vue";

const routes = [
	{
		path: "/",
		name: "main",
		component: HomeView,
	},
	{
		path: "/:name",
		name: "mod",
		component: ModView,
	},
	{
		path: "/mod-downloads",
		name: "mod-downloads",
		component: DownloadsView,
	},
	{
		path: "/mod-favourites",
		name: "mod-favourites",
		component: ModFavourites,
	},
];

const router = createRouter({
	history: createWebHistory(),
	routes,
});

export default router;
