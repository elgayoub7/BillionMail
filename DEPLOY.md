# BillionMail - Guide de Deploiement

## Pre-requis

- **Serveur** avec Docker + Docker Compose (v2+)
- **Domaine** avec DNS configure (MX, A records)
- **2 Go RAM minimum**, 4 Go recommande
- **Ports ouverts**: 25, 80, 443, 465, 587

---

## Option 1: Oracle Cloud Free Tier (ARM)

Oracle Cloud offre un serveur ARM gratuit (4 coeurs, 24 Go RAM) - **PARFAIT** pour BillionMail.

### 1. Creer le serveur

1. Aller sur [cloud.oracle.com](https://cloud.oracle.com)
2. Creer un compte ( carte bancaire requise, pas de frais)
3. **Compute > Instances > Create Instance**
4. Image: **Canonical Ubuntu 24.04** (ARM)
5. Shape: **A1 Flex** (4 cores, 24 Go RAM - gratuit)
6. SSH Key: generer ou uploader votre cle
7. Reseaux: ouvrir ports 22, 25, 80, 443, 465, 587 dans le Security List

### 2. Configurer le DNS

Chez votre registrar (Cloudflare, Namecheap, etc.):

```
Type    Name                    Value                                    TTL
A       mail                    VOTRE_IP_PUBLIQUE                        3600
MX      @                       mail.votredomaine.com                    3600  (priorite 10)
TXT     @                       v=spf1 ip4:VOTRE_IP -all                3600
TXT     _dmarc                  v=DMARC1; p=none; rua=mailto:dmarc@ votredomaine.com  3600
```

> **DKIM** : genere automatiquement au premier lancement par OpenDKIM.
> Le record DNS sera affiche dans les logs : `./deploy.sh logs postfix-billionmail | grep DKIM`
> Copiez la cle TXT dans votre zone DNS (nom: `default._domainkey`, valeur: la cle publique).
>
> **Reverse DNS (PTR)** : Chez Oracle Cloud, le PTR est automatiquement configure.
> Verifiez : `dig -x VOTRE_IP` — doit renvoyer `mail.votredomaine.com`.

### 3. Installer Docker sur le serveur

```bash
ssh ubuntu@VOTRE_IP

# Installer Docker
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER
exit
# Reconnecter
ssh ubuntu@VOTRE_IP

# Verifier
docker --version
docker compose version
```

### 4. Deployer BillionMail

```bash
# Cloner le repo
git clone https://github.com/elgayoub7/BillionMail.git
cd BillionMail

# Initialiser la config
chmod +x deploy.sh
./deploy.sh init

# Editer la config
nano .env
```

**Champs OBLIGATOIRES a modifier dans `.env`:**

```env
ADMIN_USERNAME=VotreLoginAdmin
ADMIN_PASSWORD=UnMotDePasseFort123!
SafePath=UnCheminSecret         # URL secrete pour acceder au panel

BILLIONMAIL_HOSTNAME=mail.votredomaine.com    # IMPORTANT

DBPASS=CodeGenereeAleatoirement123
REDISPASS=AutreCodeGeneree456
JWT_SECRET=openssl_rand_hex_32_ICI

# Ports - laisser par defaut sauf si conflit
SMTP_PORT=25
HTTP_PORT=80
HTTPS_PORT=443
```

Generer les mots de passe:
```bash
echo "DBPASS: $(openssl rand -hex 16)"
echo "REDISPASS: $(openssl rand -hex 16)"
echo "JWT_SECRET: $(openssl rand -hex 32)"
```

### 5. Build + Lancer

```bash
# Build l'image (10-15 min la premiere fois)
./deploy.sh build

# Lancer tous les services
./deploy.sh up
```

### 6. Verifier

```bash
# Voir les logs
./deploy.sh logs core-billionmail

# Checker les services
docker compose -f docker-compose.custom.yml ps
```

Acceder au panel:
```
https://VOTRE_IP/billion      ( SafePath par defaut )
```

---

## Option 2: Autres hebergements gratuits / pas cher

| Fournisseur | Prix | Specs | ARM? | Mail ports? |
|-------------|------|-------|------|-------------|
| **Oracle Cloud** | Gratuit | 4C/24Go ARM | Oui | Oui (25,465,587) |
| **Hetzner** | ~4 EUR/mois | 2C/4Go x86 | Non | Oui |
| **Contabo** | ~5 EUR/mois | 4C/6Go x86 | Non | Oui |
| **DigitalOcean** | ~6 USD/mois | 1C/1Go x86 | Non | Oui |
| **RunCloud + VPS** | Variable | Variable | Depends | Oui |

**Oracle Cloud est le meilleur choix** - 4 coeurs ARM + 24 Go RAM GRATUIT, et les ports mail sont ouverts.

> **Attention**: Hetzner, Contabo, DigitalOcean bloquent parfois le port 25 par defaut.
> Il faut demander le deblocage au support (delai: 1-24h).

---

## Pour votre team

### Partager l'acces

1. **Chaque membre** cree un compte sur le panel avec son propre login
2. Ou utilisez **RBAC** (Role-Based Access Control) integre pour gerer les permissions
3. Le panel est accessible via navigateur - pas besoin d'installer quoi que ce soit

### Workflow pour mises a jour

```bash
# Sur le serveur - pull + rebuild + restart
./deploy.sh update
```

### Workflow de dev (localement)

```bash
# Sur votre machine de dev
git pull origin main
# Modifier le code...
git add -A && git commit -m "feat: description"
git push origin main

# Sur le serveur
./deploy.sh update
```

---

## SSL / Certificats

BillionMail genere automatiquement des certificats SSL auto-signes au premier lancement.
Pour des certificats **Let's Encrypt** (gratuits, reconnus):

1. Assurez-vous que le port 80 est accessible
2. Le systeme ACME est integre (port 60880 en proxy)
3. Les certificats sont stockes dans `ssl/` et `ssl-self-signed/`

---

## Securite (rappel Phase 1)

Ce fork inclut les corrections de securite suivantes:
- Docker socket retire (plus d'acces root a l'hote)
- Capabilities NET_ADMIN/NET_RAW supprimees
- User non-root dans le container
- Mots de passes generes aleatoirement
- Fail2ban active par defaut

**APRES le premier login**:
1. Changer le mot de passe admin immediatement
2. Configurer l'IP whitelist si possible
3. Activer 2FA si disponible

---

## Troubleshooting

### Le build echoue
```bash
# Nettoyer et recommencer
docker compose -f docker-compose.custom.yml down
docker system prune -a
./deploy.sh build
```

### Port 25 bloque
```bash
# Tester si le port est ouvert
telnet smtp.google.com 25
# Si ca marche depuis le serveur, c'est que votre provider bloque
# Contacter le support pour debloquer
```

### Base de donnees ne demarre pas
```bash
# Verifier les permissions
ls -la postgresql-data/
# Si necessaire:
sudo chown -R 999:999 postgresql-data/
```

### Logs utiles
```bash
./deploy.sh logs core-billionmail    # Backend Go
./deploy.sh logs postfix-billionmail  # Serveur mail
./deploy.sh logs dovecot-billionmail  # IMAP/POP3
./deploy.sh logs rspamd-billionmail   # Anti-spam
```
