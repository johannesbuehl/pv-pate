<script lang="ts">
	interface Sponsorship {
		mid: string;
		name: string;
		new_name: string;
	}
</script>

<script setup lang="ts">
	import { api_call } from '@/lib';
	import { onMounted, ref, watch } from 'vue';
import { get_element_roof, get_element_string } from './BasePV.vue';
import { faDownload, faSdCard, faTrash } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import BaseButton from './BaseButton.vue';

	const sponsorships = ref<Sponsorship[]>();

	onMounted(async () => {
		const response = await api_call<Sponsorship[]>("GET", "sponsorships");

		if (response.ok) {
			sponsorships.value = await response.json()
		}
	});

	// on new sponsorships, populate new_name
	watch(sponsorships, (sponsorships) => {
		sponsorships?.forEach(sponsorship => sponsorship.new_name = sponsorship.name);
	});

	async function delete_sponsorship(mid: string) {
		if (confirm(`Patenschaft für ${get_element_roof(mid)} löschen?`)) {
			const response = await api_call<Sponsorship[]>("DELETE", "sponsorships", { mid });
			
			if (response.ok) {
				sponsorships.value = await response.json();
			}
		}
	}

	async function update_sponsorship(sponsorship: Sponsorship) {
		// if name and new_name are the same, do nothing
		if (sponsorship.name !== sponsorship.new_name) {
			if (confirm(`Patenschaft für ${get_element_roof(sponsorship.mid)} aktualiseren?\nVon "${sponsorship.name}" zu "${sponsorship.new_name}"`)) {
				const response = await api_call<Sponsorship[]>("PATCH", "sponsorships", { mid: sponsorship.mid }, { name: sponsorship.new_name });
				
				if (response.ok) {
					sponsorships.value = await response.json();
				}
			}
		}
	}

	async function get_certificate(mid: string) {
		// const response = await api_call()
	}
</script>

<template>
	<h1>Patenschaften</h1>

	<table>
		<thead>
			<tr>
				<th>Element</th>
				<th>Name</th>
				<th>Zertifikat</th>
				<th>Löschen</th>
			</tr>
		</thead>
		<tbody>
			<tr
				v-for="sponsorship of sponsorships"
				:key="sponsorship.mid"
			>
				<th>{{ get_element_string(sponsorship.mid) }}</th>
				<th class="name-cell">
					<input type="text" v-model="sponsorship.new_name" @keydown.enter="update_sponsorship(sponsorship)" />
					<BaseButton
						:disabled="sponsorship.name === sponsorship.new_name"
						@click="update_sponsorship(sponsorship)"
					>
						<FontAwesomeIcon :icon="faSdCard" />
					</BaseButton>
				</th>
				<th>
					<form
						action="/pv/api/certificates"
						target="_blank"
					>
						<input type="text" name="mid" style="display: none;" :value="sponsorship.mid" />
						<input type="submit" :id="`get-certificate-submit-${sponsorship.mid}`" style="display: none;" />
						<label :for="`get-certificate-submit-${sponsorship.mid}`">
							<BaseButton @click="get_certificate(sponsorship.mid)"><FontAwesomeIcon :icon="faDownload" /></BaseButton>
						</label>
					</form>
				</th>
				<th><BaseButton @click="delete_sponsorship(sponsorship.mid)"><FontAwesomeIcon :icon="faTrash" /></BaseButton></th>
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
