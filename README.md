1. **Génération de Room**
- Création d'une URL unique avec :
  - *Identifiant de room* : UUID v4 (ex: `/room/9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d`)
  - *Clé de chiffrement* : Dans le fragment URL (après `#`) invisible côté serveur

2. **Architecture Backend (Go/Resgate/NATS)**
- Routes :
  - `POST /room` → Génère une nouvelle room (UUID + clé aléatoire)
  - `WS /room/{id}/ws` → Endpoint WebSocket pour le flux de messages
- Communication :
  - Chaque room = sujet NATS (`room.{id}`)
  - Messages relayés en temps réel sans persistance

3. **Chiffrement End-to-End**
- Mécanisme client-side :
  - AES-GCM avec clé dérivée du fragment URL
  - Nonces générés côté client pour chaque message
- Le serveur ne manipule que des payloads binaires chiffrés

4. **Gestion des Connexions**
- Vie éphémère des rooms :
  - La room persiste tant qu'au moins 1 client connecté
  - Nettoyage automatique par Resgate après dernier disconnect
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

Points de Vigilance :
- Sécurité URL : Le fragment doit rester confidentiel (risque de fute par historique navigateur)
- Forward Secrecy : Optionnel, possible via ratchet client-side (Double Ratchet)
- Protection anti-DoS : Limiter les connexions/room par IP

Cette approche conserve la vie privée par design tout en utilisant les capacités temps-réel de NATS. Les données sensibles ne transitent jamais en clair et le serveur reste aveugle aux contenus.
