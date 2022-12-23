BEGIN;

DROP TABLE IF EXISTS "device";
DROP TABLE IF EXISTS "configuration";
DROP TABLE IF EXISTS "workload";
DROP TABLE IF EXISTS "devices_workloads";
DROP TABLE IF EXISTS "device_set";
DROP TABLE IF EXISTS "devices_sets";
DROP TABLE IF EXISTS "hardware";
DROP TABLE IF EXISTS "network_interface";

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
    current_head_sha TEXT,
    target_head_sha TEXT,
    pull_period_seconds SMALLINT DEFAULT 20,
    CHECK(pull_period_seconds >= 0) -- if 0 stop pulling
);

CREATE TABLE manifest_work (
    id TEXT PRIMARY KEY,
    repo_id TEXT REFERENCES repo(id) NOT NULL,
    path_manifest_work TEXT NOT NULL,
    content TEXT NOT NULL
);

CREATE TABLE namespace (
    id TEXT PRIMARY KEY,
    is_default BOOLEAN DEFAULT false,
    configuration_id TEXT REFERENCES configuration(id) ON DELETE CASCADE
);

CREATE TABLE device_set (
    id TEXT PRIMARY KEY,
    configuration_id TEXT REFERENCES configuration(id) ON DELETE CASCADE,
    namespace_id TEXT NOT NULL REFERENCES namespace(id) ON DELETE CASCADE
);

CREATE TABLE device (
    id TEXT PRIMARY KEY,
    enroled_at TIMESTAMP,
    registered_at TIMESTAMP,
    enroled TEXT NOT NULL DEFAULT 'not_enroled',
    registered BOOLEAN NOT NULL DEFAULT false,
    certificate_sn TEXT,
    namespace_id TEXT REFERENCES namespace(id) ON DELETE CASCADE,
    device_set_id TEXT REFERENCES device_set(id) ON DELETE CASCADE,
    configuration_id TEXT REFERENCES configuration(id) ON DELETE SET NULL
);

CREATE INDEX device_configuration_id_idx ON device (configuration_id);

CREATE TABLE devices_workloads (
    device_id TEXT REFERENCES device(id) ON DELETE CASCADE,
    manifest_work_id TEXT REFERENCES manifest_work(id) ON DELETE CASCADE,
    CONSTRAINT devices_workloads_pk PRIMARY KEY (
        device_id,
        manifest_work_id
    )
);

CREATE TABLE namespaces_workloads (
    namespace_id TEXT REFERENCES namespace(id) ON DELETE CASCADE,
    manifest_work_id TEXT REFERENCES manifest_work(id) ON DELETE CASCADE,
    CONSTRAINT namespace_manifest_work_pk PRIMARY KEY(
        namespace_id,
        manifest_work_id
    )
);

CREATE TABLE sets_workloads (
    device_set_id TEXT REFERENCES device_set(id) ON DELETE CASCADE,
    manifest_work_id TEXT REFERENCES manifest_work(id) ON DELETE CASCADE,
    CONSTRAINT device_set_manifest_work_pk PRIMARY KEY(
        device_set_id,
        manifest_work_id
    )
);

CREATE TABLE configuration_cache (
    device_id TEXT PRIMARY KEY,
    configuration BYTEA
);

COMMIT;
