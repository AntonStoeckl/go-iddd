BEGIN;

CREATE TABLE IF NOT EXISTS unique_email_addresses
(
    email_address VARCHAR(255)
        CONSTRAINT unique_email_addresses_pk
            PRIMARY KEY
);

COMMIT;