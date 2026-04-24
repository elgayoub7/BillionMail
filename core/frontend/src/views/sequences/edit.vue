<template>
  <div class="p-4">
    <n-page-header @back="router.back()" :title="isEdit ? t('sequences.editTitle') : t('sequences.createTitle')" />
    <n-card class="mt-4">
      <n-form ref="formRef" :model="form" label-placement="top">
        <div style="margin-bottom: 16px;">
          <div style="font-size: 14px; margin-bottom: 4px; padding: 0 12px 0 0; line-height: 2; color: var(--n-text-color);">Sender Email</div>
          <div style="position: relative; cursor: text;">
            <input type="text" v-model="form.addresser" placeholder="sender@example.com"
              style="width: 100%; height: 34px; padding: 0 12px; font-size: 14px;
              border: 1px solid var(--n-border-color); border-radius: 3px;
              background: var(--n-color); color: var(--n-text-color);
              outline: none; box-sizing: border-box;"
              onfocus="this.style.borderColor='var(--n-primary-color)'"
              onblur="this.style.borderColor='var(--n-border-color)'" />
          </div>
        </div>
        <n-grid :cols="2" :x-gap="24">
          <n-form-item-gi :label="t('sequences.form.name')" path="name">
            <n-input v-model:value="form.name" :placeholder="t('sequences.form.namePlaceholder')" />
          </n-form-item-gi>
          <n-form-item-gi :label="t('sequences.form.group')" path="group_id">
            <group-select v-model:value="form.group_id" />
          </n-form-item-gi>
          <n-form-item-gi :label="t('sequences.form.senderName')">
            <n-input v-model:value="form.full_name" :placeholder="t('sequences.form.senderNamePlaceholder')" />
          </n-form-item-gi>
          <n-form-item-gi :span="2" :label="t('sequences.form.description')">
            <n-input v-model:value="form.description" type="textarea" :rows="2" />
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
                <n-input-number v-model:value="step.template_id" :placeholder="t('sequences.form.templatePlaceholder')" :show-button="false" class="w-full" />
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
import GroupSelect from '@/views/contacts/subscribers/components/GroupSelect.vue'
import { NFormItem, NInput } from 'naive-ui'
import type { StepInput } from './interface'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const formRef = ref()
const saving = ref(false)

const isEdit = computed(() => !!route.params.id)

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
  steps: [
    { step_order: 1, step_type: 'email', subject: '', template_id: 0 },
  ] as StepInput[],
})

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
  form.steps.push({ step_order: nextOrder, step_type: 'email', subject: '', template_id: 0 })
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
    step.template_id = 0
    step.subject = ''
  } else if (type === 'condition') {
    step.condition_type = 'opened'
    step.condition_step_id = 1
  }
}

async function handleSave() {
  if (!form.name) {
    Message.warning(t('sequences.form.nameRequired'))
    return
  }

  saving.value = true
  try {
    if (isEdit.value) {
      await updateSequence({ id: Number(route.params.id), ...form })
    } else {
      await createSequence(form as any)
    }
    Message.success(t('sequences.form.saved'))
    router.push('/sequences')
  } catch {
    // Error handled by fetchOptions
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
    let id: number
    if (isEdit.value) {
      await updateSequence({ id: Number(route.params.id), ...form })
      id = Number(route.params.id)
    } else {
      const res = await createSequence(form as any)
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
  if (isEdit.value) {
    const res = await getSequenceDetail({ id: Number(route.params.id) })
    const detail = res
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
      steps: detail.steps.map((s) => ({
        step_order: s.step_order,
        step_type: s.step_type,
        subject: s.subject || '',
        template_id: s.template_id,
        wait_days: s.wait_days || 0,
        wait_hours: s.wait_hours || 0,
        condition_type: s.condition_type || '',
        condition_step_id: s.condition_step_id || 0,
        on_true_go_to: s.on_true_go_to || 0,
        on_false_go_to: s.on_false_go_to || 0,
      })),
    })
  }
})
</script>
