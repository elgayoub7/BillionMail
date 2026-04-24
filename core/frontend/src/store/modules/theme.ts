import { defineStore } from 'pinia'
import { GlobalThemeOverrides } from 'naive-ui'

export default defineStore(
	'ThemeStore',
	() => {
		const theme = ref<'light' | 'dark'>('dark')

		const getCssVar = (name: string) => {
			return getComputedStyle(document.documentElement).getPropertyValue(name).trim()
		}

		const getThemeOverrides = (): GlobalThemeOverrides => {
			return {
				common: {
					lineHeight: 'normal',
					fontSize: '12px',
					fontSizeSmall: '12px',
					fontSizeMedium: '12px',
					fontSizeLarge: '14px',
					borderRadius: '8px',
					baseColor: '#fff',
					textColor1: getCssVar('--color-text-1'),
					textColor2: getCssVar('--color-text-2'),
					primaryColor: getCssVar('--color-primary-1'),
					primaryColorHover: getCssVar('--color-primary-hover-1'),
					successColor: '#00d68f',
					successColorHover: '#00b377',
					warningColor: getCssVar('--color-warning-1'),
					errorColor: getCssVar('--color-error-1'),
				},
				Layout: {
					color: getCssVar('--color-bg-2'),
					textColor: getCssVar('--color-text-1'),
					headerColor: getCssVar('--color-bg-1'),
					siderColor: '#0f1117',
					siderTextColor: getCssVar('--color-text-2'),
				},
				Menu: {
					fontSize: '14px',
					itemTextColor: getCssVar('--color-menu-1'),
					itemIconColor: getCssVar('--color-menu-1'),
					itemColorActive: getCssVar('--color-menu-active-1'),
					itemTextColorActive: getCssVar('--color-menu-active-2'),
					itemIconColorActive: getCssVar('--color-menu-active-2'),
					itemColorActiveHover: getCssVar('--color-menu-active-1'),
					itemTextColorActiveHover: getCssVar('--color-menu-active-2'),
					itemIconColorActiveHover: getCssVar('--color-menu-active-2'),
					itemColorActiveCollapsed: getCssVar('--color-menu-active-1'),
				},
				Card: {
					color: getCssVar('--color-bg-1'),
					borderColor: 'transparent',
					borderRadius: '10px',
				},
				Form: {
					feedbackHeightMedium: '20px',
					feedbackHeightLarge: '22px',
					feedbackFontSizeMedium: '12px',
					feedbackFontSizeLarge: '12px',
					feedbackPadding: '2px 0 0',
					labelFontSizeTopMedium: '13px',
					labelFontSizeLeftMedium: '13px',
					labelPaddingHorizontal: '0 20px 0 0',
					labelFontWeight: '700',
				},
				Radio: {
					buttonColor: getCssVar('--color-radio-1'),
					buttonColorActive: getCssVar('--color-radio-2'),
					buttonTextColor: getCssVar('--color-text-2'),
					buttonTextColorHover: getCssVar('--color-text-1'),
					buttonTextColorActive: getCssVar('--color-text-1'),
					buttonBorderColor: 'transparent',
					buttonBorderColorHover: 'transparent',
					buttonBorderColorActive: 'transparent',
					buttonBoxShadowHover: 'none',
					buttonBoxShadowFocus: 'none',
					buttonBorderRadius: '6px',
					labelPadding: '0 0 0 8px',
				},
				Dialog: {
					contentMargin: '16px 0',
					textColor: getCssVar('--color-text-1'),
				},
				DataTable: {
					thPaddingMedium: '10px',
					tdPaddingMedium: '10px',
					thColor: getCssVar('--color-table-th-1'),
					tdColor: getCssVar('--color-table-td-1'),
					tdColorHover: getCssVar('--color-table-th-1'),
					borderRadius: '8px',
				},
				Breadcrumb: {
					fontSize: '14px',
				},
				Switch: {
					railColorActive: getCssVar('--color-primary-1'),
				},
				Tabs: {
					tabBorderRadius: '6px',
				},
				Progress: {
					textColorLineInner: '#fff',
				},
				Button: {
					color: getCssVar('--color-bg-1'),
					colorHover: getCssVar('--color-bg-1'),
				},
				Input: {
					color: getCssVar('--color-bg-1'),
					colorFocus: getCssVar('--color-bg-1'),
					borderColor: getCssVar('--color-border-1'),
					borderColorHover: getCssVar('--color-primary-1'),
					borderColorFocus: getCssVar('--color-primary-1'),
				},
			}
		}

		const themeOverrides = ref<GlobalThemeOverrides>(getThemeOverrides())

		const setTheme = (val: 'light' | 'dark') => {
			theme.value = val
		}

		const changeTheme = () => {
			const isDarkMode = theme.value === 'dark'
			document.documentElement.setAttribute('theme-mode', isDarkMode ? 'dark' : '')
			// Tailwind dark mode class
			if (isDarkMode) {
				document.documentElement.classList.add('dark')
			} else {
				document.documentElement.classList.remove('dark')
			}
			nextTick(() => {
				themeOverrides.value = getThemeOverrides()
			})
		}

		return {
			theme,
			themeOverrides,
			setTheme,
			changeTheme,
		}
	},
	{
		persist: {
			pick: ['theme'],
		},
	}
)
