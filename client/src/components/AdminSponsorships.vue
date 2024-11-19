<script lang="ts">
	interface Sponsorship {
		mid: string;
		name: string;
		new_name: string;
	}
</script>

<script setup lang="ts">
	import { api_call } from "@/lib";
	import { onMounted, ref, watch } from "vue";
	import { get_element_roof, get_element_string } from "./BasePV.vue";
	import { faDownload, faSdCard, faTrash } from "@fortawesome/free-solid-svg-icons";
	import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
	import BaseButton from "./BaseButton.vue";

	const sponsorships = ref<Sponsorship[]>();

	onMounted(async () => {
		const response = await api_call<Sponsorship[]>("GET", "sponsorships");

		if (response.ok) {
			sponsorships.value = await response.json();
		}
	});

	// on new sponsorships, populate new_name
	watch(sponsorships, (sponsorships) => {
		sponsorships?.forEach((sponsorship) => (sponsorship.new_name = sponsorship.name));
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
			if (
				confirm(
					`Patenschaft für ${get_element_roof(sponsorship.mid)} aktualiseren?\nVon "${sponsorship.name}" zu "${sponsorship.new_name}"`
				)
			) {
				const response = await api_call<Sponsorship[]>(
					"PATCH",
					"sponsorships",
					{ mid: sponsorship.mid },
					{ name: sponsorship.new_name }
				);

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

	<div class="max-w-full overflow-x-auto">
		<table>
			<thead>
				<tr class="bg-black text-white">
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
					class="odd:bg-stone-300 even:bg-stone-100"
				>
					<th>{{ get_element_string(sponsorship.mid) }}</th>
					<th class="flex items-center gap-1">
						<input
							class="rounded px-2 text-sm outline outline-2"
							type="text"
							name="name"
							v-model="sponsorship.new_name"
							@keydown.enter="update_sponsorship(sponsorship)"
							autocomplete="off"
						/>
						<BaseButton
							:disabled="sponsorship.name === sponsorship.new_name"
							:square="true"
							@click="update_sponsorship(sponsorship)"
						>
							<FontAwesomeIcon :icon="faSdCard" />
						</BaseButton>
					</th>
					<th class="mx-auto">
						<form action="/api/certificates" target="_blank">
							<input type="text" name="mid" style="display: none" :value="sponsorship.mid" />
							<input
								type="submit"
								:id="`get-certificate-submit-${sponsorship.mid}`"
								style="display: none"
							/>
							<label :for="`get-certificate-submit-${sponsorship.mid}`">
								<BaseButton class="mx-auto" @click="get_certificate(sponsorship.mid)" :square="true"
									><FontAwesomeIcon :icon="faDownload"
								/></BaseButton>
							</label>
						</form>
					</th>
					<th>
						<BaseButton class="mx-auto" @click="delete_sponsorship(sponsorship.mid)" :square="true"
							><FontAwesomeIcon :icon="faTrash"
						/></BaseButton>
					</th>
				</tr>
			</tbody>
		</table>
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
