<script setup lang="ts">
	import { ref } from "vue";

	import BasePV, {
		get_element_roof,
		get_element_type,
		type Element
	} from "./components/BasePV.vue";
	import AppLayout from "./components/AppLayout/AppLayout.vue";
	import BaseButton from "./components/BaseButton.vue";
	import { faCheck } from "@fortawesome/free-solid-svg-icons";
	import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
	import { api_call, is_element_available, type APICallResult } from "./lib";
	import { elements_db, type ElementsDB } from "./Globals";
	import BaseImageCredit from "./components/BaseImageCredit.vue";
	import BaseImageTitle from "./components/BaseImageTitle.vue";

	const reserve_form = ref<HTMLFormElement>();

	const selected_element = ref<Element & { email: string }>();

	let enter_press: boolean = false;
	async function submit(e: Event) {
		e.preventDefault();

		// if the form was submitted through the enter-key, let the user confirm
		if (enter_press) {
			enter_press = false;

			let confirm_string = `Reservierung bestätigen?\nE-Mail: ${selected_element.value?.email}\n`;

			if (
				selected_element.value?.name != undefined &&
				!!selected_element.value.name &&
				selected_element.value?.name?.length > 1
			) {
				confirm_string += `Name: ${selected_element.value?.name}`;
			} else {
				confirm_string += `Ohne Namen`;
			}

			if (!confirm(confirm_string)) {
				return;
			}
		}

		// check wether a element is selected
		if (selected_element.value !== undefined) {
			let response: APICallResult<ElementsDB>;

			// if the element is already reserved, patch it instead
			const method = is_element_available(selected_element.value.mid) ? "POST" : "PATCH";

			response = await api_call<ElementsDB>(
				method,
				"elements",
				{ mid: selected_element.value.mid },
				{
					name: selected_element.value.name,
					mail: selected_element.value.email
				}
			);

			if (response.ok) {
				elements_db.value = await response.json();

				selected_element.value = undefined;
			} else {
				alert(`Error during database write: ${await response.text()}`);
			}
		}
	}
</script>

<template>
	<AppLayout>
		<h1 class="text-center uppercase text-red-500">
			Dies ist eine interne Test-Seite. Die hier aufgelisteten Informationen sind im stetigen Wandel
			und entsprechen möglicherweise nicht der Wahrheit. Tätigen Sie keine Spenden im unten
			aufgeführten Kontext.
		</h1>
		<div id="text-content">
			<img
				id="header"
				class="mx-auto mb-8 mt-5 w-[30em] max-w-full"
				src="/external/header_logo.svg"
				alt='Logo der evangelischen Kirchengemeinde Bühl "gemeinsam glauben leben" und dem Titel "WIR BAUEN GEMEINDE - FÜR LEUTE UND MORGEN - klimafreundlich"'
			/>
			<div class="text-center">
				Ein <span class="font-bold uppercase text-blue">Klimaplus</span>-<span
					class="font-bold uppercase text-green"
					>Gemeindezentrum</span
				>
				der Evangelischen Kirchengemeinde in Bühl
			</div>
			<h2 class="mt-8 text-2xl font-bold">
				Es geht um Klimaschutz und Bewahrung der Schöpfung, aber auch um die Zukunftsfähigkeit
				unseres Gebäudes
			</h2>
			<p>
				Dazu wollen wir als Kirchengemeinde unseren Beitrag leisten. Schon früh wurde das Augenmerk
				auf die regenerative Energie gerichtet. Alle Neubauten bekommen eine Fußbodenheizung, die
				Kirche ein neues Heizregister und die energetisch sanierten Bestandsgebäude die Möglichkeit,
				die Vorlaufstemperatur zu reduzieren. Damit ist es möglich, über
				<span class="font-bold text-blue">80 kW Wärmepumpen</span>,
				<span class="font-bold text-blue">3.000 l Wasserspeicher</span> und eine
				<span class="font-bold text-blue">99 kWp PV-Anlage</span> die Versorgung nachhaltig zu
				betreiben. Im Jahresmittel produzieren wir mehr regenerative Energie als verbraucht wird und
				leisten damit einen klimapositiven Beitrag, der uns auf Dauer auch finanziell zugutekommen
				wird.
			</p>
			<div id="psalm-container" class="relative mt-8 w-full">
				<BaseImageCredit
					class="w-full"
					src="/external/solar_peterschreiber.media - istockphoto.com.webp"
					credit="peterschreiber.media - iStockphoto.com"
					alt="Aufnahme einer Fläche von Photovoltaik-Elemente. Die tiefstehende Sonne scheint von hinten über die Oberkante der Panele"
				/>
				<div id="psalm">
					<div class="font-bold text-white">
						<div class="text-left">„Vom Aufgang der Sonne</div>
						<div class="text-right">bis zu ihrem Niedergang</div>
						<div class="text-center">sei gelobet der Name des Herrn“</div>
					</div>
					<div class="text-right text-base text-orange">Psalm 113,3</div>
				</div>
			</div>
			<p>
				<BaseImageTitle
					id="real-green-award"
					class="max-w-full sm:float-right sm:ml-6 sm:mt-1 sm:w-60"
					src="/external/real-green-award.webp"
					alt='Preis mit der Aufschrift "REAL GREEN AWARD - Gewinner Jury Award 2022 Evangelische Kirchengemeinde Bühl"'
				>
					Eine mit dem <span class="font-bold text-blue">REAL GREEN AWARD 2022</span> ausgezeichnete
					Idee steht kurz vor der Realisierung
				</BaseImageTitle>
				In dieses biblische Lied stimmen wir als Kirchengemeinde gerne ein und fügen in Zukunft
				hinzu
				<span class="font-bold text-blue">„… und seien unsere Energiespeicher gut gefüllt“</span>.
				Denn in der Tat ist die Sonne ein Gottesgeschenk. Auch als wunderbarer „Lieferant“ für alle
				Energie und Wärme, die wir in unserem erneuerten Gemeindezentrum brauchen und die wir in
				Zukunft nicht mehr auf Kosten einer negativen CO<sub>2</sub>-Bilanz erzeugen wollen.
			</p>
			<p>
				Wer trägt bei zur energetischen Sanierung unseres Gemeindezentrums? Wer hilft uns dabei,
				sauber und sicher Energie über die beiden Wärmepumpen zu produieren? Wer trägt dazu bei,
				unser Gemeindezentrum gut für die Zukunft aufzustellen und – auch im übertragenen Sinn –
				mehr Licht und Wärme in unsere Stadt zu bringen?
			</p>
			<h2 class="text-2xl font-bold text-green">
				Die 7 wichtigsten Maßnahmen auf dem Web zum
				<span class="uppercase text-blue">KlimaPlus-Gemeindezentrum</span>
			</h2>
			<ol class="list-decimal pl-10">
				<li>
					Das konsequente Dämmen der Gebäudehülle mit neuen Fenster- und Glasfassaden und einer
					neuen Dämmung des Kirchen- und Saaldachs.
				</li>
				<li>
					Das Abschalten des Gasbrenners und der Installation von zwei leistungsstarken Wärmepumpen.
				</li>
				<li>Eine Fußbodenheizung in allen Neubauten.</li>
				<li>Ein neues Niedrigtemperatur-Heizregister für die Kirchenheizung.</li>
				<li>Die solare Aktivierung der Dachflächen.</li>
				<li>Rigolen und eine Zisterne, um Regenwasser auf dem Grundstück versickern zu lassen.</li>
				<li>Ein Gründach für einen positiven Beitrag für ein angenehmes Stadtklima.</li>
			</ol>
			<p>
				Die Investition für die Zukunft und für den Klimaschutz gibt es nicht umsonst. Insgesamt
				belaufen sich die 7 Maßnahmen auf eine Summe von über 800.000 €. Eine wichtige und richtige
				Investition zum Erhalt der Schöpfung.
			</p>
			<h3 class="text-2xl font-bold text-green">
				<span class="uppercase">
					Werden auch Sie <span class="text-blue">Klimaplus</span>-Pate und -Patin
				</span>
			</h3>
			<p>
				Wählen Sie in der Grafik ein Element aus, für welches Sie die Patenschaft übernehmen
				möchten. Nachdem Sie ihre E-Mail-Adresse eingegeben haben, senden wir ihnen den
				Patenschaftsablauf zu. Wenn Sie namentlich in der Grafik genannt werden möchten, können Sie
				auch optional ihren Namen eintragen.
			</p>
			<p>Für folgende Komponenten unserer PV-Anlage können Sie Pate werden:</p>
			<div class="max-w-full overflow-x-auto">
				<table id="element-list" class="w-min-80 mx-auto">
					<tbody>
						<tr>
							<td class="bg-orange text-2xl font-bold text-white">1</td>
							<td>122 St. PV-Indach-Module auf dem Kirchendach</td>
							<td class="font-bold text-blue" style="padding-inline: 1em">je</td>
							<td class="text-nowrap text-right font-bold text-blue">1.250 €</td>
						</tr>
						<tr>
							<td class="bg-orange text-2xl font-bold text-white">2</td>
							<td>53 St. PV-Aufdach-Module auf dem Ostdach vom Kindergarten</td>
							<td class="font-bold text-blue" style="padding-inline: 1em">je</td>
							<td class="text-nowrap text-right font-bold text-blue">550 €</td>
						</tr>
						<tr>
							<td class="bg-orange text-2xl font-bold text-white">3</td>
							<td>42 St. PV-Aufdach-Module auf dem Westdach vom Gemeindehaus</td>
							<td class="font-bold text-blue" style="padding-inline: 1em">je</td>
							<td class="text-nowrap text-right font-bold text-blue">550 €</td>
						</tr>
						<tr>
							<td class="bg-orange text-2xl font-bold text-white">4</td>
							<td>28 St. PV-Aufdach-Module auf dem Süddach vom Pfarrhaus</td>
							<td class="font-bold text-blue" style="padding-inline: 1em">je</td>
							<td class="text-nowrap text-right font-bold text-blue">550 €</td>
						</tr>
						<tr>
							<td class="bg-orange text-2xl font-bold text-white">5</td>
							<td>2 St. Batteriespeicher mit je 11 kWh</td>
							<td class="font-bold text-blue" style="padding-inline: 1em">je</td>
							<td class="text-nowrap text-right font-bold text-blue">12.000 €</td>
						</tr>
					</tbody>
				</table>
			</div>
		</div>
		<BasePV v-model:selected_element="selected_element">
			<template #header v-if="selected_element !== undefined">
				{{ get_element_roof(selected_element?.mid) }}
			</template>
			<template v-if="selected_element !== undefined">
				<form
					v-if="is_element_available(selected_element.mid)"
					ref="reserve_form"
					id="reserve-input-box"
					class="flex flex-col gap-1"
					@submit="submit"
					@keydown.enter="enter_press = true"
				>
					<p>
						Um Pate für {{ get_element_type(selected_element.mid, true) }} zu werden, tragen Sie
						ihre E-Mail-Adresse ein. Wir senden ihnen anschließend den genauen Ablauf zu.
					</p>
					<input
						type="email"
						name="mail"
						id="input-email"
						class="flex-1 rounded px-2 outline outline-2 invalid:text-red-500"
						v-model="selected_element.email"
						placeholder="E-Mail"
						required
					/>
					<p>
						Wenn Sie ihren Namen auf der Website (öffentlich) genannt haben möchten, tragen Sie
						diesen ein.
					</p>
					<input
						type="text"
						name="name"
						id="input-name"
						class="flex-1 rounded px-2 outline outline-2"
						maxlength="70"
						v-model="selected_element.name"
						placeholder="Name (optional)"
						autocomplete="off"
					/>
					<input type="submit" style="display: none" id="submit-reservation" />
					<label for="submit-reservation" class="text-center">
						<BaseButton
							class="bg-green font-bold text-white"
							:disabled="!reserve_form?.checkValidity()"
						>
							<FontAwesomeIcon :icon="faCheck" /> Reservieren
						</BaseButton>
					</label>
				</form>
				<div v-else-if="selected_element.reserved">Dieses Modul ist derzeit reserviert.</div>
				<div v-else id="tooltip-sold">
					<template v-if="!!selected_element.name">
						Pate für dieses Element ist {{ selected_element.name }}
					</template>
					<template v-else> Dieses Element hat bereits einen Paten </template>
				</div>
			</template>
		</BasePV>
	</AppLayout>
</template>

<style scoped>
	p {
		margin-block: 1em;
	}

	#text-content {
		max-width: min(40em, 100%);

		padding-inline: 0.25em;
	}

	#text-content > p {
		text-align: justify;
	}

	sub {
		font-size: 0.8em;
	}

	ul {
		list-style: none;
	}

	#psalm-container {
		container-name: psalm-container;
		container-type: inline-size;
	}

	@container psalm-container (width > 1px) {
		#psalm {
			font-size: 3.75cqw;

			position: absolute;
			right: 1em;
			bottom: 1em;

			display: flex;
			flex-direction: column;
			justify-content: end;
			align-items: flex-end;
		}
	}

	#element-list td:first-of-type {
		display: flex;
		align-items: center;
		justify-content: center;

		height: 1.5em;

		aspect-ratio: 1;
		border-radius: 100%;
	}

	#element-list td {
		margin: 0.25em;
	}

	#reserve-input-box p {
		margin-block: unset;
	}
</style>
