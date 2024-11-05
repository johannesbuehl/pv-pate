<script lang="ts">
	interface Reservation {
		name: string;
		reservation: string;
		mid: string;
	}
</script>

<script setup lang="ts">
	import { api_call } from '@/lib';
import { faEuro, faTrash } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import { onMounted, ref } from 'vue';
import BaseButton from './BaseButton.vue';
import { get_element_roof, get_element_string } from './BasePV.vue';

	const reservations = ref<Reservation[]>();

	onMounted(async () => {
		const response = await api_call<Reservation[]>("GET", "reservations");

		if (response.ok) {
			reservations.value = await response.json();
		}
	});

	async function confirm_reservation(mid: string) {
		if (confirm(`Reservierung für ${get_element_roof(mid)} bestätigen?`)) {
			const response = await api_call<Reservation[]>("POST", "reservations", { q: "confirm", mid })

			if (response.ok) {
				reservations.value = await response.json();
			}
		}
	}

	async function delete_reservation(mid: string) {
		if (confirm(`Reservierung für ${get_element_roof(mid)} löschen?`)) {
			const response = await api_call<Reservation[]>("DELETE", "reservations", { mid });
			
			if (response.ok) {
				reservations.value = await response.json();
			}
		}
	}
</script>

<template>
	<h1>Reservierungen</h1>

	<table>
		<thead>
			<tr>
				<th>Element</th>
				<th>Name</th>
				<th>Date</th>
				<th>Confirm</th>
				<th>Delete</th>
			</tr>
		</thead>
		<tbody>
			<tr
				v-for="reservation of reservations"
				:key="reservation.mid"
			>
				<th>{{ get_element_string(reservation.mid) }}</th>
				<th>{{ reservation.name }}</th>
				<th>{{ reservation.reservation }}</th>
				<th><BaseButton @click="confirm_reservation(reservation.mid)"><FontAwesomeIcon :icon="faEuro" /></BaseButton></th>
				<th><BaseButton @click="delete_reservation(reservation.mid)"><FontAwesomeIcon :icon="faTrash" /></BaseButton></th>
			</tr>
		</tbody>
	</table>
</template>

<style scoped>
	thead th {
		font-weight: bold;

		background-color: black;
		color: white;
	}

	tbody > tr:nth-of-type(2n) {
		background-color: hsl(0, 0%, 90%);
	}

	tbody > tr:nth-of-type(2n + 1) {
		background-color: hsl(0, 0%, 80%);
	}

	th {
		padding: 0.25em;

		font-weight: normal;
	}

	th > div.cell {
		width: 100%;

		display: flex;
		align-items: center;
		justify-content: center;
	}

	th input[type="text"] {
		flex: 1;
	}

	tbody > tr input[type="text"] {
		font-size: 0.67em;
	}
</style>
