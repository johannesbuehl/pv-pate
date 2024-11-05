<script lang="ts">
	enum WindowState {
		Login,
		Reservations,
		Sponsorships,
		Account,
		Users
	}
</script>

<script setup lang="ts">
	import { ref, watch } from 'vue';
	
	import AdminUsers from './components/AdminUsers.vue';
	import AdminLogin from './components/AdminLogin.vue';
	import AdminAccount from './components/AdminAccount.vue';
	import AppLayout from './components/AppLayout/AppLayout.vue';
	import AdminReservations from './components/AdminReservations.vue';
	import AdminSponsorships from './components/AdminSponsorships.vue';
	import { user } from './Globals';
	import { is_element_available } from './lib';
	import type { Element } from './components/BasePV.vue';

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
</script>

<template>
	<AppLayout>
		<template #header>
			<a class="navbar-item" :class="{ active: window_state === WindowState.Reservations }" @click="window_state = WindowState.Reservations">Reservierungen</a>
			<a class="navbar-item" :class="{ active: window_state === WindowState.Sponsorships }" @click="window_state = WindowState.Sponsorships">Patenschaften</a>
			<a class="navbar-item" :class="{ active: window_state === WindowState.Account }" @click="window_state = WindowState.Account">Account</a>
			<a v-if="user?.name === 'admin'" class="navbar-item" :class="{ active: window_state === WindowState.Users }" @click="window_state = WindowState.Users">Benutzer</a>
		</template>
		<AdminLogin v-if="window_state === WindowState.Login" v-model="user" />
		<AdminReservations v-else-if="window_state === WindowState.Reservations" />
		<AdminSponsorships v-else-if="window_state === WindowState.Sponsorships" />
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
