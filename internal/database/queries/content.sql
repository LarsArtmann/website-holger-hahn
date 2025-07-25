-- Technologies queries
-- name: ListTechnologies :many
SELECT * FROM technologies 
WHERE is_active = TRUE
ORDER BY category, sort_order, name;

-- name: ListTechnologiesByCategory :many
SELECT * FROM technologies 
WHERE category = ? AND is_active = TRUE
ORDER BY sort_order, name;

-- name: GetTechnology :one
SELECT * FROM technologies WHERE id = ?;

-- name: CreateTechnology :one
INSERT INTO technologies (
    name, category, icon_class, color_scheme, sort_order
) VALUES (
    ?, ?, ?, ?, ?
) RETURNING *;

-- name: UpdateTechnology :one
UPDATE technologies 
SET name = ?, category = ?, icon_class = ?, color_scheme = ?, 
    sort_order = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeleteTechnology :exec
DELETE FROM technologies WHERE id = ?;

-- Services queries
-- name: ListServices :many
SELECT * FROM services 
WHERE is_active = TRUE
ORDER BY sort_order, title;

-- name: GetService :one
SELECT * FROM services WHERE id = ?;

-- name: CreateService :one
INSERT INTO services (
    title, description, features, icon_svg, color_scheme, sort_order
) VALUES (
    ?, ?, ?, ?, ?, ?
) RETURNING *;

-- name: UpdateService :one
UPDATE services 
SET title = ?, description = ?, features = ?, icon_svg = ?, 
    color_scheme = ?, sort_order = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeleteService :exec
DELETE FROM services WHERE id = ?;

-- Experience queries
-- name: ListExperiences :many
SELECT * FROM experiences 
WHERE is_active = TRUE
ORDER BY is_current DESC, start_date DESC, sort_order;

-- name: GetExperience :one
SELECT * FROM experiences WHERE id = ?;

-- name: GetCurrentExperience :one
SELECT * FROM experiences 
WHERE is_current = TRUE AND is_active = TRUE
LIMIT 1;

-- name: CreateExperience :one
INSERT INTO experiences (
    company, position, description, achievements, technologies, 
    start_date, end_date, is_current, sort_order
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?
) RETURNING *;

-- name: UpdateExperience :one
UPDATE experiences 
SET company = ?, position = ?, description = ?, achievements = ?, 
    technologies = ?, start_date = ?, end_date = ?, is_current = ?, 
    sort_order = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeleteExperience :exec
DELETE FROM experiences WHERE id = ?;