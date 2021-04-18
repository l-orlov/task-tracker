CREATE OR REPLACE FUNCTION trigger_set_timestamp()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE users
(
    id                 BIGSERIAL PRIMARY KEY,
    email              VARCHAR(320)  NOT NULL UNIQUE,
    first_name         VARCHAR(255)  NOT NULL,
    last_name          VARCHAR(255)  NOT NULL,
    password           VARCHAR(255)  NOT NULL,
    is_email_confirmed BOOLEAN       NOT NULL DEFAULT FALSE,
    created_at         TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_timestamp
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();
