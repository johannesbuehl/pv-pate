<script lang="ts">
	enum WindowState {
		Login,
		Reservations,
		Elements,
		Account,
		Users
	}
</script>

<script setup lang="ts">
	import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
	import { faSdCard, faTrashCan, faXmark } from '@fortawesome/free-solid-svg-icons';
	import { ref, watch } from 'vue';
	
	import BasePV, { get_element_roof, get_element_type, type Element } from './components/BasePV.vue';
	import { api_call, get_element, is_element_available, type APICallResult } from './lib';
	import { elements_db, user, type ElementsDB } from './Globals';
	import AdminUsers from './components/AdminUsers.vue';
	import BaseButton from './components/BaseButton.vue';
	import AdminLogin from './components/AdminLogin.vue';
	import AdminAccount from './components/AdminAccount.vue';
	import AppLayout from './components/AppLayout/AppLayout.vue';
import AdminReservations from './components/AdminReservations.vue';

	const window_state = ref<WindowState>(WindowState.Login);
	const selected_element = ref<Element>();

	watch(user, user => {
			window_state.value = user?.logged_in ? WindowState.Reservations : WindowState.Login
	}, { deep: true });

	watch(selected_element, () => {
		if (!!selected_element.value && is_element_available(selected_element.value?.mid)) {
			console.debug("is free")

			selected_element.value = undefined;
		}
	});

	async function submit() {
		// check wether a element is selected
		if (selected_element.value !== undefined) {
			let response: APICallResult<ElementsDB>;
			
			const name = selected_element.value.name !== "" ? selected_element.value.name : undefined;
				
			// if the element is already reserved, patch it instead
			const method = is_element_available(selected_element.value.mid) ? "POST" : "PATCH";
				
			response = await api_call<ElementsDB>(method, "elements", { mid: selected_element.value.mid }, {
				name
			});
			
			if (response.ok) {
				elements_db.value = (await response.json());
				
				selected_element.value = undefined;
			} else {
				alert(`Error during database write: ${await response.text()}`);
			}
		}
	}

	async function delete_reservation() {
		// only proceed if there is a valid element-selection
		if (selected_element.value !== undefined) {
			const element_name = selected_element.value?.mid.match(/\w?\d+/)?.["0"].toUpperCase()
			
			if (confirm(`Patenschaft für ${get_element_type(selected_element.value.mid)} ${element_name} mit dem Namen "${get_element(selected_element.value.mid)?.name}" löschen?`)) {
				const response = await api_call<ElementsDB>("DELETE", "elements", { mid: selected_element.value.mid });
				
				if (response.ok) {
					elements_db.value = (await response.json());

					selected_element.value.name = undefined;
				}
			}
		}
	}
</script>

<template>
	<AppLayout>
		<template #header>
			<a class="navbar-item" :class="{ active: window_state === WindowState.Reservations }" @click="window_state = WindowState.Reservations">Reservierungen</a>
			<a class="navbar-item" :class="{ active: window_state === WindowState.Elements }" @click="window_state = WindowState.Elements">Elemente</a>
			<a class="navbar-item" :class="{ active: window_state === WindowState.Account }" @click="window_state = WindowState.Account">Account</a>
			<a v-if="user?.name === 'admin'" class="navbar-item" :class="{ active: window_state === WindowState.Users }" @click="window_state = WindowState.Users">Benutzer</a>
		</template>
		<AdminLogin v-if="window_state === WindowState.Login" v-model="user" />
		<AdminReservations v-else-if="window_state === WindowState.Reservations" />
		<BasePV
			v-else-if="window_state === WindowState.Elements"
			v-model:selected_element="selected_element"
		>
			<template #header
				v-if="selected_element !== undefined"
			>
				{{ get_element_roof(selected_element.mid) }}
			</template>
			<div
				v-if="selected_element"
				id="tooltip-content"
			>
				<div
					v-if="selected_element.reserved"
				>
					Modul ist reserviert.<br>
					<a @click="window_state = WindowState.Reservations">
						Zu Reservierungen springen
					</a>
				</div>
				<template v-else-if="!is_element_available(selected_element.mid)">
					<BaseButton @click="submit">
						<FontAwesomeIcon :icon="faSdCard" />
					</BaseButton>
					<input type="text" id="input-name" v-model="selected_element.name" placeholder="Anonym" @keydown.enter="submit" />
					<BaseButton
						v-if="is_element_available(selected_element.mid)"
						@click="selected_element.name = undefined"
					>
						<FontAwesomeIcon :icon="faXmark" />
					</BaseButton>
					<BaseButton
						v-else
						@click="delete_reservation"
					>
						<FontAwesomeIcon :icon="faTrashCan" />
					</BaseButton>
				</template>
			</div>
		</BasePV>
		<AdminAccount v-else-if="window_state === WindowState.Account" />
		<AdminUsers v-else-if="window_state === WindowState.Users" />
	</AppLayout>
</template>

<style scoped>
	#main-view {
		width: 100%;
		height: 100%;
	}

	.navbar-item.active {
		text-decoration: underline;

		font-weight: bold;
	}

	#tooltip-content {
		display: flex;
		
		align-items: center;

		gap: 0.25em;
	}

	#input-name {
		transition: opacity 0.2s;

		width: 100%;
	}

	#input-name:disabled {
		cursor: not-allowed;
		opacity: 50%;
	}

	#tooltip-content a {
		text-decoration: underline;
		font-style: italic;
	}
</style>
