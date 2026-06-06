-- 001_init.up.sql
-- MiFinca: Complete initial schema

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ─────────────────────────────────────────
-- TENANTS
-- ─────────────────────────────────────────
CREATE TABLE tenants (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name                VARCHAR(200) NOT NULL,
    plan                VARCHAR(50)  NOT NULL DEFAULT 'basic',
    global_sync_version BIGINT       NOT NULL DEFAULT 0,
    is_active           BOOLEAN      NOT NULL DEFAULT true,
    created_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- ─────────────────────────────────────────
-- USERS
-- ─────────────────────────────────────────
CREATE TABLE users (
    id            UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id     UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    full_name     VARCHAR(200) NOT NULL,
    role          VARCHAR(20)  NOT NULL CHECK (role IN ('owner','vet','worker')),
    is_active     BOOLEAN      NOT NULL DEFAULT true,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_users_tenant ON users(tenant_id);

-- ─────────────────────────────────────────
-- DEVICES
-- ─────────────────────────────────────────
CREATE TABLE devices (
    id                UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id         UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id           UUID        NOT NULL REFERENCES users(id),
    device_uid        VARCHAR(255) NOT NULL,
    last_sync_version BIGINT       NOT NULL DEFAULT 0,
    last_seen_at      TIMESTAMPTZ,
    created_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, device_uid)
);

-- ─────────────────────────────────────────
-- FARMS
-- ─────────────────────────────────────────
CREATE TABLE farms (
    id           UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id    UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name         VARCHAR(200) NOT NULL,
    address      TEXT,
    latitude     DECIMAL(10,8),
    longitude    DECIMAL(11,8),
    notes        TEXT,
    sync_version BIGINT      NOT NULL DEFAULT 0,
    is_deleted   BOOLEAN     NOT NULL DEFAULT false,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_farms_tenant ON farms(tenant_id);
CREATE INDEX idx_farms_sync   ON farms(tenant_id, sync_version);

-- ─────────────────────────────────────────
-- ANIMALS
-- ─────────────────────────────────────────
CREATE TABLE animals (
    id           UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id    UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    farm_id      UUID        NOT NULL REFERENCES farms(id),
    visual_tag   VARCHAR(50) NOT NULL,
    rfid_tag     VARCHAR(100),
    species      VARCHAR(20) NOT NULL CHECK (species IN ('bovino','ovino','caprino','porcino','equino')),
    breed        VARCHAR(100),
    sex          CHAR(1)     NOT NULL CHECK (sex IN ('M','F')),
    birth_date   DATE        NOT NULL,
    weight_kg    DECIMAL(7,2),
    color        VARCHAR(100),
    status       VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active','sold','dead','lost')),
    sire_id      UUID        REFERENCES animals(id),
    dam_id       UUID        REFERENCES animals(id),
    photo_url    TEXT,
    notes        TEXT,
    status_date  DATE,
    status_notes TEXT,
    sync_version BIGINT      NOT NULL DEFAULT 0,
    is_deleted   BOOLEAN     NOT NULL DEFAULT false,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, visual_tag),
    UNIQUE(rfid_tag)
);
CREATE INDEX idx_animals_tenant ON animals(tenant_id);
CREATE INDEX idx_animals_farm   ON animals(tenant_id, farm_id) WHERE NOT is_deleted;
CREATE INDEX idx_animals_sync   ON animals(tenant_id, sync_version);
CREATE INDEX idx_animals_status ON animals(tenant_id, status)  WHERE NOT is_deleted;
CREATE INDEX idx_animals_search ON animals(tenant_id, visual_tag, rfid_tag);

-- ─────────────────────────────────────────
-- HEALTH_EVENTS
-- ─────────────────────────────────────────
CREATE TABLE health_events (
    id            UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id     UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    animal_id     UUID        NOT NULL REFERENCES animals(id),
    event_type    VARCHAR(30) NOT NULL CHECK (event_type IN ('VACCINE','TREATMENT','DIAGNOSIS','DEWORMING','SURGERY','CHECKUP')),
    event_date    DATE        NOT NULL,
    product_name  VARCHAR(200),
    dose_ml       DECIMAL(7,2),
    lot_number    VARCHAR(100),
    vet_name      VARCHAR(200),
    next_date     DATE,
    notes         TEXT,
    cost          DECIMAL(10,2),
    bulk_batch_id UUID,
    sync_version  BIGINT      NOT NULL DEFAULT 0,
    is_deleted    BOOLEAN     NOT NULL DEFAULT false,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_health_animal ON health_events(tenant_id, animal_id) WHERE NOT is_deleted;
CREATE INDEX idx_health_sync   ON health_events(tenant_id, sync_version);
CREATE INDEX idx_health_next   ON health_events(tenant_id, next_date) WHERE next_date IS NOT NULL AND NOT is_deleted;
CREATE INDEX idx_health_batch  ON health_events(bulk_batch_id) WHERE bulk_batch_id IS NOT NULL;

-- ─────────────────────────────────────────
-- MOVEMENTS
-- ─────────────────────────────────────────
CREATE TABLE movements (
    id            UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id     UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    animal_id     UUID        NOT NULL REFERENCES animals(id),
    from_farm_id  UUID        REFERENCES farms(id),
    to_farm_id    UUID        REFERENCES farms(id),
    movement_date DATE        NOT NULL,
    reason        VARCHAR(30) NOT NULL CHECK (reason IN ('SALE_TRANSFER','GRAZING','VETERINARY','RETURN','PURCHASE','OTHER')),
    notes         TEXT,
    sync_version  BIGINT      NOT NULL DEFAULT 0,
    is_deleted    BOOLEAN     NOT NULL DEFAULT false,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_movements_animal ON movements(tenant_id, animal_id);
CREATE INDEX idx_movements_farm   ON movements(tenant_id, from_farm_id);
CREATE INDEX idx_movements_sync   ON movements(tenant_id, sync_version);

-- ─────────────────────────────────────────
-- NOTIFICATIONS
-- ─────────────────────────────────────────
CREATE TABLE notifications (
    id               UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id        UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    health_event_id  UUID        REFERENCES health_events(id),
    type             VARCHAR(50) NOT NULL,
    message          TEXT        NOT NULL,
    is_read          BOOLEAN     NOT NULL DEFAULT false,
    scheduled_date   DATE        NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(health_event_id, type)
);

-- ─────────────────────────────────────────
-- FUNCTION: increment sync_version atomically
-- ─────────────────────────────────────────
CREATE OR REPLACE FUNCTION increment_tenant_sync_version(p_tenant_id UUID)
RETURNS BIGINT AS $$
DECLARE v BIGINT;
BEGIN
    UPDATE tenants
    SET global_sync_version = global_sync_version + 1
    WHERE id = p_tenant_id
    RETURNING global_sync_version INTO v;
    RETURN v;
END;
$$ LANGUAGE plpgsql;

-- ─────────────────────────────────────────
-- TRIGGER: auto-advance sync_version on UPDATE
-- ─────────────────────────────────────────
CREATE OR REPLACE FUNCTION trg_sync_version() RETURNS TRIGGER AS $$
BEGIN
    NEW.sync_version := increment_tenant_sync_version(NEW.tenant_id);
    NEW.updated_at   := NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER animals_sync       BEFORE UPDATE ON animals       FOR EACH ROW EXECUTE FUNCTION trg_sync_version();
CREATE TRIGGER health_sync        BEFORE UPDATE ON health_events  FOR EACH ROW EXECUTE FUNCTION trg_sync_version();
CREATE TRIGGER movements_sync     BEFORE UPDATE ON movements      FOR EACH ROW EXECUTE FUNCTION trg_sync_version();
CREATE TRIGGER farms_sync         BEFORE UPDATE ON farms          FOR EACH ROW EXECUTE FUNCTION trg_sync_version();

-- ─────────────────────────────────────────
-- ROW LEVEL SECURITY (multi-tenancy)
-- ─────────────────────────────────────────
ALTER TABLE animals       ENABLE ROW LEVEL SECURITY;
ALTER TABLE health_events ENABLE ROW LEVEL SECURITY;
ALTER TABLE movements     ENABLE ROW LEVEL SECURITY;
ALTER TABLE farms         ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_iso ON animals       USING (tenant_id = current_setting('app.current_tenant_id', true)::uuid);
CREATE POLICY tenant_iso ON health_events USING (tenant_id = current_setting('app.current_tenant_id', true)::uuid);
CREATE POLICY tenant_iso ON movements     USING (tenant_id = current_setting('app.current_tenant_id', true)::uuid);
CREATE POLICY tenant_iso ON farms         USING (tenant_id = current_setting('app.current_tenant_id', true)::uuid);
