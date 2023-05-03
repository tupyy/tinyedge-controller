BEGIN;

CREATE TABLE repo (
    id varchar(255) PRIMARY KEY,
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

CREATE TYPE ref_type as ENUM ('workload', 'configuration');

CREATE TABLE manifest (
    id varchar(255) PRIMARY KEY,
    ref_type ref_type NOT NULL,
    name varchar(255) NOT NULL,
    repo_id varchar(255) NOT NULL REFERENCES repo(id) ON DELETE CASCADE,
    path TEXT NOT NULL
);

CREATE TABLE namespace (
    id TEXT PRIMARY KEY,
    is_default BOOLEAN DEFAULT false,
    configuration_manifest_id varchar(255) REFERENCES manifest(id) ON DELETE CASCADE
);

CREATE TABLE device_set (
    id TEXT PRIMARY KEY,
    configuration_manifest_id varchar(255) REFERENCES manifest(id) ON DELETE SET NULL,
    namespace_id varchar(255) NOT NULL REFERENCES namespace(id) ON DELETE CASCADE
);

CREATE TABLE device (
    id TEXT PRIMARY KEY,
    enroled_at TIMESTAMP,
    registered_at TIMESTAMP,
    enroled TEXT NOT NULL DEFAULT 'not_enroled',
    registered BOOLEAN NOT NULL DEFAULT false,
    certificate_sn TEXT,
    namespace_id varchar(255) NOT NULL REFERENCES namespace(id) ON DELETE SET NULL,
    device_set_id varchar(255) REFERENCES device_set(id) ON DELETE SET NULL,
    configuration_manifest_id varchar(255) REFERENCES manifest(id) ON DELETE SET NULL
);

CREATE TABLE devices_manifests (
    device_id varchar(255) REFERENCES device(id) ON DELETE CASCADE,
    manifest_id varchar(255) REFERENCES manifest(id) ON DELETE CASCADE,
    CONSTRAINT devices_manifests_pk PRIMARY KEY (
        device_id,
        manifest_id
    )
);

CREATE TABLE namespaces_manifests (
    namespace_id varchar(255) REFERENCES namespace(id) ON DELETE CASCADE,
    manifest_id varchar(255) REFERENCES manifest(id) ON DELETE CASCADE,
    CONSTRAINT namespace_manifest_pk PRIMARY KEY(
        namespace_id,
        manifest_id
    )
);

CREATE TABLE sets_manifests (
    device_set_id varchar(255) REFERENCES device_set(id) ON DELETE CASCADE,
    manifest_id varchar(255) REFERENCES manifest(id) ON DELETE CASCADE,
    CONSTRAINT device_set_manifest_pk PRIMARY KEY(
        device_set_id,
        manifest_id
    )
);

CREATE TABLE secret (
    id varchar(255) PRIMARY KEY,
    path varchar(255) NOT NULL,
    current_hash TEXT NOT NULL,
    target_hash TEXT NOT NULL
);

CREATE TABLE secrets_manifests (
    secret_id varchar(255) REFERENCES secret(id),
    manifest_id varchar(255) REFERENCES manifest(id),
    CONSTRAINT secret_manifest_pk PRIMARY KEY(
        secret_id,
        manifest_id
    )
);

COMMIT;
