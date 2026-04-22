# BillionMail Cold — Plan d'exécution

## Contexte
BillionMail est un mail server open-source (Go + Vue 3) forké par elgayoub7.
Le fork `dev` contient déjà: séquences multi-step, warmup IP, multi-IP/domain, 25 fixes sécurité.

**Objectif**: Transformer BillionMail en plateforme cold mail lean (~150 emails/jour, multi-domaine, haute délivrabilité).

## Architecture existante (GARDER)
- **Backend**: Go (GoFrame) — binaire compilé
- **Frontend**: Vue 3 + Naive UI + TypeScript
- **DB**: PostgreSQL (61 tables)
- **Infra**: Docker Compose (7 conteneurs)
- **Relay**: Brevo SMTP (port 587) sur VM Oracle

## Ce qui existe DÉJÀ (ne pas recoder)
| Feature | Localisation | État |
|---------|-------------|------|
| Séquences multi-step | `service/sequence/` | ✅ Fonctionnel |
| Warmup IP progressif | `service/warmup/` | ✅ Rate limiter + daily ramp |
| Multi-IP/domain rotation | `service/multi_ip_domain/` | ✅ Config manager |
| Per-recipient analytics | `service/maillog_stat/` | ✅ Opens/clicks/bounces |
| Sécurité durcie | 25 commits | ✅ CORS, JWT, XSS, etc. |
| Contacts + groupes | `contacts/`, `subscribe_list/` | ✅ |
| Templates email | `template/` | ✅ |
| Relay SMTP | Postfix + Brevo | ✅ Sur VM |
| Domain management | `domain/` | ✅ |

## Ce qu'il MANQUE pour cold mail

### P0 — Critique
1. **AB Test Engine** — Split test sur sujet/body, auto winner
2. **Lead Scoring** — Score basé sur engagement (opens, clicks, replies)
3. **Cold Dashboard** — Vue simplifiée orientée outreach
4. **Bounce/Reply handler** — Auto-pause sur hard bounce, avance sur reply
5. **Spam check pré-envoi** — Score via Rspamd avant envoi

### P1 — Important
6. **Warmup dashboard** — Visualiser progression warmup par domaine
7. **Séquence conditionnelle avancée** — Branchement sur opens/clicks/replies
8. **Personalization variables** — {{first_name}}, {{company}}, etc. avancées
9. **Health check domaines** — SPF/DKIM/DMARC status en temps réel
10. **Email validation** — Vérifier adresse avant envoi (MX check)

### P2 — Nice to have
11. **LinkedIn enrichment** — Auto-remplir company/title depuis profil
12. **Timezone scheduling** — Envoyer à l'heure locale du prospect
13. **Unified inbox** — Vue replies centralisée
14. **Export/rapport PDF** — Analytics exportables

---

## Phase 1 — Architecture & Schema (2-3h)

### DB Migrations
```sql
-- AB Tests
CREATE TABLE ab_tests (
    id SERIAL PRIMARY KEY,
    sequence_step_id INT NOT NULL,
    name VARCHAR(255),
    status INT DEFAULT 0, -- 0=draft, 1=running, 2=completed
    winner_variant INT, -- 0=A, 1=B
    winner_criteria VARCHAR(50) DEFAULT 'open_rate', -- open_rate, click_rate, reply_rate
    sample_size INT, -- % of audience per variant
    confidence_threshold REAL DEFAULT 0.95,
    created_at INT DEFAULT EXTRACT(EPOCH FROM NOW()),
    completed_at INT
);

CREATE TABLE ab_test_variants (
    id SERIAL PRIMARY KEY,
    ab_test_id INT REFERENCES ab_tests(id),
    variant INT NOT NULL, -- 0=A, 1=B
    subject TEXT,
    body_html TEXT,
    body_text TEXT,
    sent_count INT DEFAULT 0,
    open_count INT DEFAULT 0,
    click_count INT DEFAULT 0,
    reply_count INT DEFAULT 0,
    bounce_count INT DEFAULT 0
);

-- Lead Scoring
CREATE TABLE lead_scores (
    id SERIAL PRIMARY KEY,
    contact_email TEXT NOT NULL UNIQUE,
    score INT DEFAULT 0,
    engagement_level VARCHAR(20) DEFAULT 'cold', -- cold, warm, hot, converted
    total_opens INT DEFAULT 0,
    total_clicks INT DEFAULT 0,
    total_replies INT DEFAULT 0,
    total_bounces INT DEFAULT 0,
    last_engagement_at INT,
    last_scored_at INT DEFAULT EXTRACT(EPOCH FROM NOW()),
    sequence_id INT,
    metadata JSONB DEFAULT '{}'
);

-- Sequence conditions avancées
ALTER TABLE sequence_steps ADD COLUMN IF NOT EXISTS condition_type VARCHAR(50) DEFAULT 'time';
-- time, opened, clicked, replied, bounced, not_opened, not_clicked, not_replied

-- Bounce tracking
CREATE TABLE bounce_rules (
    id SERIAL PRIMARY KEY,
    rule_type VARCHAR(50) NOT NULL, -- hard_bounce, soft_bounce, complaint
    pattern VARCHAR(255),
    action VARCHAR(50) DEFAULT 'pause_contact', -- pause_contact, pause_sequence, notify
    created_at INT DEFAULT EXTRACT(EPOCH FROM NOW())
);

-- Spam scores pré-envoi
CREATE TABLE email_spam_scores (
    id SERIAL PRIMARY KEY,
    email_task_id INT,
    score REAL,
    report JSONB,
    checked_at INT DEFAULT EXTRACT(EPOCH FROM NOW())
);

-- Domain health
CREATE TABLE domain_health (
    id SERIAL PRIMARY KEY,
    domain VARCHAR(255) NOT NULL,
    spf_status VARCHAR(20), -- pass, fail, missing
    dkim_status VARCHAR(20),
    dmarc_status VARCHAR(20),
    mx_status VARCHAR(20),
    last_checked INT DEFAULT EXTRACT(EPOCH FROM NOW())
);
```

## Phase 2 — Backend Go (4-6h)

### Story 2.1: AB Test Engine
- Fichiers: `service/abtest/abtest.go`, `service/abtest/abtest_engine.go`
- API: `controller/abtest/abtest.go`
- Logique:
  1. Créer test A/B sur un step de séquence
  2. Split recipients 50/50 (ou configurable)
  3. Tracker opens/clicks/replies par variante
  4. Auto-déterminer winner quand statistiquement significatif
  5. Continuer avec la variante gagnante

### Story 2.2: Lead Scoring Engine
- Fichiers: `service/scoring/scoring.go`, `service/scoring/scoring_engine.go`
- API: `controller/scoring/scoring.go`
- Logique:
  1. Score = (opens×2 + clicks×5 + replies×10) - (bounces×5)
  2. Niveaux: cold (0-5), warm (6-20), hot (21-50), converted (51+)
  3. Hook dans `maillog_stat/tracker.go` pour update en temps réel
  4. API GET /scoring/leads?page=1&level=hot

### Story 2.3: Bounce Handler Amélioré
- Fichiers: `service/bounce/bounce_handler.go`
- Hook dans `sequence_event_handler.go`
- Logique:
  1. Parse logs Postfix → identifier hard vs soft bounce
  2. Hard bounce → auto-désinscrire de toutes séquences actives
  3. Complaint → flagger contact + notify
  4. Soft bounce → retry + pause si 3 soft bounces consécutifs

### Story 2.4: Spam Check Pré-envoi
- Fichiers: `service/spamcheck/spamcheck.go`
- Logique:
  1. Appeler Rspamd API locale pour scorer le contenu
  2. Si score > 6/10 → bloquer envoi + alerter
  3. Afficher suggestions (remove link, reduce images, etc.)

### Story 2.5: Séquence Conditions Avancées
- Modification: `service/sequence/sequence_engine.go`
- Ajouter condition_type: opened, clicked, replied, not_opened, not_clicked
- Branchement: si step N = "attend reply" → si reply détecté → avancer, sinon → step suivant après délai

### Story 2.6: Domain Health Check
- Fichiers: `service/domainhealth/domainhealth.go`
- Logique:
  1. DNS lookups: SPF, DKIM, DMARC, MX
  2. Stocker résultat dans domain_health
  3. Cron toutes les 6h
  4. Alerte si un record change ou est manquant

## Phase 3 — Frontend Vue (4-6h)

### Story 3.1: Cold Dashboard (`views/cold-dashboard/`)
- Vue simplifiée avec:
  - Séquences actives (cards)
  - Reply rate, open rate, bounce rate
  - Lead scoring heatmap
  - Warmup status par domaine
  - Problèmes (domaines en erreur, bounces élevés)

### Story 3.2: AB Test UI (`views/sequences/edit.vue` modification)
- Onglet "A/B Test" dans l'éditeur de step
- Variantes côte à côte
- Métriques en temps réel
- Sélection manuelle du winner

### Story 3.3: Lead Scoring UI (`views/contacts/` modification)
- Colonne score dans la liste contacts
- Filtre par niveau (cold/warm/hot)
- Détail contact avec timeline d'engagement

### Story 3.4: Domain Health UI (`views/domain/` modification)
- Status visuel (🟢🟡🔴) pour SPF/DKIM/DMARC/MX
- Bouton "Re-check"
- Alerte si problème

### Story 3.5: Séquence Conditions UI (`views/sequences/edit.vue` modification)
- Dropdown type de condition
- Visualisation branchée (flow chart simple)

## Phase 4 — Tests & Déploiement VM (2-3h)

### Story 4.1: Tests unitaires
- AB test engine
- Scoring engine
- Bounce handler

### Story 4.2: Tests d'intégration VM
- Deploy sur Oracle VM
- Test envoi via Brevo relay
- Vérifier tracking opens/clicks
- Tester séquence multi-step end-to-end
- Tester AB test avec vrais emails

### Story 4.3: Smoke tests
- Créer séquence 3 steps
- Importer 10 contacts test
- Vérifier warmup ramp
- Vérifier scoring
- Vérifier bounce handling

## Ordre d'exécution (critique)

```
Phase 1 (Schema)     →  2h  →  Bloquant pour tout
Phase 2.1 (AB Test)  →  2h  →  Peut être en parallèle avec 2.2
Phase 2.2 (Scoring)  →  1.5h →  Peut être en parallèle avec 2.1
Phase 2.3 (Bounce)   →  1h  →  Dépend de Phase 1
Phase 2.4 (Spam)     →  1h  →  Indépendant
Phase 2.5 (Conditions)→ 1.5h →  Dépend de Phase 1
Phase 2.6 (Health)   →  1h  →  Indépendant
Phase 3.1 (Dashboard)→  2h  →  Dépend de 2.1, 2.2
Phase 3.2-3.5 (UI)   →  3h  →  Dépend des backends respectifs
Phase 4 (Tests+Deploy)→ 3h  →  Tout doit être mergé
```

**Total estimé**: ~18h de travail agent
**4 agents en parallèle** → ~5h temps réel
