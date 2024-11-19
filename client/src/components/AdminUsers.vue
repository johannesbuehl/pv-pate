<script lang="ts">
	export function validate_password(password: string): boolean {
		const password_length = password?.length;

		return password_length >= 12 && password_length <= 64;
	}
</script>

<script setup lang="ts">
	import { faPlus, faSdCard, faTrash } from "@fortawesome/free-solid-svg-icons";
	import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
	import { onMounted, ref } from "vue";

	import BaseButton from "./BaseButton.vue";
	import { api_call } from "@/lib";
	import type { User } from "@/Globals";

	interface PasswordUser extends User {
		password: string;
	}

	const add_user_name_input = ref<string>("");
	const add_user_password_input = ref<string>("");
	const users = ref<PasswordUser[]>([]);

	onMounted(async () => {
		const response = await api_call<User[]>("GET", "users");

		if (response.ok) {
			store_users(await response.json());
		}
	});

	function store_users(new_user: User[]) {
		users.value = new_user.map((user) => {
			return { ...user, password: "" };
		});
	}

	function validate_new_user(): boolean {
		if (add_user_name_input.value.length === 0) {
			return false;
		}

		if (!validate_password(add_user_password_input.value)) {
			return false;
		}

		return true;
	}

	async function add_user() {
		if (validate_new_user()) {
			const response = await api_call<User[]>("POST", "users", undefined, {
				name: add_user_name_input.value,
				password: add_user_password_input.value
			});

			if (response.ok) {
				store_users(await response.json());

				// clear the input-boxes
				add_user_name_input.value = "";
				add_user_password_input.value = "";
			}
		}
	}

	async function delete_user(user: PasswordUser) {
		if (user.name !== "admin") {
			if (window.confirm(`Delete user '${user.name}'?`)) {
				const response = await api_call<User[]>("DELETE", "users", { uid: user.uid });

				if (response.ok) {
					store_users(await response.json());
				}
			}
		}
	}

	async function modify_user(user: PasswordUser) {
		if (validate_password(user.password)) {
			const response = await api_call<User[]>(
				"PATCH",
				"users",
				{ uid: user.uid },
				{ password: user.password }
			);

			if (response.ok) {
				users.value = await response.json();
			}
		}
	}
</script>

<template>
	<div class="m-2 flex max-w-full flex-col items-center gap-2">
		<h1>Benutzer</h1>
		<div class="flex items-center gap-4">
			<div class="flex grid-cols-[auto_auto] flex-col sm:grid sm:gap-2">
				Benutzername:
				<input
					class="flex-1 rounded px-2 outline outline-1 invalid:text-red-500"
					type="text"
					name="username"
					autocomplete="off"
					v-model="add_user_name_input"
					placeholder="username"
				/>
				Passwort:
				<input
					class="flex-1 rounded px-2 outline outline-1 invalid:text-red-500"
					type="text"
					name="password"
					v-model="add_user_password_input"
					placeholder="password"
				/>
			</div>
			<BaseButton :disabled="!validate_new_user()" :square="true" @click="add_user">
				<FontAwesomeIcon :icon="faPlus" />
			</BaseButton>
		</div>
		<div class="max-w-full overflow-x-auto">
			<table class="max-w-160">
				<thead class="bg-black text-white">
					<tr>
						<th>UID</th>
						<th>Name</th>
						<th>Passwort</th>
						<th>Bestätigen</th>
						<th>Löschen</th>
					</tr>
				</thead>
				<tbody>
					<tr v-for="user of users" :key="user.uid" class="odd:bg-stone-300 even:bg-stone-100">
						<th>{{ user.uid }}</th>
						<th>{{ user.name }}</th>
						<th>
							<input
								class="rounded px-2 text-sm outline outline-2"
								type="text"
								name="new-password"
								v-model="user.password"
								placeholder="Neues Passwort"
							/>
						</th>
						<th>
							<BaseButton
								class="button mx-auto"
								:disabled="!validate_password(user.password)"
								:square="true"
								@click="modify_user(user)"
								><FontAwesomeIcon :icon="faSdCard"
							/></BaseButton>
						</th>
						<th>
							<BaseButton
								class="button mx-auto"
								:disabled="user.name === 'admin'"
								:square="true"
								@click="delete_user(user)"
								><FontAwesomeIcon :icon="faTrash"
							/></BaseButton>
						</th>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</template>

<style scoped>
	th {
		@apply p-1;
	}

	tbody th {
		@apply font-normal;
	}
</style>
