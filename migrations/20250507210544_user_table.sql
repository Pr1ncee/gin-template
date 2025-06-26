-- +goose Up
-- +goose StatementBegin
CREATE TYPE user_role AS ENUM ('User', 'Admin');
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "user" (
                                      "id" SERIAL PRIMARY KEY,
                                      "first_name" VARCHAR(255) NOT NULL,
    "last_name" VARCHAR(255) NOT NULL,
    "email" VARCHAR(255) NOT NULL UNIQUE,
    "password" BYTEA NOT NULL,
    "phone_number" VARCHAR(50),
    "role" user_role NOT NULL DEFAULT 'User',
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
                               );
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER update_user_updated_at
    BEFORE UPDATE ON "user"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER update_user_updated_at ON "user";

-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE "user";
DROP TYPE "user_role";
-- +goose StatementEnd
