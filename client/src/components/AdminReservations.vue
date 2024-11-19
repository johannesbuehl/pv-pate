<script lang="ts">
	interface Reservation {
		name: string;
		new_name: string;
		reservation: string;
		mid: string;
	}
</script>

<script setup lang="ts">
	import { api_call } from "@/lib";
	import { faEuro, faSdCard, faTrash } from "@fortawesome/free-solid-svg-icons";
	import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
	import { onMounted, ref, watch } from "vue";
	import BaseButton from "./BaseButton.vue";
	import { get_element_roof, get_element_string } from "./BasePV.vue";

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
		reservations?.forEach((reservation) => (reservation.new_name = reservation.name));
	});

	async function confirm_reservation(mid: string) {
		if (confirm(`Reservierung für ${get_element_roof(mid)} bestätigen?`)) {
			const response = await api_call<Reservation[]>("POST", "reservations", { mid });

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
			if (
				confirm(
					`Reservierung für ${get_element_roof(reservation.mid)} aktualiseren?\nVon "${reservation.name}" zu "${reservation.new_name}"`
				)
			) {
				const response = await api_call<Reservation[]>(
					"PATCH",
					"reservations",
					{ mid: reservation.mid },
					{ name: reservation.new_name }
				);

				if (response.ok) {
					reservations.value = await response.json();
				}
			}
		}
	}
</script>

<template>
	<h1>Reservierungen</h1>

	<div class="max-w-full overflow-x-auto">
		<table>
			<thead class="bg-black text-white">
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
					class="odd:bg-stone-300 even:bg-stone-100"
				>
					<th>{{ get_element_string(reservation.mid) }}</th>
					<th class="flex items-center gap-1">
						<input
							class="rounded px-2 text-sm outline outline-2"
							type="text"
							name="name"
							v-model="reservation.new_name"
							@keydown.enter="update_reservation(reservation)"
							autocomplete="off "
						/>
						<BaseButton
							:disabled="reservation.name === reservation.new_name"
							:square="true"
							@click="update_reservation(reservation)"
						>
							<FontAwesomeIcon :icon="faSdCard" />
						</BaseButton>
					</th>
					<th>{{ reservation.reservation }}</th>
					<th>
						<BaseButton class="mx-auto" @click="confirm_reservation(reservation.mid)" :square="true"
							><FontAwesomeIcon :icon="faEuro"
						/></BaseButton>
					</th>
					<th>
						<BaseButton class="mx-auto" @click="delete_reservation(reservation.mid)" :square="true"
							><FontAwesomeIcon :icon="faTrash"
						/></BaseButton>
					</th>
				</tr>
			</tbody>
		</table>
	</div>
	<a ref="download_ref" style="display: none" />
</template>

<style scoped>
	th {
		@apply p-1;
	}

	tbody th {
		@apply font-normal;
	}
</style>
