BEGIN;

CREATE TABLE IF NOT EXISTS unique_email_addresses
(
    email_address VARCHAR(255)
        CONSTRAINT unique_email_addresses_pk
            PRIMARY KEY,
    customer_id VARCHAR(255) DEFAULT NULL NOT NULL
);

CREATE INDEX IF NOT EXISTS email_addresses_customer_id_idx
    on unique_email_addresses (customer_id);

COMMIT;