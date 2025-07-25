-- name: CreateAnalyticsEvent :one
INSERT INTO analytics_events (
    event_type, page_path, user_agent, ip_address, session_id, referrer, metadata
) VALUES (
    ?, ?, ?, ?, ?, ?, ?
) RETURNING *;

-- name: GetAnalyticsEvent :one
SELECT * FROM analytics_events WHERE id = ?;

-- name: ListAnalyticsEvents :many
SELECT * FROM analytics_events 
ORDER BY created_at DESC 
LIMIT ? OFFSET ?;

-- name: ListAnalyticsEventsByType :many
SELECT * FROM analytics_events 
WHERE event_type = ?
ORDER BY created_at DESC 
LIMIT ? OFFSET ?;

-- name: GetPageViewStats :many
SELECT 
    page_path,
    COUNT(*) as view_count,
    COUNT(DISTINCT session_id) as unique_visitors,
    DATE(created_at) as date
FROM analytics_events 
WHERE event_type = 'page_view'
    AND created_at >= datetime('now', '-30 days')
GROUP BY page_path, DATE(created_at)
ORDER BY date DESC, view_count DESC;

-- name: GetEventCountsByType :many
SELECT 
    event_type,
    COUNT(*) as event_count,
    DATE(created_at) as date
FROM analytics_events 
WHERE created_at >= datetime('now', '-30 days')
GROUP BY event_type, DATE(created_at)
ORDER BY date DESC, event_count DESC;

-- name: DeleteOldAnalyticsEvents :exec
DELETE FROM analytics_events 
WHERE created_at < datetime('now', '-90 days');