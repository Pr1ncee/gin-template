-- =============================================================================
-- USER TABLE CRUD OPERATIONS
-- =============================================================================

-- name: CreateUser :one
INSERT INTO "user" (
    "first_name",
    "last_name",
    "email",
    "password",
    "phone_number",
    "role"
) VALUES (
             $1, $2, $3, $4, $5, $6
         ) RETURNING *;

-- name: ListUsers :many
SELECT
    "id",
    "first_name",
    "last_name",
    "email",
    "phone_number",
    "role",
    "created_at",
    "updated_at"
FROM "user"
WHERE
  -- Role filter (optional - pass NULL to ignore)
    ($1::text IS NULL OR $1 = '' OR "role"::text = $1)
  AND
  -- Search filter (optional - pass NULL or empty string to ignore)
    ($2::text IS NULL OR $2 = '' OR
     "first_name" ILIKE '%' || $2 || '%' OR
     "last_name" ILIKE '%' || $2 || '%' OR
     "email" ILIKE '%' || $2 || '%')
ORDER BY
    CASE
        WHEN $2::text IS NOT NULL AND $2 != '' THEN "last_name"
        ELSE NULL
END DESC,
    CASE
        WHEN $2::text IS NOT NULL AND $2 != '' THEN "first_name"
        ELSE NULL
END DESC,
    CASE
        WHEN $2::text IS NULL OR $2 = '' THEN "created_at"
        ELSE NULL
END DESC
LIMIT $3 OFFSET $4;

-- name: GetUserByID :one
SELECT
    "id",
    "first_name",
    "last_name",
    "email",
    "phone_number",
    "role",
    "created_at",
    "updated_at"
FROM "user"
WHERE "id" = $1;

-- name: GetUserByEmail :one
SELECT
    "id",
    "first_name",
    "last_name",
    "email",
    "phone_number",
    "role",
    "created_at",
    "updated_at"
FROM "user"
WHERE "email" = $1;

-- name: GetUserWithPasswordById :one
SELECT
    "id",
    "first_name",
    "last_name",
    "email",
    "password",
    "phone_number",
    "role",
    "created_at",
    "updated_at"
FROM "user"
WHERE "id" = $1;

-- name: UpdateUserPartial :one
UPDATE "user"
SET
    "first_name" = COALESCE(sqlc.narg('first_name'), "first_name"),
    "last_name" = COALESCE(sqlc.narg('last_name'), "last_name"),
    "email" = COALESCE(sqlc.narg('email'), "email"),
    "phone_number" = COALESCE(sqlc.narg('phone_number'), "phone_number"),
    "password" = COALESCE(sqlc.narg('password'), password),
    "role" = COALESCE(sqlc.narg('role'), "role"),
    "updated_at" = CURRENT_TIMESTAMP
WHERE "id" = sqlc.arg('id')
    RETURNING *;

-- name: DeleteUserByID :exec
DELETE FROM "user"
WHERE "id" = $1;

-- name: DeleteUserByEmail :exec
DELETE FROM "user"
WHERE "email" = $1;

-- name: CheckUserExists :one
SELECT EXISTS(
    SELECT 1 FROM "user" WHERE "id" = $1
) as "user_exists";
