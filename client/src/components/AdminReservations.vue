<script lang="ts">
	interface Reservation {
		name: string;
		new_name: string;
		reservation: string;
		mid: string;
	}
</script>

<script setup lang="ts">
	import { api_call } from '@/lib';
import { faEuro, faSdCard, faTrash } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import { onMounted, ref, watch } from 'vue';
import BaseButton from './BaseButton.vue';
import { get_element_roof, get_element_string } from './BasePV.vue';

	const reservations = ref<Reservation[]>();

	onMounted(async () => {
		const response = await api_call<Reservation[]>("GET", "reservations");

		if (response.ok) {
			reservations.value = await response.json();
		}
	});

	// on new reservations, populate new_name
	watch(reservations, (reservations) => {
		reservations?.forEach(reservation => reservation.new_name = reservation.name);
	});

	async function confirm_reservation(mid: string) {
		if (confirm(`Reservierung für ${get_element_roof(mid)} bestätigen?`)) {
			const response = await api_call<Reservation[]>("POST", "reservations", { mid })

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

	async function update_reservation(reservation: Reservation) {
		// if name and new_name are the same, do nothing
		if (reservation.name !== reservation.new_name) {
			if (confirm(`Reservierung für ${get_element_roof(reservation.mid)} aktualiseren?\nVon "${reservation.name}" zu "${reservation.new_name}"`)) {
				const response = await api_call<Reservation[]>("PATCH", "reservations", { mid: reservation.mid }, { name: reservation.new_name });
				
				if (response.ok) {
					reservations.value = await response.json();
				}
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
				<th>Reservierungsdatum</th>
				<th>Bestätigen</th>
				<th>Löschen</th>
			</tr>
		</thead>
		<tbody>
			<tr
				v-for="reservation of reservations"
				:key="reservation.mid"
			>
				<th>{{ get_element_string(reservation.mid) }}</th>
				<th class="name-cell">
					<input type="text" v-model="reservation.new_name" @keydown.enter="update_reservation(reservation)" />
					<BaseButton
						:disabled="reservation.name === reservation.new_name"
						@click="update_reservation(reservation)"
					>
						<FontAwesomeIcon :icon="faSdCard" />
					</BaseButton>
				</th>
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

	th.name-cell {
		display: flex;
		align-items: center;

		gap: 0.25em;
	}
</style>
