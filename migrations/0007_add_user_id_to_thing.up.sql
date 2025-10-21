ALTER TABLE thing ADD COLUMN user_id BIGINT NOT NULL;
ALTER TABLE thing ADD CONSTRAINT fk_thing_user FOREIGN KEY (user_id) REFERENCES users(id);
CREATE INDEX idx_thing_user_id ON thing(user_id);