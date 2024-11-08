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
	const download_ref = ref<HTMLAnchorElement>();

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
			const response = await api_call<Reservation[]>("POST", "reservations", { mid });

			if (response.ok) {
				reservations.value = await response.json();

				const urlsearchparams = new URLSearchParams({ mid });

				// request the certificate
				const response_cert = await fetch(window.origin + "/pv/api/certificates?" + urlsearchparams.toString(), {
					headers: {
						// eslint-disable-next-line @typescript-eslint/naming-convention
						"Content-Type": "application/json; charset=UTF-8"
					},
					credentials: "include",
					method: "GET"
				});

				if (response_cert.ok) {
					if (download_ref.value !== undefined) {
						// set the file-name
						const header = response_cert.headers.get('Content-Disposition');
						const parts = header?.split(';');
						const filename = parts?.[1].split('=')[1]?.replaceAll("\"", "");

						if (filename !== undefined) {
							download_ref.value.download = filename;
						}

						const url = window.URL.createObjectURL(await response_cert.blob());

						download_ref.value.href = url;
						download_ref.value.click();
						window.URL.revokeObjectURL(url);
					}
				}
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

	<div id="table-wrapper">
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
						<input type="text" name="name" v-model="reservation.new_name" @keydown.enter="update_reservation(reservation)" autocomplete="off "/>
						<BaseButton
							:disabled="reservation.name === reservation.new_name"
							:square="true" 
							@click="update_reservation(reservation)"
						>
							<FontAwesomeIcon :icon="faSdCard" />
						</BaseButton>
					</th>
					<th>{{ reservation.reservation }}</th>
					<th><BaseButton class="center" @click="confirm_reservation(reservation.mid)" :square="true"><FontAwesomeIcon :icon="faEuro" /></BaseButton></th>
					<th><BaseButton class="center" @click="delete_reservation(reservation.mid)" :square="true"><FontAwesomeIcon :icon="faTrash" /></BaseButton></th>
				</tr>
			</tbody>
		</table>
	</div>
	<a ref="download_ref" style="display: none;" />
</template>

<style scoped>
	#table-wrapper {
		overflow-x: auto;
		max-width: 100%;
	}

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

	.center {
		margin-inline: auto;
	}
</style>
