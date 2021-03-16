BEGIN;

CREATE TABLE IF NOT EXISTS unique_identities
(
    email_address VARCHAR(255)
        CONSTRAINT unique_identities_pk
            PRIMARY KEY,
    identity_id VARCHAR(255) DEFAULT NULL NOT NULL
);

CREATE INDEX IF NOT EXISTS unique_identities_identity_id_idx
    on unique_identities (identity_id);

COMMIT;