import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { isDev } from '@/utils'

// Routes reflect list
const routesReflectList = [
	'Overview',
	'Inbox',
	'template',
	'Send API',
	'Contacts',
	'Sequences',
	'Leads',
	'Enrichment',
	'MailDomain',
	'MailBoxes',
	'SMTP',
	'Logs',
	'Settings',
	'Automation',
	'Video Outreach',
	'Analytics',
]

// Explicitly import route modules (import.meta.glob/webpackContext broken by source.define)
import overview from './modules/overview'
import coldDashboard from './modules/cold-dashboard'
import sequences from './modules/sequences'
import template from './modules/template'
import contacts from './modules/contacts'
import mailbox from './modules/mailbox'
import domain from './modules/domain'
import smtp from './modules/smtp'
import logs from './modules/logs'
import settings from './modules/settings'
import api from './modules/api'
import automation from './modules/automation'
import inbox from './modules/inbox'
import domainHealth from './modules/domain-health'
import scoring from './modules/scoring'
import videoOutreach from './modules/video-outreach'
import analytics from './modules/analytics'
import market from './modules/market'
import deliverability from './modules/deliverability'

const allModules: RouteRecordRaw[] = [
	overview,
	coldDashboard,
	sequences,
	template,
	contacts,
	mailbox,
	domain,
	smtp,
	logs,
	settings,
	api,
	automation,
	inbox,
	domainHealth,
	scoring,
	videoOutreach,
	analytics,
	market,
	deliverability,
]

// Module routes
export let menuList: RouteRecordRaw[] = allModules.filter(Boolean)

// Sort module routes
menuList = menuList.reduce((p: RouteRecordRaw[], v: RouteRecordRaw) => {
	const routeIndex = routesReflectList.findIndex(item => item == v.meta?.title)
	if (routeIndex >= 0) p[routeIndex] = v
	return p
}, [] as RouteRecordRaw[]).filter(Boolean) as RouteRecordRaw[]

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
