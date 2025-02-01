BEGIN;

DROP INDEX IF EXISTS idx_messages_room;
DROP INDEX IF EXISTS idx_messages_timestamp;
DROP INDEX IF EXISTS idx_rooms_activity;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS rooms;
DROP EXTENSION IF EXISTS "uuid-ossp";

COMMIT;
