-- name: CreateMovement :one
INSERT INTO movements (
    id, tenant_id, animal_id, from_farm_id, to_farm_id, movement_date, reason, notes, sync_version
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, increment_tenant_sync_version($2)
) RETURNING *;

-- name: ListMovementsByAnimal :many
SELECT * FROM movements
WHERE tenant_id = $1 AND animal_id = $2 AND NOT is_deleted
ORDER BY movement_date DESC;

-- name: GetMovementsSince :many
SELECT * FROM movements
WHERE tenant_id = $1 AND sync_version > $2
ORDER BY sync_version ASC;
