import type { Config } from 'tailwindcss'

export default {
	content: ['./src/**/*.{vue,ts,tsx,jsx,html}'],
	darkMode: 'class',
	theme: {
		extend: {
			colors: {
				surface: {
					950: '#0a0b10',
					900: '#0f1117',
					800: '#1a1d27',
					700: '#232630',
					600: '#2a2d3a',
					500: '#363a4a',
				},
				accent: {
					DEFAULT: '#6c5ce7',
					hover: '#7d6ff0',
					light: '#a29bfe',
				},
			},
			fontFamily: {
				sans: ['Inter', 'system-ui', '-apple-system', 'sans-serif'],
			},
		},
	},
	// Prefix to avoid conflicts with Naive UI classes
	corePlugins: {
		preflight: false, // Naive UI + our SCSS handles resets
	},
	plugins: [],
} satisfies Config
