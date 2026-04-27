import { instance } from '@/api'

export interface ScoreSubjectResponse {
	score: number
	warnings: string[]
}

export function scoreSubject(subject: string) {
	return instance.post<ScoreSubjectResponse>('/batch_mail/score_subject', { subject })
}
