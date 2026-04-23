import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { isDev } from '@/utils'

// Routes reflect list (used for sidebar ordering)
const routesReflectList = [
	'Overview',
	'Cold Dashboard',
	'Email Marketing',
	'template',
	'Send API',
	'Contacts',
	'Sequences',
	'Lead Scoring',
	'AB Tests',
	'MailDomain',
	'Domain Health',
	'MailBoxes',
	'SMTP',
	'Logs',
	'Settings',
	'Automation',
	'Video Outreach',
]

// Static imports for all route modules
// Replaces broken import.meta.webpackContext (only loads 1/17 modules in Rspack)
import overviewRoute from './modules/overview'
import coldDashboardRoute from './modules/cold-dashboard'
import marketRoute from './modules/market'
import templateRoute from './modules/template'
import apiRoute from './modules/api'
import contactsRoute from './modules/contacts'
import sequencesRoute from './modules/sequences'
import scoringRoute from './modules/scoring'
import abtestRoute from './modules/abtest'
import mailDomainRoute from './modules/domain'
import domainHealthRoute from './modules/domain-health'
import mailboxRoute from './modules/mailbox'
import smtpRoute from './modules/smtp'
import logsRoute from './modules/logs'
import settingsRoute from './modules/settings'
import automationRoute from './modules/automation'
import videoOutreachRoute from './modules/video-outreach'

const moduleRoutes: RouteRecordRaw[] = [
	overviewRoute,
	coldDashboardRoute,
	marketRoute,
	templateRoute,
	apiRoute,
	contactsRoute,
	sequencesRoute,
	scoringRoute,
	abtestRoute,
	mailDomainRoute,
	domainHealthRoute,
	mailboxRoute,
	smtpRoute,
	logsRoute,
	settingsRoute,
	automationRoute,
	videoOutreachRoute,
]

// Sort module routes according to routesReflectList order, filter gaps
export let menuList: RouteRecordRaw[] = moduleRoutes
	.reduce((p: RouteRecordRaw[], v: RouteRecordRaw) => {
		const routeIndex = routesReflectList.findIndex(item => item == v.meta!.title)
		if (routeIndex >= 0) {
			p[routeIndex] = v
		}
		return p
	}, [] as RouteRecordRaw[])
	.filter(Boolean)

const otherArray: RouteRecordRaw[] = []

if (isDev) {
	otherArray.push({
		path: '/test',
		name: 'Test',
		component: () => import('@/views/test/index.vue'),
	})
}

export const routes: RouteRecordRaw[] = [
	{
		path: '/login',
		name: 'Login',
		component: () => import('@/views/login/index.vue'),
	},
	{
		path: '/',
		redirect: '/overview',
	},
	...menuList,
	...otherArray,
]

const router = createRouter({
	history: createWebHistory('/'),
	routes,
	strict: false,
	scrollBehavior: () => ({ left: 0, top: 0 }),
})

export default router
