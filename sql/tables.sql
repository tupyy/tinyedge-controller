BEGIN;

-- table configuration holds the configuration for the device's agent
CREATE TABLE configuration (
    id TEXT PRIMARY KEY,
    heartbeat_period_seconds SMALLINT DEFAULT 30, -- 30s
    log_level TEXT DEFAULT 'info',
    CHECK (heartbeat_period_seconds > 0)
);

CREATE TABLE repo (
    id TEXT PRIMARY KEY,
    url TEXT NOT NULL,
    branch TEXT, -- should be an enum allowing only "master" or "main"
    local_path TEXT,
    auth_type varchar(20),
    auth_secret_path varchar(20),
    current_head_sha TEXT,
    target_head_sha TEXT,
    pull_period_seconds SMALLINT DEFAULT 20,
    CHECK(pull_period_seconds >= 0) -- if 0 stop pulling
);

CREATE TABLE manifest_reference (
    id TEXT PRIMARY KEY,
    repo_id TEXT NOT NULL REFERENCES repo(id) ON DELETE CASCADE,
    valid BOOLEAN NOT NULL,
    hash TEXT NOT NULL,
    path_manifest_reference TEXT NOT NULL
);

CREATE TABLE secret (
    id TEXT PRIMARY KEY,
    path TEXT NOT NULL,
    current_hash TEXT NOT NULL,
    target_hash TEXT NOT NULL
);

CREATE TABLE secrets_manifests (
    secret_id TEXT REFERENCES secret(id),
    manifest_reference_id TEXT REFERENCES manifest_reference(id),
    CONSTRAINT secret_manifest_reference_pk PRIMARY KEY(
        secret_id,
        manifest_reference_id
    )
);

CREATE TABLE namespace (
    id TEXT PRIMARY KEY,
    is_default BOOLEAN DEFAULT false,
    configuration_id TEXT NOT NULL REFERENCES configuration(id) ON DELETE SET NULL
);

CREATE TABLE device_set (
    id TEXT PRIMARY KEY,
    configuration_id TEXT REFERENCES configuration(id) ON DELETE SET NULL,
    namespace_id TEXT NOT NULL REFERENCES namespace(id) ON DELETE CASCADE
);

CREATE TABLE device (
    id TEXT PRIMARY KEY,
    enroled_at TIMESTAMP,
    registered_at TIMESTAMP,
    enroled TEXT NOT NULL DEFAULT 'not_enroled',
    registered BOOLEAN NOT NULL DEFAULT false,
    certificate_sn TEXT,
    namespace_id TEXT REFERENCES namespace(id) ON DELETE SET NULL,
    device_set_id TEXT REFERENCES device_set(id) ON DELETE SET NULL,
    configuration_id TEXT REFERENCES configuration(id) ON DELETE SET NULL
);

CREATE INDEX device_configuration_id_idx ON device (configuration_id);

CREATE TABLE devices_references (
    device_id TEXT REFERENCES device(id) ON DELETE CASCADE,
    manifest_reference_id TEXT REFERENCES manifest_reference(id) ON DELETE CASCADE,
    CONSTRAINT devices_references_pk PRIMARY KEY (
        device_id,
        manifest_reference_id
    )
);

CREATE TABLE namespaces_references (
    namespace_id TEXT REFERENCES namespace(id) ON DELETE CASCADE,
    manifest_reference_id TEXT REFERENCES manifest_reference(id) ON DELETE CASCADE,
    CONSTRAINT namespace_reference_pk PRIMARY KEY(
        namespace_id,
        manifest_reference_id
    )
);

CREATE TABLE sets_references (
    device_set_id TEXT REFERENCES device_set(id) ON DELETE CASCADE,
    manifest_reference_id TEXT REFERENCES manifest_reference(id) ON DELETE CASCADE,
    CONSTRAINT device_set_reference_pk PRIMARY KEY(
        device_set_id,
        manifest_reference_id
    )
);

COMMIT;
