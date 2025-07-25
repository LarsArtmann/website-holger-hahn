-- name: CreateContact :one
INSERT INTO contacts (
    name, email, company, message, subject, source
) VALUES (
    ?, ?, ?, ?, ?, ?
) RETURNING *;

-- name: GetContact :one
SELECT * FROM contacts WHERE id = ?;

-- name: GetContactByEmail :one
SELECT * FROM contacts WHERE email = ? ORDER BY created_at DESC LIMIT 1;

-- name: ListContacts :many
SELECT * FROM contacts 
ORDER BY created_at DESC 
LIMIT ? OFFSET ?;

-- name: ListContactsByStatus :many
SELECT * FROM contacts 
WHERE status = ?
ORDER BY created_at DESC 
LIMIT ? OFFSET ?;

-- name: UpdateContactStatus :one
UPDATE contacts 
SET status = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeleteContact :exec
DELETE FROM contacts WHERE id = ?;

-- name: CountContacts :one
SELECT COUNT(*) FROM contacts;

-- name: CountContactsByStatus :one
SELECT COUNT(*) FROM contacts WHERE status = ?;