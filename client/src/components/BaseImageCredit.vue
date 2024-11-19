<script setup lang="ts">
	defineProps<{
		src: string;
		credit: string;
		alt?: string;
	}>();
</script>

<template>
	<div class="image-credit">
		<img :src="src" :alt="alt" />
		<div id="credit">{{ credit }}</div>
	</div>
</template>

<style scoped>
	.image-credit {
		position: relative;

		container-name: image-container;
		container-type: inline-size;
	}

	img {
		width: 100%;
		display: block;
	}

	@container image-container (width > 1px) {
		#credit {
			position: absolute;

			bottom: 0;
			left: 0;

			color: white;

			padding: 0.25em;

			backdrop-filter: brightness(0.5);

			font-size: 1.75cqw;
		}
	}
</style>
