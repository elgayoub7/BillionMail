<template>
  <div class="inbox-placement-page">
    <n-h1>Inbox Placement</n-h1>

    <n-space vertical :size="20">
      <!-- Controls -->
      <n-space>
        <n-button type="primary" @click="triggerTest">Test Now</n-button>
        <n-button @click="showAddSeed = true">Add Seed</n-button>
      </n-space>

      <!-- Seed List -->
      <n-card title="Seed List" size="small">
        <n-data-table :columns="seedColumns" :data="seeds" :pagination="false" />
      </n-card>

      <!-- Scores Table -->
      <n-card title="Inbox Rates by Sender" size="small">
        <n-data-table :columns="scoreColumns" :data="scores" :pagination="false" />
      </n-card>
    </n-space>

    <!-- Add Seed Modal -->
    <n-modal v-model:show="showAddSeed">
      <n-card style="width: 400px" title="Add Seed Email" :bordered="false">
        <n-form>
          <n-form-item label="Email">
            <n-input v-model:value="newSeed.email" type="text" placeholder="test@gmail.com" />
          </n-form-item>
          <n-form-item label="Provider">
            <n-select v-model:value="newSeed.provider" :options="providerOptions" />
          </n-form-item>
        </n-form>
        <template #footer>
          <n-space justify="end">
            <n-button @click="showAddSeed = false">Cancel</n-button>
            <n-button type="primary" @click="addSeed">Add</n-button>
          </n-space>
        </template>
      </n-card>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { NButton, NDataTable, NCard, NSpace, NModal, NForm, NFormItem, NInput, NSelect, NTag, NH1 } from 'naive-ui'

const seeds = ref([])
const scores = ref([])
const showAddSeed = ref(false)
const newSeed = ref({ email: '', provider: 'gmail' })

const providerOptions = [
  { label: 'Gmail', value: 'gmail' },
  { label: 'Outlook', value: 'outlook' },
  { label: 'Yahoo', value: 'yahoo' },
  { label: 'iCloud', value: 'icloud' },
]

const seedColumns = [
  { title: 'Email', key: 'email' },
  { title: 'Provider', key: 'provider' },
  {
    title: 'Active',
    key: 'active',
    render(row: any) {
      return h(NTag, { type: row.active ? 'success' : 'default' }, () => row.active ? 'Active' : 'Inactive')
    }
  },
  {
    title: 'Actions',
    key: 'actions',
    render(row: any) {
      return h(NButton, { size: 'small', type: 'error', onClick: () => removeSeed(row.email) }, () => 'Remove')
    }
  },
]

const scoreColumns = [
  { title: 'Sender', key: 'username' },
  { title: 'Email', key: 'email' },
  {
    title: 'Inbox Rate',
    key: 'inbox_rate',
    render(row: any) {
      const pct = row.inbox_rate + '%'
      const type = row.status === 'good' ? 'success' : row.status === 'warning' ? 'warning' : 'error'
      return h(NTag, { type }, () => pct)
    }
  },
  { title: 'Status', key: 'status' },
  { title: 'Last Test', key: 'last_test' },
]

async function loadData() {
  const [seedsRes, scoresRes] = await Promise.all([
    fetch('/api/deliverability/seeds'),
    fetch('/api/deliverability/scores'),
  ])
  seeds.value = await seedsRes.json()
  scores.value = await scoresRes.json()
}

async function triggerTest() {
  await fetch('/api/deliverability/trigger-test', { method: 'POST' })
}

async function addSeed() {
  await fetch('/api/deliverability/seeds', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(newSeed.value),
  })
  showAddSeed.value = false
  newSeed.value = { email: '', provider: 'gmail' }
  await loadData()
}

async function removeSeed(email: string) {
  await fetch(`/api/deliverability/seeds/${encodeURIComponent(email)}`, { method: 'DELETE' })
  await loadData()
}

onMounted(loadData)

// Auto-refresh every 60s
setInterval(loadData, 60000)
</script>