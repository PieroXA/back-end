-- name: CreateHealthEvent :one
INSERT INTO health_events (
    id, tenant_id, animal_id, event_type, event_date, product_name, dose_ml, lot_number, vet_name, next_date, notes, cost, bulk_batch_id, sync_version
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, increment_tenant_sync_version($2)
) RETURNING *;

-- name: ListHealthEventsByAnimal :many
SELECT * FROM health_events
WHERE tenant_id = $1 AND animal_id = $2 AND NOT is_deleted
ORDER BY event_date DESC;

-- name: GetHealthEventsSince :many
SELECT * FROM health_events
WHERE tenant_id = $1 AND sync_version > $2
ORDER BY sync_version ASC;
