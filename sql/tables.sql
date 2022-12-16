BEGIN;

DROP TABLE IF EXISTS "device";
DROP TABLE IF EXISTS "configuration";
DROP TABLE IF EXISTS "workload";
DROP TABLE IF EXISTS "devices_workloads";
DROP TABLE IF EXISTS "device_set";
DROP TABLE IF EXISTS "devices_sets";
DROP TABLE IF EXISTS "hardware";
DROP TABLE IF EXISTS "network_interface";

CREATE TABLE configuration (
    id TEXT PRIMARY KEY,
    hardware_profile_scope TEXT DEFAULT 'full', -- full scope
    hardware_profile_include BOOLEAN DEFAULT true,
    heartbeat_period_seconds SMALLINT DEFAULT 30, -- 30s
    CHECK (heartbeat_period_seconds > 0)
);

-- hardware section
CREATE TABLE os_information (
    id TEXT PRIMARY KEY
);

CREATE TABLE system_vendor (
    id TEXT PRIMARY KEY
);

CREATE TABLE hardware (
    id TEXT PRIMARY KEY,
    os_information_id TEXT REFERENCES os_information ON DELETE CASCADE,
    system_vendor_id TEXT REFERENCES system_vendor ON DELETE CASCADE
);

CREATE TABLE network_interface (
    id TEXT PRIMARY KEY,
    hardware_id TEXT REFERENCES hardware,
    name TEXT NOT NULL,
    mac_address MACADDR8 NOT NULL,
    has_carrier BOOLEAN NOT NULL,
    ip4 INET[]
);

CREATE TABLE workload (
    id TEXT PRIMARY KEY,
    path TEXT NOT NULL
);

CREATE TABLE namespace (
    id TEXT PRIMARY KEY,
    is_default BOOLEAN DEFAULT false,
    configuration_id TEXT REFERENCES configuration(id) ON DELETE CASCADE
);

CREATE TABLE device_set (
    id TEXT PRIMARY KEY,
    configuration_id TEXT REFERENCES configuration(id) ON DELETE CASCADE,
    namespace_id TEXT REFERENCES namespace(id) ON DELETE CASCADE
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
    configuration_id TEXT REFERENCES configuration(id) ON DELETE SET NULL,
    hardware_id TEXT REFERENCES hardware(id) ON DELETE SET NULL
);

CREATE INDEX device_configuration_id_idx ON device (configuration_id);

CREATE TABLE devices_workloads (
    device_id TEXT REFERENCES device ON DELETE CASCADE,
    workload_id TEXT REFERENCES workload ON DELETE CASCADE,
    CONSTRAINT devices_workloads_pk PRIMARY KEY (
        device_id,
        workload_id
    )
);

CREATE TABLE namespaces_workloads (
    namespace_id TEXT REFERENCES namespace(id) ON DELETE CASCADE,
    workload_id TEXT REFERENCES workload(id) ON DELETE CASCADE,
    CONSTRAINT namespace_workload_pk PRIMARY KEY(
        namespace_id,
        workload_id
    )
);

CREATE TABLE sets_workloads (
    device_set_id TEXT REFERENCES device_set(id) ON DELETE CASCADE,
    workload_id TEXT REFERENCES workload(id) ON DELETE CASCADE,
    CONSTRAINT device_set_workload_pk PRIMARY KEY(
        device_set_id,
        workload_id
    )
);

COMMIT;
