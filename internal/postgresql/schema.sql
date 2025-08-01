CREATE TYPE user_role AS ENUM ('User', 'Admin');

CREATE TABLE "user" (
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
