BEGIN;

INSERT INTO namespace (id, is_default) VALUES ('default', true);
INSERT INTO device_set(id,namespace_id) VALUES ('default', 'default');

COMMIT;
