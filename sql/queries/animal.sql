-- name: CreateAnimal :one
INSERT INTO animals (
    id, tenant_id, farm_id, visual_tag, rfid_tag, species, breed, sex,
    birth_date, weight_kg, color, sire_id, dam_id, notes, sync_version
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, increment_tenant_sync_version($2)
) RETURNING *;

-- name: GetAnimalByID :one
SELECT * FROM animals
WHERE tenant_id = $1 AND id = $2 AND NOT is_deleted;

-- name: GetAnimalByIDForUpdate :one
SELECT * FROM animals
WHERE tenant_id = $1 AND id = $2
FOR UPDATE;

-- name: SoftDeleteAnimal :exec
UPDATE animals
SET is_deleted = true,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2;

-- name: ListAnimals :many
SELECT * FROM animals
WHERE tenant_id = $1
  AND NOT is_deleted
  AND (sqlc.narg('farm_id')::uuid IS NULL OR farm_id = sqlc.narg('farm_id'))
  AND (sqlc.narg('species')::varchar IS NULL OR species = sqlc.narg('species'))
  AND (sqlc.narg('sex')::char IS NULL OR sex = sqlc.narg('sex'))
  AND (sqlc.narg('status')::varchar IS NULL OR status = sqlc.narg('status'))
  AND (sqlc.narg('search')::varchar IS NULL OR visual_tag ILIKE '%' || sqlc.narg('search') || '%' OR rfid_tag ILIKE '%' || sqlc.narg('search') || '%')
  AND (sqlc.narg('cursor')::varchar IS NULL OR visual_tag > sqlc.narg('cursor'))
ORDER BY visual_tag ASC
LIMIT $2;

-- name: GetAnimalsSince :many
SELECT * FROM animals
WHERE tenant_id = $1 AND sync_version > $2
ORDER BY sync_version ASC;

-- name: CountAnimalsByTenant :one
SELECT COUNT(*) FROM animals
WHERE tenant_id = $1 AND NOT is_deleted;

-- name: GetPedigreeRows :many
WITH RECURSIVE pedigree AS (
    SELECT a.id, a.visual_tag, a.species, a.breed, a.sex, a.birth_date, a.sire_id, a.dam_id, 1 AS depth
    FROM animals a WHERE a.id = $1 AND a.tenant_id = $2
    UNION ALL
    SELECT a.id, a.visual_tag, a.species, a.breed, a.sex, a.birth_date, a.sire_id, a.dam_id, p.depth + 1
    FROM animals a
    INNER JOIN pedigree p ON a.id = p.sire_id OR a.id = p.dam_id
    WHERE p.depth < 4
)
SELECT p.id, p.visual_tag, p.species, p.breed, p.sex, p.birth_date, p.sire_id, p.dam_id, p.depth FROM pedigree p;
