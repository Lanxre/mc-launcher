import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";

export function useScrollManager() {
	const router = useRouter();
	const scrollPosition = ref(0);

	const saveMainPageScroll = () => {
		if (router.currentRoute.value.path === "/") {
			scrollPosition.value = window.scrollY;
			sessionStorage.setItem("mainPageScroll", scrollPosition.value.toString());
		}
	};

	const restoreMainPageScroll = () => {
		if (router.currentRoute.value.path === "/") {
			const savedPosition = sessionStorage.getItem("mainPageScroll");
			if (savedPosition) {
				setTimeout(() => {
					window.scrollTo({
						top: parseInt(savedPosition),
						behavior: "auto",
					});
					sessionStorage.removeItem("mainPageScroll");
				}, 50);
			}
		}
	};

	const saveScrollBeforeLeave = () => {
		if (router.currentRoute.value.path === "/") {
			scrollPosition.value = window.scrollY;
			sessionStorage.setItem("mainPageScroll", scrollPosition.value.toString());
		}
	};

	onMounted(() => {
		if (router.currentRoute.value.path === "/") {
			restoreMainPageScroll();
		}
	});

	return {
		saveMainPageScroll,
		restoreMainPageScroll,
		saveScrollBeforeLeave,
	};
}
