# Phase 2 - Séquences de Relance Automatisées

## Date : 2026-04-15

## Résumé

Ajout d'un système complet de séquences de relance automatisées type Lemlist. Chaque étape d'une séquence crée un `email_task` standard traité par le `TaskExecutor` existant. Un Sequence Engine (timer 30s) avance les enrollments et évalue les conditions.

---

## Fichiers Créés (11)

### Backend Go

| Fichier | Description |
|---------|-------------|
| `core/internal/service/database_initialization/sequence.go` | 4 tables + 9 indexes |
| `core/internal/model/entity/sequence.go` | 4 entity structs (Sequence, SequenceStep, SequenceEnrollment, SequenceEmailTask) |
| `core/api/sequence/v1/sequence.go` | 12 endpoints (req/res types GoFrame v2) |
| `core/api/sequence/sequence.go` | Interface ISequenceV1 |
| `core/internal/controller/sequence/sequence.go` | ControllerV1 avec 12 handlers |
| `core/internal/service/sequence/sequence.go` | CRUD + enrollment + helpers |
| `core/internal/service/sequence/sequence_engine.go` | Moteur timer (advance + conditions + batch emails) |
| `core/internal/service/sequence/sequence_event_handler.go` | RecordEvent pour open/click |

### Frontend Vue/TS

| Fichier | Description |
|---------|-------------|
| `core/frontend/src/router/modules/sequences.ts` | 4 routes (list, create, edit, detail) |
| `core/frontend/src/api/modules/sequences/sequence.ts` | API client (12 fonctions) |
| `core/frontend/src/views/sequences/interface.ts` | Types TypeScript |
| `core/frontend/src/views/sequences/index.vue` | Liste des séquences |
| `core/frontend/src/views/sequences/edit.vue` | Builder d'étapes (email/wait/condition) |
| `core/frontend/src/views/sequences/detail.vue` | Détail + enrollments (2 tabs) |

---

## Fichiers Modifiés (5)

### `core/internal/cmd/cmd.go`
- **Ligne ~25** : Import `sequence_ctrl "billionmail-core/internal/controller/sequence"`
- **Ligne ~271** : `sequence_ctrl.NewV1()` ajouté au `group.Bind(...)`

### `core/internal/service/middlewares/rbac.go`
- **Ligne 27** : `"sequence"` ajouté au slice `modules` pour les permissions RBAC

### `core/internal/service/timers/timers.go`
- **Ligne 5** : Import `"billionmail-core/internal/service/sequence"`
- **Après ligne 100** : 2 timers ajoutés :
  - `gtimer.Add(30*time.Second, ...)` → `sequence.ProcessSequenceEnrollments(ctx)`
  - `gtimer.Add(1*time.Minute, ...)` → `sequence.UpdateEnrollmentsFromCompletedTasks(ctx)`

### `core/internal/service/maillog_stat/tracker.go`
- **Ligne 5** : Import `"billionmail-core/internal/service/sequence"`
- **Après ligne 81** (case "open") : `sequence.RecordEvent(ctx, "open", data.Recipient, data.CampaignId, data.MessageId)`
- **Après ligne 128** (case "click") : `sequence.RecordEvent(ctx, "click", data.Recipient, data.CampaignId, data.MessageId)`

### `core/frontend/src/i18n/lang/en.json`
- Section `"sequences"` étendue : ~60 nouvelles clés i18n (status, form, conditions, enrollments, stats, loading)

---

## Base de Données - 4 Tables

### `bm_sequences`
Définitions des séquences.
```sql
id SERIAL PRIMARY KEY,
name VARCHAR(255) UNIQUE,
description TEXT,
status SMALLINT DEFAULT 0,          -- 0:draft, 1:active, 2:paused, 3:archived
addresser VARCHAR(320),             -- email expéditeur
full_name VARCHAR(255),             -- nom affiché
group_id INTEGER,                   -- groupe de contacts
tag_ids TEXT,                       -- filtre tags (JSON)
tag_logic VARCHAR(10) DEFAULT 'AND',
track_open SMALLINT DEFAULT 1,
track_click SMALLINT DEFAULT 1,
unsubscribe SMALLINT DEFAULT 1,
total_enrolled / total_completed / total_unsubscribed / total_bounced INTEGER,
create_time / update_time INTEGER
```

### `bm_sequence_steps`
Étapes d'une séquence.
```sql
id SERIAL PRIMARY KEY,
sequence_id INTEGER REFERENCES bm_sequences(id) ON DELETE CASCADE,
step_order INTEGER,
UNIQUE(sequence_id, step_order),
step_type VARCHAR(50) DEFAULT 'email',  -- 'email', 'wait', 'condition'
subject TEXT,                         -- pour email
template_id INTEGER,                  -- pour email
wait_days / wait_hours INTEGER,       -- pour wait
condition_type VARCHAR(50),           -- 'opened', 'not_opened', 'clicked', 'not_clicked', 'bounced'
condition_step_id INTEGER,            -- étape à vérifier
on_true_go_to / on_false_go_to INTEGER, -- branching
sent_count / opened_count / clicked_count / bounced_count INTEGER
```

### `bm_sequence_enrollments`
Contacts inscrits dans une séquence.
```sql
id SERIAL PRIMARY KEY,
sequence_id INTEGER REFERENCES bm_sequences(id) ON DELETE CASCADE,
contact_id INTEGER,
email VARCHAR(320),
UNIQUE(sequence_id, email),
current_step INTEGER DEFAULT 1,       -- étape actuelle (1-indexed)
status SMALLINT DEFAULT 0,            -- 0:active, 1:completed, 2:paused, 3:exited
enrolled_at / current_step_entered_at / last_email_sent_at / completed_at INTEGER,
total_emails_sent / total_opens / total_clicks INTEGER
```

### `bm_sequence_email_tasks`
Liaison entre étapes séquence et email_tasks existants.
```sql
id SERIAL PRIMARY KEY,
sequence_id / enrollment_id / step_id INTEGER,
email_task_id INTEGER,                -- FK vers email_tasks.id
contact_email VARCHAR(320),
status SMALLINT DEFAULT 0,            -- 0:pending, 1:sent, 2:failed
sent_at INTEGER, message_id TEXT
```

---

## API Endpoints (12)

| Méthode | Chemin | Description |
|---------|--------|-------------|
| GET | `/sequence/list` | Lister séquences (pagination, filtres) |
| GET | `/sequence/find` | Détail séquence + étapes |
| POST | `/sequence/create` | Créer séquence + étapes |
| POST | `/sequence/update` | Modifier séquence + étapes |
| POST | `/sequence/delete` | Supprimer séquence |
| POST | `/sequence/activate` | Activer (draft → active) |
| POST | `/sequence/pause` | Mettre en pause |
| POST | `/sequence/resume` | Reprendre |
| POST | `/sequence/enroll` | Inscrire des contacts |
| GET | `/sequence/enrollments` | Lister enrollments |
| POST | `/sequence/remove_enrollment` | Retirer un contact |
| POST | `/sequence/send_test` | Envoyer email test |

---

## Architecture du Sequence Engine

```
┌─────────────────────────────────────────────────┐
│           Timer (toutes les 30s)                 │
│  ProcessSequenceEnrollments(ctx)                │
│                                                  │
│  1. advanceReadyEnrollments()                   │
│     → Trouve enrollments dont le wait est        │
│       expiré, passe à l'étape suivante          │
│                                                  │
│  2. evaluateConditionSteps()                     │
│     → Évalue conditions (opened, clicked,        │
│       bounced) via tables mailstat               │
│                                                  │
│  3. createBatchStepEmails()                      │
│     → Crée un email_task par (séquence, étape)  │
│     → Insert recipients dans recipient_info      │
│     → Le timer ProcessEmailTasks existant        │
│       prend le relais automatiquement            │
└─────────────────────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────┐
│         Timer (toutes les 60s)                  │
│  UpdateEnrollmentsFromCompletedTasks(ctx)       │
│                                                  │
│  → Marque séquence_email_tasks comme sent        │
│  → Incrémente compteurs enrollment              │
│  → Détecte bounces → marque enrollment exited   │
└─────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────┐
│         Events (tracker.go)                      │
│                                                  │
│  case "open":  → sequence.RecordEvent("open")   │
│  case "click": → sequence.RecordEvent("click")  │
│                                                  │
│  → Incrémente total_opens / total_clicks         │
│    dans bm_sequence_enrollments                  │
└─────────────────────────────────────────────────┘
```

---

## Types d'Étapes

### Email
Envoie un email avec template et subject. Crée un `email_task` standard.

### Wait
Attend X jours et Y heures avant de passer à l'étape suivante. Le timer `advanceReadyEnrollments` vérifie toutes les 30s.

### Condition
Branche conditionnelle basée sur les events d'une étape précédente :
- `opened` / `not_opened` → vérifie `mailstat_opened`
- `clicked` / `not_clicked` → vérifie `mailstat_clicked`
- `bounced` → vérifie `mailstat_send_mails` (status='bounced')

Le branching se configure via `on_true_go_to` et `on_false_go_to`.

---

## Frontend

### Routes
- `/sequences` → Liste des séquences
- `/sequences/create` → Création
- `/sequences/:id/edit` → Modification
- `/sequences/:id` → Détail (2 tabs : étapes + enrollments)

### Composants
- **index.vue** : Table avec filtres (status, recherche), actions (activate/pause/resume/edit/delete)
- **edit.vue** : Builder vertical avec step cards (ajout/suppression/réordonnancement), formulaire séquence + step types
- **detail.vue** : 2 tabs — Steps (stats par étape) + Enrollments (table paginée avec remove)

### Menu
L'entrée "Sequences" apparaît automatiquement dans la sidebar (icon `i-mdi-email-sync-outline` déjà mappé, `routesReflectList` contient déjà `'Sequences'`).
