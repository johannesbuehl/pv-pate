<script lang="ts">
	export interface Element {
		mid: string;
		name?: string;
		reserved?: boolean;
	}

	const element_type_map: Record<string, string> = {
		pv: "PV-Modul",
		bs: "Batteriespeicher"
	};
	const dieses_element_type_map: Record<string, string> = {
		pv: "dieses PV-Modul",
		bs: "diesen Batteriespeicher"
	};

	export function get_element_type(mid: string, dieses: boolean = false): string {
		if (dieses) {
			return dieses_element_type_map[mid.slice(0, 2)];
		} else {
			return element_type_map[mid.slice(0, 2)];
		}
	}

	const row_roof_map = {
		d: "Kirchendach",
		l: "Ostdach",
		r: "Westdach",
		v: "Süddach"
	};

	export function get_element_string(mid: string): string {
		// handle batteries extra
		if (mid.slice(0, 2) === "bs") {
			return `${get_element_type(mid)} ${mid.slice(3).toUpperCase()}`;
		} else {
			const mid_row = mid.slice(3, 4);

			for (let row of Object.keys(row_roof_map)) {
				if (mid_row <= row) {
					return `${get_element_type(mid)} ${mid.slice(3).toUpperCase()}`;
				}
			}

			return mid_row;
		}
	}

	export function get_element_roof(mid: string): string {
		let return_string = get_element_string(mid);

		// add the roof for pv-modules
		if (mid.slice(0, 2) === "pv") {
			const mid_row = mid.slice(3, 4);

			for (let [row, roof] of Object.entries(row_roof_map)) {
				if (mid_row <= row) {
					return_string += ` (${roof})`;

					break;
				}
			}
		}
		return return_string;
	}
</script>

<script setup lang="ts">
	import { onBeforeMount, onMounted, onUnmounted, ref } from "vue";

	import BaseTooltip from "./BaseTooltip.vue";
	import { is_element_available, get_element } from "@/lib";

	const svg = ref<string>();

	const svg_path = "external/elements.svg";

	const svg_wrapper = ref<HTMLDivElement>();
	const tooltip = ref<HTMLDivElement>();
	const svg_selected_element = ref<SVGRectElement>();

	const selected_element = defineModel<Element | undefined>("selected_element");

	onBeforeMount(async () => {
		const svg_request = fetch(svg_path);

		if ((await svg_request).ok) {
			svg.value = await (await svg_request).text();
		}
	});

	function hide_tooltip() {
		svg_selected_element.value = undefined;
		selected_element.value = undefined;
	}

	function prepare_svg(r: string): string {
		const parser = new DOMParser();
		const svg_dom = parser.parseFromString(r, "image/svg+xml").documentElement;

		svg_dom.removeAttribute("width");
		svg_dom.removeAttribute("height");
		// svg_dom.removeAttribute("viewBox")
		svg_dom.id = "main-content";

		const prepare_element = (ele: SVGPathElement, classname: string) => {
			ele.classList.add("element");

			ele.classList.add(classname);

			if (!is_element_available(ele.id)) {
				ele.classList.add("sold");
			}
		};

		// select all elements
		const elements: SVGPathElement[] = Array.from(
			svg_dom.querySelectorAll<SVGPathElement>("[id^='pv-']")
		);
		elements.forEach((element) => prepare_element(element, "module"));

		// select all batteries
		const batteries: SVGPathElement[] = Array.from(
			svg_dom.querySelectorAll<SVGPathElement>("[id^='bs-']")
		);
		batteries.forEach((element) => prepare_element(element, "battery"));

		let svg_string = new XMLSerializer().serializeToString(svg_dom);

		return svg_string;
	}

	onMounted(() => {
		document.addEventListener("click", on_click);
	});

	onUnmounted(() => {
		document.removeEventListener("click", on_click);
	});

	function on_click(e: MouseEvent) {
		const target = (e.target as SVGElement).closest(".element");

		if (target) {
			// clear the previous selected element
			svg_selected_element.value = target as SVGRectElement;
			const mid = svg_selected_element.value?.id;

			selected_element.value = get_element(mid);
		}
	}
</script>

<template>
	<div id="wrapper">
		<div v-if="!!svg" id="div-svg" ref="svg_wrapper" v-html="prepare_svg(svg)"></div>
	</div>
	<Transition>
		<div v-if="selected_element" id="tooltip-wrapper" ref="tooltip" @click="hide_tooltip">
			<BaseTooltip id="tooltip" @close="hide_tooltip">
				<template #header>
					<slot name="header"></slot>
				</template>
				<slot></slot>
			</BaseTooltip>
		</div>
	</Transition>
</template>

<style scoped>
	#wrapper {
		display: flex;

		flex-direction: column;
		justify-content: center;
		overflow-x: auto;

		max-width: 100%;
	}

	#div-svg {
		margin-inline: auto;
	}

	#tooltip-wrapper {
		position: fixed;

		inset: 0;

		backdrop-filter: blur(0.125em);

		display: flex;

		align-items: center;
		justify-content: center;
	}

	#tooltip {
		max-width: 25em;
	}

	.v-enter-active,
	.v-leave-active {
		transition: filter 0.2s;
	}

	.v-enter-from,
	.v-leave-to {
		filter: opacity(0);
	}
</style>

<style>
	svg#main-content {
		width: 50em;
	}

	/* @media screen and (max-width: 900px) {
		svg#main-content {
			max-width: none;
		}
	} */

	svg#main-content * {
		user-select: none;
	}

	svg#main-content .fill {
		cursor: pointer;

		transition: filter 0.1s;
	}

	/* module */
	svg#main-content .module:hover .fill,
	svg#main-content .module:hover.fill {
		filter: var(--filter-module-hover);
	}

	/* element - sold */
	svg#main-content .module.sold .fill,
	svg#main-content .module.sold.fill {
		filter: var(--filter-module-sold);
	}

	svg#main-content .module.sold:hover .fill,
	svg#main-content .module.sold:hover.fill {
		filter: var(--filter-module-sold-hover);
	}

	/* battery */
	svg#main-content .battery:hover .fill,
	svg#main-content .battery:hover.fill {
		filter: var(--filter-battery-hover);
	}

	/* battery - sold */
	svg#main-content .battery.sold .fill,
	svg#main-content .battery.sold.fill {
		filter: var(--filter-battery-sold);
	}

	svg#main-content .battery.sold:hover .fill,
	svg#main-content .battery.sold:hover.fill {
		filter: var(--filter-battery-sold-hover);
	}
</style>
