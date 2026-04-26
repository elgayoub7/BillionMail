<template>
  <div class="p-4">
    <n-page-header @back="router.back()" :title="isEdit ? t('sequences.editTitle') : t('sequences.createTitle')" />
    <n-card class="mt-4">
      <n-form ref="formRef" :model="form" label-placement="top">
        <n-grid :cols="2" :x-gap="24">
          <n-form-item-gi :label="t('sequences.form.name')" path="name">
            <n-input v-model:value="form.name" :placeholder="t('sequences.form.namePlaceholder')" />
          </n-form-item-gi>
          <n-form-item-gi :label="t('sequences.form.group')" path="group_id">
            <group-select v-model:value="form.group_id" />
          </n-form-item-gi>
          <n-form-item-gi :span="2" :label="t('sequences.form.description')">
            <n-input v-model:value="form.description" type="textarea" :rows="2" />
          </n-form-item-gi>
        </n-grid>

        <!-- Sender Pool -->
        <n-divider>Senders</n-divider>
        <n-grid :cols="2" :x-gap="24">
          <n-form-item-gi :span="2" label="Sender Pool">
            <n-select
              v-model:value="form.selected_senders"
              :options="mailboxOptions"
              multiple
              filterable
              placeholder="Select sender mailboxes..."
              @update:value="onSenderPoolChange"
            />
          </n-form-item-gi>
          <n-form-item-gi label="Daily limit per sender">
            <n-input-number v-model:value="form.daily_limit_per_sender" :min="0" placeholder="0 = unlimited" class="w-full" />
          </n-form-item-gi>
          <n-form-item-gi label="Sender Display Name">
            <n-input v-model:value="form.full_name" placeholder="Auto-filled from mailbox" />
          </n-form-item-gi>
        </n-grid>

        <!-- Scheduling -->
        <n-divider>Planning</n-divider>
        <n-grid :cols="3" :x-gap="24">
          <n-form-item-gi label="Delay between emails (sec)">
            <n-select
              v-model:value="form.send_delay"
              :options="delayOptions"
              placeholder="Select delay..."
            />
          </n-form-item-gi>
          <n-form-item-gi label="Schedule start hour">
            <n-input-number v-model:value="form.schedule_start_hour" :min="0" :max="23" class="w-full" />
          </n-form-item-gi>
          <n-form-item-gi label="Schedule end hour">
            <n-input-number v-model:value="form.schedule_end_hour" :min="1" :max="24" class="w-full" />
          </n-form-item-gi>
          <n-form-item-gi :span="3" label="Active days">
            <n-checkbox-group v-model:value="form.schedule_days">
              <n-space>
                <n-checkbox v-for="d in dayOptions" :key="d.value" :value="d.value" :label="d.label" />
              </n-space>
            </n-checkbox-group>
          </n-form-item-gi>
        </n-grid>
      </n-form>

      <n-divider>{{ t('sequences.form.steps') }}</n-divider>

      <div class="flex flex-col gap-3">
        <div v-for="(step, index) in form.steps" :key="index" class="relative">
          <div v-if="index > 0" class="absolute left-6 -top-3 w-0.5 h-3 bg-gray-300" />
          <n-card size="small" :title="`${t('sequences.form.step')} ${step.step_order}`" class="ml-8">
            <template #header-extra>
              <n-button text type="error" size="small" @click="removeStep(index)" :disabled="form.steps.length <= 1">
                {{ t('common.actions.delete') }}
              </n-button>
            </template>
            <n-space vertical>
              <n-select
                v-model:value="step.step_type"
                :options="stepTypeOptions"
                @update:value="(val: string) => onStepTypeChange(step, val)"
              />
              <template v-if="step.step_type === 'email'">
                <n-input v-model:value="step.subject" :placeholder="t('sequences.form.subjectPlaceholder')" />
                <n-select
                  v-model:value="step.template_id"
                  :options="templateOptions"
                  filterable
                  placeholder="Select a template..."
                  clearable
                />
              </template>
              <template v-if="step.step_type === 'wait'">
                <n-space>
                  <n-input-number v-model:value="step.wait_days" :min="0" :placeholder="t('sequences.form.days')">
                    <template #suffix>{{ t('sequences.form.days') }}</template>
                  </n-input-number>
                  <n-input-number v-model:value="step.wait_hours" :min="0" :max="23" :placeholder="t('sequences.form.hours')">
                    <template #suffix>{{ t('sequences.form.hours') }}</template>
                  </n-input-number>
                </n-space>
              </template>
              <template v-if="step.step_type === 'condition'">
                <n-select v-model:value="step.condition_type" :options="conditionOptions" />
                <n-input-number v-model:value="step.condition_step_id" :min="1" :placeholder="t('sequences.form.conditionStep')" class="w-full" />
              </template>
            </n-space>
          </n-card>
        </div>
      </div>

      <n-button dashed block class="mt-4" @click="addStep">
        {{ t('sequences.form.addStep') }}
      </n-button>
    </n-card>

    <div class="mt-4 flex justify-end gap-3">
      <n-button @click="router.back()">{{ t('common.actions.cancel') }}</n-button>
      <n-button @click="handleSave">{{ t('sequences.form.saveDraft') }}</n-button>
      <n-button type="primary" :loading="saving" @click="handleSaveAndActivate">
        {{ t('sequences.form.saveAndActivate') }}
      </n-button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { Message } from '@/utils'
import { getSequenceDetail, createSequence, updateSequence, activateSequence } from '@/api/modules/sequences/sequence'
import { getTemplateAll } from '@/api/modules/market/template'
import GroupSelect from '@/views/contacts/subscribers/components/GroupSelect.vue'
import { NFormItem, NInput } from 'naive-ui'
import type { StepInput } from './interface'
import { instance } from '@/api'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const formRef = ref()
const saving = ref(false)

const isEdit = computed(() => !!route.params.id)

const templateOptions = ref<any[]>([])
const mailboxOptions = ref<any[]>([])

const delayOptions = [
  { label: 'No delay', value: 0 },
  { label: '30 seconds', value: 30 },
  { label: '1 minute', value: 60 },
  { label: '2 minutes', value: 120 },
  { label: '5 minutes', value: 300 },
  { label: '10 minutes', value: 600 },
  { label: '15 minutes', value: 900 },
  { label: '20 minutes', value: 1200 },
  { label: '30 minutes', value: 1800 },
  { label: '1 hour', value: 3600 },
]

const dayOptions = [
  { label: 'Mon', value: 1 },
  { label: 'Tue', value: 2 },
  { label: 'Wed', value: 3 },
  { label: 'Thu', value: 4 },
  { label: 'Fri', value: 5 },
  { label: 'Sat', value: 6 },
  { label: 'Sun', value: 7 },
]

const form = reactive({
  name: '',
  description: '',
  addresser: '',
  full_name: '',
  group_id: 0,
  tag_ids: [] as number[],
  tag_logic: 'AND',
  track_open: 1,
  track_click: 1,
  unsubscribe: 1,
  selected_senders: [] as string[],
  sender_pool: '[]',
  daily_limit_per_sender: 0,
  send_delay: 0,
  schedule_start_hour: 0,
  schedule_end_hour: 24,
  schedule_days: [1, 2, 3, 4, 5, 6, 7] as number[],
  steps: [
    { step_order: 1, step_type: 'email', subject: '', template_id: null as number | null },
  ] as StepInput[],
})

async function loadTemplates() {
  try {
    const res = await getTemplateAll()
    const list = Array.isArray(res) ? res : (res?.data || res?.list || [])
    templateOptions.value = list.map((t: any) => ({
      label: t.temp_name || t.name || `Template #${t.id}`,
      value: t.id,
    }))
  } catch { /* ignore */ }
}

async function loadMailboxes() {
  try {
    const res = await instance.get('/mailbox/all')
    const list = Array.isArray(res) ? res : (res?.data || [])
    mailboxOptions.value = list.map((m: any) => ({
      label: `${m.full_name || m.username} <${m.username}>`,
      value: m.username,
      email: m.username,
      full_name: m.full_name || '',
      domain: m.domain || '',
    }))
  } catch { /* ignore */ }
}

function onSenderPoolChange(ids: string[]) {
  if (ids.length === 0) {
    form.sender_pool = '[]'
    form.addresser = ''
    form.full_name = ''
    return
  }
  const selected = mailboxOptions.value.filter(m => ids.includes(m.value as unknown as string))
  const pool = selected.map(m => ({ email: m.email, name: m.full_name || m.email.split('@')[0] }))
  form.sender_pool = JSON.stringify(pool)
  form.addresser = pool[0].email
  form.full_name = pool[0].name
}

const stepTypeOptions = [
  { label: t('sequences.stepType.email'), value: 'email' },
  { label: t('sequences.stepType.wait'), value: 'wait' },
  { label: t('sequences.stepType.condition'), value: 'condition' },
]

const conditionOptions = [
  { label: t('sequences.condition.opened'), value: 'opened' },
  { label: t('sequences.condition.notOpened'), value: 'not_opened' },
  { label: t('sequences.condition.clicked'), value: 'clicked' },
  { label: t('sequences.condition.notClicked'), value: 'not_clicked' },
  { label: t('sequences.condition.bounced'), value: 'bounced' },
]

function addStep() {
  const nextOrder = form.steps.length + 1
  form.steps.push({ step_order: nextOrder, step_type: 'email', subject: '', template_id: null as number | null })
}

function removeStep(index: number) {
  form.steps.splice(index, 1)
  form.steps.forEach((step, i) => {
    step.step_order = i + 1
  })
}

function onStepTypeChange(step: StepInput, type: string) {
  if (type === 'wait') {
    step.wait_days = step.wait_days || 3
    step.wait_hours = 0
  } else if (type === 'email') {
    step.template_id = null as any
    step.subject = ''
  } else if (type === 'condition') {
    step.condition_type = 'opened'
    step.condition_step_id = 1
  }
}

function buildPayload() {
  let addr = form.addresser
  let fn = form.full_name
  if (form.sender_pool && form.sender_pool !== '[]') {
    try {
      const pool = JSON.parse(form.sender_pool)
      if (pool.length > 0) {
        addr = pool[0].email
        fn = pool[0].name
      }
    } catch {}
  }
  return {
    name: form.name,
    description: form.description,
    addresser: addr,
    full_name: fn,
    group_id: form.group_id,
    tag_ids: form.tag_ids,
    tag_logic: form.tag_logic,
    track_open: form.track_open,
    track_click: form.track_click,
    unsubscribe: form.unsubscribe,
    sender_pool: form.sender_pool,
    daily_limit_per_sender: form.daily_limit_per_sender,
    send_delay: form.send_delay,
    schedule_start_hour: form.schedule_start_hour,
    schedule_end_hour: form.schedule_end_hour,
    schedule_days: JSON.stringify(form.schedule_days),
    steps: form.steps.map(s => ({ ...s, template_id: s.template_id || 0 })),
  }
}

async function handleSave() {
  if (!form.name) {
    Message.warning(t('sequences.form.nameRequired'))
    return
  }

  saving.value = true
  try {
    const payload = buildPayload()
    if (isEdit.value) {
      await updateSequence({ id: Number(route.params.id), ...payload })
    } else {
      await createSequence(payload as any)
    }
    Message.success(t('sequences.form.saved'))
    router.push('/sequences')
  } catch {
  } finally {
    saving.value = false
  }
}

async function handleSaveAndActivate() {
  if (!form.name) {
    Message.warning(t('sequences.form.nameRequired'))
    return
  }

  saving.value = true
  try {
    const payload = buildPayload()
    let id: number
    if (isEdit.value) {
      await updateSequence({ id: Number(route.params.id), ...payload })
      id = Number(route.params.id)
    } else {
      const res = await createSequence(payload as any)
      id = res.id
    }
    await activateSequence({ id })
    Message.success(t('sequences.form.activated'))
    router.push('/sequences')
  } catch {
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  loadTemplates()
  loadMailboxes()

  if (isEdit.value) {
    const res = await getSequenceDetail({ id: Number(route.params.id) })
    const detail = res
    let sp: any[] = []
    try { sp = detail.sender_pool ? JSON.parse(detail.sender_pool) : [] } catch {}
    let sd: number[] = [1,2,3,4,5,6,7]
    try { sd = detail.schedule_days ? JSON.parse(detail.schedule_days) : [1,2,3,4,5,6,7] } catch {}

    Object.assign(form, {
      name: detail.name,
      description: detail.description,
      addresser: detail.addresser,
      full_name: detail.full_name,
      group_id: detail.group_id,
      tag_ids: detail.tag_ids || [],
      tag_logic: detail.tag_logic,
      track_open: detail.track_open,
      track_click: detail.track_click,
      unsubscribe: detail.unsubscribe,
      sender_pool: detail.sender_pool || '[]',
      daily_limit_per_sender: detail.daily_limit_per_sender || 0,
      send_delay: detail.send_delay || 0,
      schedule_start_hour: detail.schedule_start_hour ?? 0,
      schedule_end_hour: detail.schedule_end_hour ?? 24,
      schedule_days: sd,
      steps: detail.steps.map((s) => ({
        step_order: s.step_order,
        step_type: s.step_type,
        subject: s.subject || '',
        template_id: s.template_id || null,
        wait_days: s.wait_days || 0,
        wait_hours: s.wait_hours || 0,
        condition_type: s.condition_type || '',
        condition_step_id: s.condition_step_id || 0,
        on_true_go_to: s.on_true_go_to || 0,
        on_false_go_to: s.on_false_go_to || 0,
      })),
    })

    // After mailboxes loaded, match sender pool to mailbox IDs
    if (sp.length > 0) {
      const interval = setInterval(() => {
        if (mailboxOptions.value.length === 0) return
        clearInterval(interval)
        const ids = mailboxOptions.value
          .filter(m => sp.some((s: any) => s.email === m.email))
          .map(m => m.value as unknown as string)
        form.selected_senders = ids
      }, 300)
    }
  }
})
</script>