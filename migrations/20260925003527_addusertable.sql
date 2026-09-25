-- Create "users" table
CREATE TABLE "users" (
  "id" character varying NOT NULL,
  "name" character varying NOT NULL,
  "email" character varying NOT NULL,
  "password_hash" character varying NOT NULL,
  "role" character varying NOT NULL DEFAULT 'user',
  "active" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "user_created_at" to table: "users"
CREATE INDEX "user_created_at" ON "users" ("created_at");
-- Create index "user_email" to table: "users"
CREATE UNIQUE INDEX "user_email" ON "users" ("email");
