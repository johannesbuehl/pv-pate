/** @type {import('tailwindcss').Config} */
export default {
	theme: {
		extend: {
			colors: {
				blue: "hsl(208, 75%, 41%)",
				green: "hsl(95, 59%, 43%)",
				orange: "hsl(36, 100%, 48%)",

			}
		},
	},
	plugins: [],
	content: ["./src/**/*.{vue,js,ts,html}"]
}

