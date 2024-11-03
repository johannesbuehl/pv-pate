<script setup lang="ts">
	import { ref } from 'vue';

	import BasePV, { get_element_roof, get_element_type, type Element } from './components/BasePV.vue';
	import { reserved_elements } from './Globals';
import AppLayout from './components/AppLayout/AppLayout.vue';

	const selected_element = ref<Element>();
</script>

<template>
	<AppLayout>
		<div id="header" class="center">
			<div class="blue bold upper stretch">
				Wir bauen Gemeinde
			</div>
			<div class="green upper">
				Für Leute und morgen
			</div>
			<div class="orange">
				klimafreundlich
			</div>
		</div>
		<div id="text-content">
			<p>
				Eine mit dem REAL GREEN AWARD 2022 ausgezeichnete Idee steht kurz vor der Realisierung.
			</p>
			<p>
				Ein <span class="blue upper bold">Klimaplus</span>-<span class="green upper bold">Gemeindezentrum</span> der Evangelischen Kirchengemeinde in Bühl
			</p>
			<p>
				<span class="blue bold">Klimaschutz ist Bewahrung der Schöpfung</span><br>
				Und davon hat auch die Kirche ihren Beitrag zu leisten. Schon in der frühen Planungsphase wurde das Augenmerk auf die regenerative Energie gerichtet. Alle Neubauten bekommen eine Fußbodenheizung, die Kirche ein neues Heizregister und die energetisch sanierten Bestandsgebäude die Möglichkeit die Vorlaufstemperatur zu reduzieren. Damit ist es möglich über 80 kW Wärmepumpen, 3.000 l Wasserspeicher und eine 99 kWp PV Anlage die Versorgung nachhaltig zu betreiben. Im Jahresmittel produzieren wir mehr regenerative Energie als verbraucht wird und leisten damit einen klimapositiven Beitrag.
			</p>
			<div id="psalm-verse">
				<div id="psalm" class="blue bold">
					<div class="left">
						„Vom Aufgang der Sonne
					</div>
					<div class="right">
						bis zu ihrem Niedergang
					</div>
					<div class="center">
						sei gelobet der Name des Herrn“
					</div>
				</div>
				<div id="citation" class="black right">
					Psalm 113,3
				</div>
			</div>
			<p>
				So lautet ein bekannter Kanon (Psalm 113,3).
			</p>
			<p>
				Wir stimmen dankbar ein und fügen in Zukunft gerne hinzu <span class="bold blue">„…und seien unsere Energiespeicher gut gefüllt“</span>. Denn in der Tat ist die Sonne ein Gottesgeschenk. Auch als wunderbarer „Lieferant“ für alle Energie und Wärme, die wir in unserem erneuerten Gemeindezentrum brauchen und die wir in Zukunft nicht mehr auf Kosten einer negativen CO2-Bilanz erzeugen wollen.
			</p>
			<p>
				Wer finanziert und spendet uns ein Solarpanel (oder mehrere)? Wer hilft uns dabei, sauber und sicher Energie zu gewinnen – für unsere Kirchengemeinde und für andere, indem wir unsere Überschüsse verkaufen, um Erträge für die Finanzierung des Gemeindezentrums zu generieren bzw. für den Stromkauf im Winter. Wer trägt dazu bei, unser Gemeindezentrum gut für die Zukunft aufzustellen und – auch im übertragenen Sinn – mehr Licht und Wärme in unsere Stadt zu bringen?
			</p>
			<p class="green bold">
				<span class="upper">
					Und so werden Sie <span class="blue">Klimaplus</span>-Pate und -Patin
				</span>
				für ein CO<sub>2</sub>-freies Gemeindezentrum
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
				<div
					v-if="reserved_elements[selected_element.mid] !== undefined"
					id="tooltip-sold"
				>
					Gespendet von {{ selected_element.name }}
				</div>
				<div
					v-else
					id="tooltip-buy"
				>
					Dieses {{ get_element_type(selected_element.mid) }} spenden<br>
					<a href="https://www.evkirchebuehl.de" target="_blank" rel="noopener noreferrer">Dummy-Link</a>
				</div>
			</template>
		</BasePV>
	</AppLayout>
</template>

<style scoped>
	#header {
		font-size: 2em;
	}

	#header > .stretch {
		letter-spacing: 0.05em;
	}

	#header > .condense {
		letter-spacing: -0.05em;
	}

	#text-content {
		max-width: 40em;

		display: flex;
		flex-direction: column;
		gap: 0.75em;
	}

	ul {
		list-style: none;
	}

	#psalm-verse {
		display: flex;
		flex-direction: column;
		gap: 0.67em;
	}
	
	#psalm-verse {
		margin-inline: auto;
		width: 21em;
		max-width: 90cqw;
	}

	#citation {
		margin-right: 3em;
		font-size: 0.8em;
	}

	#element-list {
		margin-block-start: 1em;
	}

	#element-list td:first-of-type {
		background-color: orange;
		color: var(--v-color-white);
		font-weight: bold;
		padding: 0.45em;

		display: flex;
		align-items: center;
		justify-content: center;

		height: 2em;

		aspect-ratio: 1;
		border-radius: 100%;
	}

	#element-list td {
		margin: 0.5em
	}

	a {
		text-decoration: underline;
		font-style: italic;
	}
</style>
