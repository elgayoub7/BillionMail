import DOMPurify from 'dompurify'

/**
 * Sanitize HTML content to prevent XSS attacks.
 * Used with v-html directives to safely render user/AI-generated content.
 */
export function sanitizeHtml(html: string, options?: DOMPurify.Config): string {
    return DOMPurify.sanitize(html, {
        ALLOWED_TAGS: [
            'b', 'i', 'em', 'strong', 'a', 'p', 'br', 'div', 'span',
            'h1', 'h2', 'h3', 'h4', 'h5', 'h6',
            'ul', 'ol', 'li', 'code', 'pre', 'blockquote',
            'table', 'thead', 'tbody', 'tr', 'th', 'td', 'img', 'hr',
        ],
        ALLOWED_ATTR: ['href', 'class', 'style', 'src', 'alt', 'target', 'rel'],
        ...options,
    })
}
