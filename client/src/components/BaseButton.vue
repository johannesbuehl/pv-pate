<script setup lang="ts">

	defineProps<{
		disabled?: boolean;
		square?: boolean;
	}>();
</script>

<template>
	<div class="base-button" :class="{ disabled, square }" @click="(e) => disabled ? e.stopImmediatePropagation() : undefined">
		<slot></slot>
	</div>
</template>

<style scoped>
	.base-button {
		cursor: pointer;

		border-radius: 0.125em;		
		padding: 0.25em;
		display: inline-block;

		outline: 0.0625em solid black;

		box-shadow: 0.0625em 0.0625em 0.0625em black;
	}

	.base-button:active {
		box-shadow: 0.0625em 0.0625em 0.0625em black inset;
	}

	.base-button:hover {
		filter: brightness(1.1);
	}

	.base-button.disabled {
		filter: contrast(0.25) brightness(1.5);

		cursor: not-allowed;
	}

	.base-button.square {
		display: flex;
		
		aspect-ratio: 1;
		height: 1.5em;
		
		align-items: center;
		justify-content: center;
	}
</style>
