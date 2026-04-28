import tailwindcss from 'tailwindcss'
import autoprefixer from 'autoprefixer'
import UnoCSS from 'unocss/postcss'

export default {
	plugins: [UnoCSS(), tailwindcss(), autoprefixer()],
}
