<script setup lang="ts">
import BaseStatic from '@/BaseStatic.vue';

</script>

<template>
	<BaseStatic>
		<h1>Datenschutz</h1>
		<h2>Datenschutzhinweise - Evangelische Kirchengemeinde Bühl</h2>
		<p>
			Wir nehmen den Schutz Ihrer personenbezogenen Daten ernst und beachten die Regeln der anwendbaren Datenschutzgesetze. Es ist uns wichtig, dass Sie jederzeit wissen, welche personenbezogenen Daten während Ihres Besuchs auf unserer Webseite und bei Nutzung unserer Angebote erhoben werden und wie wir sie verwenden.<br>
			Personenbezogene Daten sind Einzelangaben über persönliche und sachliche Verhältnisse einer bestimmten oder bestimmbaren natürlichen Person.<br>
			Mit den nachfolgenden Hinweisen möchten wir Sie über die Erhebung und Verarbeitung Ihrer Daten auf unseren Webseiten informieren:
		</p>

		<h2>Cookies</h2>
		<p>
			Diese Seite nutzt Cookies ausschließlich für den Login im Admin-Panel. Dieses Session-Cookie beinhaltet Session-Token, User-ID, Token-ID. Als technisch notwendiges Cookie benötigt dies keiner Zustimmung.
		</p>
		<p>
			Bei regulären Besuchern werden keine Cookies gespeichert.
		</p>

		<h2>Datenschutzrechtlich verantwortliche Stelle</h2>
		Die datenschutzrechtlich verantwortliche Stelle ist im Impressum unserer Webseite näher bezeichnet. Sprechen Sie uns bei Fragen gerne über die dort hinterlegten Kontaktdaten an.
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

	h3 {
		color: unset;
	}

	ul {
		margin-block: 1em;
	}
</style>
