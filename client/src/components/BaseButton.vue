<script setup lang="ts">
	defineProps<{
		disabled?: boolean;
		square?: boolean;
	}>();
</script>

<template>
	<div
		class="base-button inline-block cursor-pointer rounded-sm p-1 shadow-[0.125rem_0.125rem_0.125rem_black] outline outline-1"
		:class="{
			'cursor-not-allowed brightness-150 contrast-[0.25]': disabled,
			'flex aspect-square h-6 items-center justify-center': square,
			'hover:contrast-75 active:shadow-[0.125rem_0.125rem_0.125rem_inset_black]': !disabled
		}"
		@click="(e) => (disabled ? e.stopImmediatePropagation() : undefined)"
	>
		<slot></slot>
	</div>
</template>

<style scoped></style>
