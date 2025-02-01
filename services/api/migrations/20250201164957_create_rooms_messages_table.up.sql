BEGIN; 

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Table des salles de chat
CREATE TABLE rooms (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    last_activity TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE
);

-- Table des messages
CREATE TABLE messages (
    id BIGSERIAL PRIMARY KEY,
    room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    encrypted_content BYTEA NOT NULL,
    nonce BYTEA NOT NULL CHECK (octet_length(nonce) = 12),
    timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Index pour optimiser les requêtes fréquentes
CREATE INDEX idx_messages_room ON messages(room_id);
CREATE INDEX idx_messages_timestamp ON messages(timestamp);
CREATE INDEX idx_rooms_activity ON rooms(last_activity) WHERE is_active = TRUE;

COMMIT;
