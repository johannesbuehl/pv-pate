<script lang="ts">
	export function validate_password(
		password_new: string,
		password_repeat?: string,
		password_old?: string
	): string[] {
		const result: string[] = [];

		if (password_old !== undefined && password_old.length === 0) {
			result.push("Bisheriges Passwort fehlt");
		}

		if (password_new.length < 12) {
			result.push("Passwort muss mindestens 12 Zeichen lang sein");
		}

		if (password_new.length > 64) {
			result.push("Passwort darf höchstens 64 Zeichen lang sein");
		}

		if (password_repeat !== undefined && password_new !== password_repeat) {
			result.push("Passwörter stimmen nicht überein");
		}

		return result;
	}
</script>

<script setup lang="ts">
	import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
	import BaseButton from "./BaseButton.vue";
	import { faSdCard } from "@fortawesome/free-solid-svg-icons";
	import { ref } from "vue";

	import { api_call } from "@/lib";

	const password_current = ref<string>("");
	const password_new = ref<string>("");
	const password_repeat = ref<string>("");

	async function change_password() {
		if (
			validate_password(password_new.value, password_repeat.value, password_current.value)
				.length === 0
		) {
			const response = await api_call<{}>("PATCH", "user/password", undefined, {
				password: password_new.value
			});

			if (response.ok) {
				alert("Passwort erfolgreich geändert");
				password_current.value = "";
				password_new.value = "";
				password_repeat.value = "";

				// reload the page since the users get's logged out
				location.reload();
			}
		}
	}
</script>

<template>
	<div class="m-2 flex flex-col items-center">
		<h1>Account</h1>
		<div class="flex flex-col items-center gap-4">
			<form class="flex grid-cols-[auto_auto] flex-col sm:grid sm:gap-2">
				<input style="display: none" type="text" name="username" autocomplete="username" />
				Bisheriges Passwort
				<input
					type="password"
					class="flex-1 rounded px-2 outline outline-2 invalid:text-red-500"
					name="current-password"
					autocomplete="current-password"
					v-model="password_current"
				/>
				Neues Passwort
				<input
					class="flex-1 rounded px-2 outline outline-2 invalid:text-red-500"
					type="password"
					name="new-password"
					autocomplete="new-password"
					v-model="password_new"
				/>
				Neues Passwort wiederholen
				<input
					class="flex-1 rounded px-2 outline outline-2 invalid:text-red-500"
					type="password"
					name="new-password-confirm"
					autocomplete="new-password"
					v-model="password_repeat"
				/>
			</form>
			<div
				v-if="validate_password(password_new, password_repeat, password_current).length > 0"
				class="text-red-500"
			>
				<div
					v-for="e in validate_password(password_new, password_repeat, password_current)"
					:key="e"
				>
					{{ e }}
				</div>
			</div>
			<BaseButton
				:disabled="validate_password(password_new, password_repeat, password_current).length > 0"
				@click="change_password"
				><FontAwesomeIcon :icon="faSdCard" /> Passwort ändern</BaseButton
			>
		</div>
	</div>
</template>

<style scoped></style>
