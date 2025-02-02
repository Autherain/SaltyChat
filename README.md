Je vais restructurer le README pour intégrer les précisions et le nouveau schéma SQL. Voici la proposition révisée :

---

# Secure Real-time Chat - Architecture Technique

## 1. Génération de Room
- **URL unique** composée de :
  - `UUID v4` comme identifiant (ex: `/room/9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d`)
  - **Clé de chiffrement** stockée dans le fragment URL (partie après `#`), invisible côté serveur
- **Accès éphémère** : Le lien devient inutilisable après fermeture de tous les onglets clients

## 2. Chiffrement End-to-End
- **Mécanisme client-side** :
  - AES-GCM avec clé dérivée du fragment URL
  - Nonces uniques (12 bytes) générés à chaque message
- **Zero-knowledge serveur** :
  - Le serveur manipule exclusivement des blobs binaires (`encrypted_content` + `nonce`)
  - Aucune exposition des métadonnées utilisateur

## 3. Gestion des Connexions
- **Cycle de vie des rooms** :
  - Activée au premier client connecté via Resgate
  - Désactivée automatiquement après dernier `WebSocket disconnect`
  - Données historiques conservées en base (hors clé de chiffrement)
- **Sessions utilisateurs** :
  - Identifiées par la connexion WebSocket active
  - Détection instantanée des déconnexions (participants visibles en temps réel)

## 4. Flux de Données
1. Client A → Chiffre message + génère nonce
2. Publication sur la ressource NATS `room.{id}.messages`
3. Resgate diffuse à tous les subscribers de la room
4. Clients B/C → Déchiffrement via la clé locale

## 5. Frontend Minimaliste
- **Fonctionnalités clés** :
  - Lecture dynamique du fragment URL pour l'initialisation
  - Abonnement automatique à `room.{id}.messages` via Resgate
  - Purge mémoire des clés sur événement `beforeunload`
  - UI reactive avec masquage des messages après déconnexion

## 6. Schéma de Base de Données (PostgreSQL)
```sql
-- Table des salles (metadata seulement)
CREATE TABLE rooms (
    id UUID PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    last_activity TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP, -- Dernière interaction
    is_active BOOLEAN DEFAULT TRUE -- État géré par Resgate
);

-- Table des messages chiffrés
CREATE TABLE messages (
    id UUID PRIMARY KEY,
    room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    encrypted_content BYTEA NOT NULL, -- Payload chiffré
    nonce BYTEA NOT NULL CHECK (octet_length(nonce) = 12), -- 96 bits pour AES-GCM
    timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Index d'optimisation
CREATE INDEX idx_messages_room ON messages(room_id);
CREATE INDEX idx_messages_timestamp ON messages(timestamp);
```

**Explications techniques** :
- **Sécurité renforcée** :
  - Les colonnes sensibles utilisent des types binaires natifs (`BYTEA`)
  - Contraintes de taille strictes pour les nonces
  - Cascade de suppression pour l'archivage automatique
- **Performances** :
  - Indexation ciblée sur les requêtes temporelles et par room
  - Séparation metadata/payload pour l'optimisation stockage
- **Audit** :
  - Horodatage UTC avec précision microseconde (`TIMESTAMPTZ`)
  - Trace d'activité via `last_activity`

## 7. Architecture Serveur
- **Resgate** : Gère les subscriptions temps-réel et le cycle de vie des rooms
- **NATS** : Bus de messages pour la diffusion globale
- **Base de données** : Stockage persistant des messages chiffrés (hors clé)

Cette architecture garantit :
- 🔒 **Confidentialité** par chiffrement client-to-client
- ⚡ **Réactivité** grâce au stack NATS/Resgate
- 🧹 **Auto-nettoyage** des ressources inactives
