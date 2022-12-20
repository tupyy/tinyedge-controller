BEGIN;

INSERT INTO configuration (id,heartbeat_period_seconds) VALUES ('default', 10);
INSERT INTO namespace (id, is_default, configuration_id) VALUES ('default', true, 'default');
INSERT INTO device_set(id,namespace_id, configuration_id) VALUES ('default', 'default', 'default');

COMMIT;
