ALTER TABLE thing
    ADD COLUMN deleted BOOLEAN NOT null default false;
ALTER TABLE thing
    ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL;

CREATE INDEX idx_thing_deleted ON thing(deleted);
CREATE INDEX idx_thing_deleted_at ON thing(deleted_at);




