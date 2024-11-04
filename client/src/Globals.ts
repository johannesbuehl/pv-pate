import { ref } from "vue";
import { api_call, HTTPStatus } from "./lib";

export interface ElementsDB {
	taken: Record<string, string>;
	reserved: string[];
}

export const elements_db = ref<ElementsDB>({
	reserved: [],
	taken: {}
});

void (async () => {
	const reserved_elements_request = api_call<ElementsDB>("GET", "elements");

	if ((await reserved_elements_request).ok) {
		elements_db.value = (await (await reserved_elements_request).json())
	}
})();

export interface User {
	uid: number;
	name: string;
}

export interface UserLogin extends User {
	logged_in: boolean;
}

export const user = ref<UserLogin>();

void (async () => {
	const response = await api_call<UserLogin>("GET", "welcome");

	if (response.status === HTTPStatus.OK) {
		user.value = await response.json();
	}
})()