<script setup lang="ts">
	import { ref } from 'vue';

	import BasePV, { get_element_roof, get_element_type, type Element } from './components/BasePV.vue';
	import AppLayout from './components/AppLayout/AppLayout.vue';
	import BaseButton from './components/BaseButton.vue';
	import { faCheck } from '@fortawesome/free-solid-svg-icons';
	import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
	import { api_call, is_element_available, type APICallResult } from './lib';
	import { elements_db, type ElementsDB } from './Globals';
import BaseImageCredit from './components/BaseImageCredit.vue';
import BaseImageTitle from './components/BaseImageTitle.vue';

	const reserve_form = ref<HTMLFormElement>();

	const selected_element = ref<Element & { email: string }>();

	let enter_press: boolean = false;
	async function submit(e: Event) {
		e.preventDefault();

		// if the form was submitted through the enter-key, let the user confirm
		if (enter_press) {
			enter_press = false;

			let confirm_string = `Reservierung bestätigen?\nE-Mail: ${selected_element.value?.email}\n`;

			if (selected_element.value?.name != undefined && !selected_element.value.name && selected_element.value?.name?.length > 1) {
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
				
			response = await api_call<ElementsDB>(method, "elements", { mid: selected_element.value.mid }, {
				name: selected_element.value.name,
				mail: selected_element.value.email
			});
			
			if (response.ok) {
				elements_db.value = (await response.json());
				
				selected_element.value = undefined;
			} else {
				alert(`Error during database write: ${await response.text()}`);
			}
		}
	}
</script>

<template>
	<AppLayout>
		<div id="text-content">
			<div id="header" class="center">
				<div class="blue bold upper stretch" style="font-size: 1.175em;">
					Wir bauen Gemeinde
				</div>
				<div class="green upper" style="font-size: 1.1em;">
					Für Leute und morgen
				</div>
				<div class="orange">
					klimafreundlich
				</div>
			</div>
			<p class="center">
				Ein <span class="blue upper bold">Klimaplus</span>-<span class="green upper bold">Gemeindezentrum</span> der Evangelischen Kirchengemeinde in Bühl
			</p>
			<h2 class="bold">Klimaschutz ist Bewahrung der Schöpfung</h2>
			<p>
				Und davon hat auch die Kirche ihren Beitrag zu leisten. Schon in der frühen Planungsphase wurde das Augenmerk auf die regenerative Energie gerichtet. Alle Neubauten bekommen eine Fußbodenheizung, die Kirche ein neues Heizregister und die energetisch sanierten Bestandsgebäude die Möglichkeit die Vorlaufstemperatur zu reduzieren. Damit ist es möglich über <span class="blue bold">80 kW Wärmepumpen</span>, <span class="blue bold">3.000 l Wasserspeicher</span> und eine <span class="blue bold">99 kWp PV-Anlage</span> die Versorgung nachhaltig zu betreiben. Im Jahresmittel produzieren wir mehr regenerative Energie als verbraucht wird und leisten damit einen klimapositiven Beitrag.
			</p>
			<div id="psalm-container">
				<BaseImageCredit style="width: 100%;" src="/solar_peterschreiber.media - istockphoto.com.webp" credit="peterschreiber.media - iStockphoto.com" />
				<div id="psalm">
					<div class="bold white">
						<div class="left">„Vom Aufgang der Sonne</div>
						<div class="right">bis zu ihrem Niedergang</div>
						<div class="center">sei gelobet der Name des Herrn“</div>
					</div>
					<div id="citation" class="orange right">
						Psalm 113,3
					</div>
				</div>
			</div>
			<p>
				So lautet ein bekannter Kanon (Psalm 113,3).
			</p>
			<p>
				<BaseImageTitle id="real-green-award" class="float-right" src="/real-green-award.webp">
					Eine mit dem <span class="blue bold">REAL GREEN AWARD 2022</span> ausgezeichnete Idee steht kurz vor der Realisierung
				</BaseImageTitle>
				Wir stimmen dankbar ein und fügen in Zukunft gerne hinzu <span class="bold blue">„…und seien unsere Energiespeicher gut gefüllt“</span>. Denn in der Tat ist die Sonne ein Gottesgeschenk. Auch als wunderbarer „Lieferant“ für alle Energie und Wärme, die wir in unserem erneuerten Gemeindezentrum brauchen und die wir in Zukunft nicht mehr auf Kosten einer negativen CO<sub>2</sub>-Bilanz erzeugen wollen.
			</p>
			<p>
				Wer finanziert und spendet uns ein Solarpanel (oder mehrere)? Wer hilft uns dabei, sauber und sicher Energie zu gewinnen – für unsere Kirchengemeinde und für andere, indem wir unsere Überschüsse verkaufen, um Erträge für die Finanzierung des Gemeindezentrums zu generieren bzw. für den Stromkauf im Winter. Wer trägt dazu bei, unser Gemeindezentrum gut für die Zukunft aufzustellen und – auch im übertragenen Sinn – mehr Licht und Wärme in unsere Stadt zu bringen?
			</p>
			<h3 class="green bold">
				<span class="upper">
					Und so werden Sie <span class="blue">Klimaplus</span>-Pate und -Patin
				</span>
				für ein CO<sub>2</sub>-freies Gemeindezentrum
			</h3>
			<div class="side-split">
				<img id="banner-energiespenderin" src="/public/banner_energiespenderin.svg"/>
				<p>
					Wählen Sie in der Grafik ein Element aus, für welches Sie die Patenschaft übernehmen möchten.
					Nachdem Sie ihre E-Mail-Adresse eingegeben haben, senden wir ihnen den Spendenablauf zu.
					Wenn Sie namentlich in der Grafik genannt werden möchten, können Sie auch optional ihren Namen eintragen.
				</p>
			</div>
			<p>
				Für folgende Komponenten unserer PV-Anlage können Sie Pate werden:
			</p>
			<table id="element-list">
				<tbody>
					<tr>
						<td>1</td>
						<td>122 St. PV-Indach-Module auf dem Kirchendach</td>
						<td class="bold blue right nowrap">1.250 €</td>
					</tr>
					<tr>
						<td>2</td>
						<td>53 St. PV-Aufdach-Module auf dem Ostdach vom Kindergarten</td>
						<td class="bold blue right nowrap">550 €</td>
					</tr>
					<tr>
						<td>3</td>
						<td>42 St. PV-Aufdach-Module auf dem Westdach vom Gemeindehaus</td>
						<td class="bold blue right nowrap">550 €</td>
					</tr>
					<tr>
						<td>4</td>
						<td>28 St. PV-Aufdach-Module auf dem Süddach vom Pfarrhaus</td>
						<td class="bold blue right nowrap">550 €</td>
					</tr>
					<tr>
						<td>5</td>
						<td>2 St. Batteriespeicher mit je 11 kWh</td>
						<td class="bold blue right nowrap">12.000 €</td>
					</tr>
				</tbody>
			</table>
		</div>
		<BasePV
			v-model:selected_element="selected_element"
		>
			<template #header
				v-if="selected_element !== undefined"
			>
				{{ get_element_roof(selected_element?.mid) }}
			</template>
			<template
				v-if="selected_element !== undefined"
			>
				<form
					v-if="is_element_available(selected_element.mid)"
					ref="reserve_form"
					id="reserve-input-box"
					@submit="submit"
					@keydown.enter="enter_press = true"
				>
					<p>
						Um Pate für {{ get_element_type(selected_element.mid, true) }} zu werden, tragen Sie ihre E-Mail-Adresse ein. Wir senden ihnen anschließend den genauen Ablauf zu.
					</p>
					<input type="email" name="mail" id="input-email" v-model="selected_element.email" placeholder="E-Mail" required />
					<p>
						Wenn Sie ihren Namen auf der Website (öffentlich) genannt haben möchten, tragen Sie diesen ein.
					</p>
					<input type="text" name="name" id="input-name" maxlength="70" v-model="selected_element.name" placeholder="Name (optional)" autocomplete="off" />
					<input type="submit" style="display: none;" id="submit-reservation" />
					<label for="submit-reservation" class="center">
						<BaseButton id="button-reserve" :disabled="!reserve_form?.checkValidity()">
							<FontAwesomeIcon :icon="faCheck" /> Reservieren
						</BaseButton>
					</label>
				</form>
				<div
					v-else-if="selected_element.reserved"
				>
					Dieses Modul ist derzeit reserviert.
				</div>
				<div
					v-else
					id="tooltip-sold"
				>
					<template v-if="!!selected_element.name">
						Pate für dieses Element ist {{ selected_element.name }}
					</template>
					<template v-else>
						Dieses Element hat bereits einen Paten
					</template>
				</div>
			</template>
		</BasePV>
	</AppLayout>
</template>

<style scoped>
	h2, h3 {
		margin-top: 3em;
		margin-bottom: 0em;
	}

	p {
		margin-block: 1em;
	}

	#header {
		font-family: var(--font-header);

		font-size: 2.5em;

		margin-top: 2em;
	}

	#header > .stretch {
		letter-spacing: 0.02em;
	}

	#text-content {
		max-width: 40em;

		padding-inline: 0.25em;

		text-align: justify;
	}

	sub {
		font-size: 0.8em;
	}

	ul {
		list-style: none;
	}
	
	#real-green-award {
		max-width: 100%;
	}

	.side-split {
		display: flex;
		align-items: center;

		gap: 1em;
	}

	#banner-energiespenderin {
		display: none;
	}

	@media screen and (min-width: 875px) {
		#banner-energiespenderin {
			display: unset;

			height: 15em;
			height: 4em;

			transform: rotate(-3deg);
		}
	}

	@media screen and (min-width: 600px) {
		#real-green-award {
			width: 15em;
		}
	}

	#psalm-container {
		margin-top: 2em;

		position: relative;

		width: 100%;

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

	#citation {
		font-size: 0.8em;
	}

	#element-list {
		max-width: 30em;

		margin-inline: auto;
	}

	#element-list td:first-of-type {
		background-color: orange;
		color: var(--v-color-white);
		font-weight: bold;
		font-size: 1.5em;

		display: flex;
		align-items: center;
		justify-content: center;

		height: 1.5em;

		aspect-ratio: 1;
		border-radius: 100%;
	}

	#element-list td {
		margin: 0.25em
	}

	#reserve-input-box > input {
		flex: 1;
	}
	
	#input-email:invalid {
		color: red;
	}

	#reserve-input-box {
		display: flex;
		flex-direction: column;

		gap: 0.25em;
	}
	
	#button-reserve {
		background-color: var(--v-color-green);
		color: white;

		font-weight: bold
	}
</style>
