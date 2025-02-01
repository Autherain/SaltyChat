1. **Génération de Room**
- Création d'une URL unique avec :
  - *Identifiant de room* : UUID v4 (ex: `/room/9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d`)
  - *Clé de chiffrement* : Dans le fragment URL (après `#`) invisible côté serveur

3. **Chiffrement End-to-End**
- Mécanisme client-side :
  - AES-GCM avec clé dérivée du fragment URL
  - Nonces générés côté client pour chaque message
- Le serveur ne manipule que des payloads binaires chiffrés

4. **Gestion des Connexions**
- Vie éphémère des rooms :
  - La room persiste tant qu'au moins 1 client connecté
  - Nettoyage automatique par Resgate après dernier disconnect
- La personne ne peut pas ré-utiliser le lien qu'on lui a envoyé pour accéder au chat et discuter. Cela ne dure que temps l'onglet n'est pas fermé. 
- Session utilisateur :
  - Identifiée par le WebSocket ouvert
  - Fermeture = suppression immédiate de la liste des participants

5. **Flux de Données**
- Client A envoie message → Chiffrement → Publie sur `room.{id}.messages`
- NATS diffuse à tous les subscribers → Déchiffrement côté clients B, C...

6. **Frontend Minimaliste**
- Éléments clés :
  - Lecture dynamique du fragment URL pour la clé
  - Connexion WebSocket auto-init à l'ouverture
  - Destruction des clés en mémoire à l'`onbeforeunload`

Cette approche conserve la vie privée par design tout en utilisant les capacités temps-réel de NATS. Les données sensibles ne transitent jamais en clair et le serveur reste aveugle aux contenus.


```sql
PRAGMA foreign_keys = ON;

CREATE TABLE rooms (
    id TEXT PRIMARY KEY,          -- UUID v4
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    room_id TEXT NOT NULL,
    encrypted_content BLOB NOT NULL,  -- Données chiffrées (AES-GCM)
    nonce BLOB NOT NULL,              -- Valeur aléatoire (12 bytes pour AES-GCM)
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE
);

CREATE INDEX idx_messages_room ON messages(room_id);
CREATE INDEX idx_messages_timestamp ON messages(timestamp);
```

Explications :

1. **Table `rooms`** :
- Stocke l'UUID unique de chaque room
- `created_at` pour audit (optionnel)
- Suppression automatique des messages liés via le `ON DELETE CASCADE`

2. **Table `messages`** :
- `encrypted_content` : Message chiffré (format binaire)
- `nonce` : Vecteur d'initialisation pour AES-GCM
- Les métadonnées utilisateur (username) sont incluses dans le payload chiffré
- Indexation sur `room_id` pour les requêtes par salle

3. **Sécurité** :
- Aucune donnée sensible en clair
- La clé de chiffrement reste dans le fragment URL (jamais stockée)
- Les nonces sont uniques par message

4. **Gestion de la durée de vie** :
- Les rooms/messages sont automatiquement nettoyés par SQLite via les contraintes de clé étrangère
- La suppression d'une room entraîne la suppression de tous ses messages
