BEGIN;

INSERT INTO configuration (id,heartbeat_period_seconds, log_level) VALUES ('default', 10, 'debug');
INSERT INTO namespace (id, is_default, configuration_id) VALUES ('default', true, 'default');
INSERT INTO device_set(id,namespace_id, configuration_id) VALUES ('default', 'default', 'default');

COMMIT;
