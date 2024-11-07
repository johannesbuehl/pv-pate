<script setup lang="ts">
	defineProps<{
		src: string;
	}>();
</script>

<template>
	<div>
		<img :src="src" />
		<div id="title"><slot></slot></div>
	</div>
</template>

<style scoped>
	img {
		width: 100%;
		display: block;
	}

	#title {
		text-align: center
	}
</style>
