# BillionMail Cold — Agent Context

## Codebase Patterns (SUIVRE EXACTEMENT)

### 1. DB Initialization
- Fichier: `core/internal/service/database_initialization/<feature>.go`
- Pattern: Fonction dans `init()`, utilise `g.DB().Exec(ctx, sql)` (GoFrame)
- Tables préfixées `bm_` pour les nouvelles (ex: `bm_sequence_steps`)
- Appelé depuis `database_initialization.go` → ajouter dans la liste d'init

### 2. API Interface
- Fichier: `core/api/<feature>/<feature>.go`
- Pattern: Interface `IFeatureV1` avec méthodes typées
- Request/Response structs dans `core/api/<feature>/<feature>.go`
- Utilise `g.Meta` pour les tags de route

### 3. Controller
- Fichier: `core/internal/controller/<feature>/<feature>.go`
- Pattern: `ControllerV1` qui implémente l'interface API
- Appelle le service layer, jamais la DB directement

### 4. Service Layer
- Fichier: `core/internal/service/<feature>/<feature>.go`
- Pattern: Singleton via `var service = New()` + fonction d'accès
- Utilise `g.DB()` pour les queries (GoFrame ORM)
- Context passé partout

### 5. Entity/Model
- Fichier: `core/internal/model/entity/<feature>.go`
- Pattern: Structs Go avec tags `json` et `dc` (description)
- Pas de GORM — GoFrame gdb

### 6. Router Registration
- Fichier: `core/internal/cmd/cmd.go`
- Pattern: Ajouter les routes dans le group approprié

### 7. Frontend
- Vue 3 + Naive UI + TypeScript
- Fichier: `core/frontend/src/views/<feature>/`
- Router: `core/frontend/src/router/index.ts`
- API calls via axios

### 8. Timer/Scheduler
- Fichier: `core/internal/service/timers/timers.go`
- Pattern: Cron jobs pour les tâches périodiques

## Conventions
- Imports GoFrame: `github.com/gogf/gf/v2/...`
- DB: `g.DB()` (GoFrame)
- Log: `g.Log()`
- Config: `g.Cfg()`
- Errors: `gerror.New()` de GoFrame
- Frontend: Composition API, `<script setup lang="ts">`

## Structure DB existante
- 61 tables dans `public` schema
- Contacts: `contacts`, `contact_groups`, `contact_group_detail`
- Email: `email_tasks`, `recipient_info`, `mailstat_send_mails`
- Séquences (fork): `bm_sequences`, `bm_sequence_steps`, `bm_sequence_enrollments`, `bm_sequence_email_tasks`
- Warmup: `sender_ip_warmup`, `sender_ip_mail_provider`

## VM Oracle (pour tests)
- IP: 141.253.118.83
- SSH: ubuntu avec clé `C:\Users\ayoub\Documents\Ubuntu Key\ssh-key-2026-04-20.key`
- Brevo relay actif sur port 587
- Port 25 bloqué en local mais OK via Brevo sur VM
