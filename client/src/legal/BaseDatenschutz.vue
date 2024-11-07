<script setup lang="ts">
import BaseStatic from '@/BaseStatic.vue';

</script>

<template>
	<BaseStatic>
		<h1>Datenschutz</h1>
		<h2>I. Informationen über die Verarbeitung Ihrer Daten gemäß Art. 13 der Datenschutz-Grundverordnung (DS-GVO)</h2>
		<h3>1. Verantwortlicher und Datenschutzbeauftragter</h3>

		<p>
			Verantwortlich für diese Website ist
		</p>
		<p>
			<span class="bold"> Dr. Götz Häuser</span><br>
			Evangelische Kirchengemeinde Bühl<br>
			Johannesstr. 5<br>
			77815 Bühl<br>
			<a href="mailto:pfarrer@evkirchebuehl.de">pfarrer@evkirchebuehl.de</a>
		</p>
		<p>
			Den Datenschutzbeauftragten erreichen Sie per E-Mail unter
		</p>
		<p>
			<a href="mailto:datenschutz@evkirchebuehl.de">datenschutz@evkirchebuehl.de</a>
		</p>

		<h3>2. Daten, die für die Bereitstellung der Website und die Erstellung der Protokolldateien verarbeitet werden</h3>
		<h4>a. Welche Daten werden für welchen Zweck verarbeitet?</h4>
		<p>
			Bei der Reservierung eines Elements für eine Patenschaft werden folgende personenbezogene Daten gespeichert:
		</p>
		<ul>
			<li>Identifikationsnummer des Elements</li>
			<li>Name, auf den die Patenschaft gelistet wird (sofern angegeben)</li>
			<li>Datum und Uhrzeit der Reservierung (wird nach erfolgreicher Patenschaft oder Abbruch gelöscht)</li>
		</ul>

		<p>
			Wenn Sie der Ansicht sind, dass die Verarbeitung Ihrer personenbezogenen Daten gegen Datenschutzrecht verstößt, haben Sie nach Ar. 77 Abs. 1 DS-GVO das Recht, sich bei einer Datenschutzaufsichtsbehörde eigener Wahl zu beschweren.
			Hierzu gehört auch die für den Verantwortlichen zuständige Datenschutzaufsichtsbehörde: Landesbeauftragter für Datenschutz und Informationsfreiheit Baden-Württemberg, <a href="https://www.baden-wuerttemberg.datenschutz.de/beschwerde">https://www.baden-wuerttemberg.datenschutz.de/beschwerde.</a>
		</p>
		<h4>g. Recht auf Datenübertragbarkeit</h4>
		<p>
			Für den Fall, dass die Voraussetzungen des Art. 20 Abs. 1 DS-GVO vorliegen, steht Ihnen das Recht zu, sich Daten, die wir auf Grundlage Ihrer Einwilligung oder in Erfüllung eines Vertrags automatisiert verarbeiten, an sich oder an Dritte aushändigen zu lassen.
			Die Erfassung der Daten zur Bereitstellung der Website sind für den Betrieb der Internetseite zwingend erforderlich.
			Sie beruhen daher nicht auf einer Einwilligung nach Art. 6 Abs. 1 Buchstabe a DS-GVO oder auf einem Vertrag nach Art. 6 Abs. 1 Buchstabe b DS-GVO, sondern sind nach Art. 6 Abs. 1 Buchstabe f DS-GVO gerechtfertigt. Die Voraussetzungen des Art. 20 Abs. 1 DS-GVO sind demnach insoweit nicht erfüllt.
		</p>
		<h2>II. Recht auf Widerspruch gemäß Art. 21 Abs. 1 DS-GVO</h2>
		<p>
			Sie haben das Recht, aus Gründen, die sich aus Ihrer besonderen Situation ergeben, jederzeit gegen die Verarbeitung Ihrer personenbezogenen Daten, die aufgrund von Artikel 6 Abs. 1 Buchstabe f DS-GVO erfolgt, Widerspruch einzulegen.
			Der Verantwortliche verarbeitet die personenbezogenen Daten dann nicht mehr, es sei denn, er kann zwingende schutzwürdige Gründe für die Verarbeitung nachweisen, die die Interessen, Rechte und Freiheiten der betroffenen Person überwiegen, oder die Verarbeitung dient der Geltendmachung, Ausübung oder Verteidigung von Rechtsansprüchen.
			Die Erfassung der Daten zur Bereitstellung der Website sind für den Betrieb der Internetseite zwingend erforderlich.
		</p>
		<h2>Cookies</h2>
		<p>
			Diese Seite nutzt Cookies ausschließlich für den Login im Admin-Panel.
			Dieses Session-Cookie beinhaltet Session-Token, User-ID, Token-ID.
			Als technisch notwendiges Cookie benötigt dies keiner Zustimmung.
			Bei regulären Besuchern werden keine Cookies gespeichert.
		</p>
	</BaseStatic>	
</template>

<style scoped>
	p {
		margin-block: 0.55em;
	}

	p.container {
		margin-block: 1em;
	}

	a {
		color: var(--v-color-blue);
	}

	h2, h3 {
		margin-top: 1em;
		margin-bottom: 0.25em;
	}

	ul {
		margin-block: 1em;
	}
</style>
