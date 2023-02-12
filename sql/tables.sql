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

CREATE TYPE ref_type as ENUM ('manifest', 'configuration');

CREATE TABLE reference (
    id varchar(255) PRIMARY KEY,
    ref_type ref_type NOT NULL,
    name varchar(255) NOT NULL,
    repo_id varchar(255) NOT NULL REFERENCES repo(id) ON DELETE CASCADE,
    valid BOOLEAN NOT NULL,
    hash TEXT NOT NULL,
    path_reference TEXT NOT NULL
);

CREATE TABLE secret (
    id varchar(255) PRIMARY KEY,
    path varchar(255) NOT NULL,
    current_hash TEXT NOT NULL,
    target_hash TEXT NOT NULL
);

CREATE TABLE secrets_manifests (
    secret_id varchar(255) REFERENCES secret(id),
    reference_id varchar(255) REFERENCES reference(id),
    CONSTRAINT secret_reference_pk PRIMARY KEY(
        secret_id,
        reference_id
    )
);

CREATE TABLE namespace (
    id TEXT PRIMARY KEY,
    is_default BOOLEAN DEFAULT false,
    reference_id varchar(255) REFERENCES reference(id) ON DELETE SET NULL
);

CREATE TABLE device_set (
    id TEXT PRIMARY KEY,
    reference_id varchar(255) REFERENCES reference(id) ON DELETE SET NULL,
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
    reference_id varchar(255) REFERENCES reference(id) ON DELETE SET NULL
);

CREATE TABLE devices_references (
    device_id varchar(255) REFERENCES device(id) ON DELETE CASCADE,
    reference_id varchar(255) REFERENCES reference(id) ON DELETE CASCADE,
    CONSTRAINT devices_references_pk PRIMARY KEY (
        device_id,
        reference_id
    )
);

CREATE TABLE namespaces_references (
    namespace_id varchar(255) REFERENCES namespace(id) ON DELETE CASCADE,
    reference_id varchar(255) REFERENCES reference(id) ON DELETE CASCADE,
    CONSTRAINT namespace_reference_pk PRIMARY KEY(
        namespace_id,
        reference_id
    )
);

CREATE TABLE sets_references (
    device_set_id varchar(255) REFERENCES device_set(id) ON DELETE CASCADE,
    reference_id varchar(255) REFERENCES reference(id) ON DELETE CASCADE,
    CONSTRAINT device_set_reference_pk PRIMARY KEY(
        device_set_id,
        reference_id
    )
);

COMMIT;
