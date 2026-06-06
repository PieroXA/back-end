-- 001_init.down.sql
-- Rollback the entire initial schema

DROP TRIGGER IF EXISTS farms_sync     ON farms;
DROP TRIGGER IF EXISTS movements_sync ON movements;
DROP TRIGGER IF EXISTS health_sync    ON health_events;
DROP TRIGGER IF EXISTS animals_sync   ON animals;

DROP FUNCTION IF EXISTS trg_sync_version();
DROP FUNCTION IF EXISTS increment_tenant_sync_version(UUID);

DROP POLICY IF EXISTS tenant_iso ON farms;
DROP POLICY IF EXISTS tenant_iso ON movements;
DROP POLICY IF EXISTS tenant_iso ON health_events;
DROP POLICY IF EXISTS tenant_iso ON animals;

DROP TABLE IF EXISTS notifications CASCADE;
DROP TABLE IF EXISTS movements CASCADE;
DROP TABLE IF EXISTS health_events CASCADE;
DROP TABLE IF EXISTS animals CASCADE;
DROP TABLE IF EXISTS farms CASCADE;
DROP TABLE IF EXISTS devices CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS tenants CASCADE;
