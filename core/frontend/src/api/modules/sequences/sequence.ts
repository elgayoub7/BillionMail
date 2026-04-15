import { instance } from '@/api'
import { i18n } from '@/i18n'

const { t } = i18n.global

export function getSequenceList(params: { page: number; page_size: number; keyword?: string; status?: number }) {
  return instance.get('/sequence/list', { params })
}

export function getSequenceDetail(params: { id: number }) {
  return instance.get('/sequence/find', { params })
}

export function createSequence(data: Record<string, unknown>) {
  return instance.post('/sequence/create', data, {
    fetchOptions: {
      loading: t('sequences.loading.creating'),
      successMessage: true,
    },
  })
}

export function updateSequence(data: Record<string, unknown>) {
  return instance.post('/sequence/update', data, {
    fetchOptions: {
      loading: t('sequences.loading.updating'),
      successMessage: true,
    },
  })
}

export function deleteSequence(data: { id: number }) {
  return instance.post('/sequence/delete', data, {
    fetchOptions: {
      loading: t('sequences.loading.deleting'),
      successMessage: true,
    },
  })
}

export function activateSequence(data: { id: number }) {
  return instance.post('/sequence/activate', data, {
    fetchOptions: { successMessage: true },
  })
}

export function pauseSequence(data: { id: number }) {
  return instance.post('/sequence/pause', data, {
    fetchOptions: { successMessage: true },
  })
}

export function resumeSequence(data: { id: number }) {
  return instance.post('/sequence/resume', data, {
    fetchOptions: { successMessage: true },
  })
}

export function enrollContacts(data: { sequence_id: number; contact_ids?: number[] }) {
  return instance.post('/sequence/enroll', data, {
    fetchOptions: { successMessage: true },
  })
}

export function getEnrollments(params: { sequence_id: number; page: number; page_size: number; status?: number }) {
  return instance.get('/sequence/enrollments', { params })
}

export function removeEnrollment(data: { enrollment_id: number }) {
  return instance.post('/sequence/remove_enrollment', data, {
    fetchOptions: { successMessage: true },
  })
}

export function sendTestStep(data: { sequence_id: number; step_id: number; test_email: string }) {
  return instance.post('/sequence/send_test', data, {
    fetchOptions: { successMessage: true },
  })
}
