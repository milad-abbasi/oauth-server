CREATE TABLE IF NOT EXISTS apps(
   id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
   secret VARCHAR NOT NULL,
   owner_id UUID NOT NULL,
   created_at TIMESTAMP DEFAULT NOW() NOT NULL,
   updated_at TIMESTAMP DEFAULT NOW() NOT NULL,
   CONSTRAINT fk_apps_users
       FOREIGN KEY (owner_id)
       REFERENCES users(id)
       ON DELETE CASCADE
);
