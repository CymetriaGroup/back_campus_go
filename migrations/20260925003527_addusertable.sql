-- Create "users" table
CREATE TABLE IF NOT EXISTS "users" (
  "id" character varying NOT NULL,
  "tenant_id" character varying NOT NULL,
  "auth_id" character varying NOT NULL,
  "name" character varying NOT NULL,
  "email" character varying NOT NULL,
  "role" character varying NOT NULL DEFAULT 'user',
  "status" character varying NOT NULL DEFAULT 'ACTIVE',
  "active" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
CREATE INDEX IF NOT EXISTS "user_created_at" ON "users" ("created_at");
CREATE UNIQUE INDEX IF NOT EXISTS "user_email" ON "users" ("email");
CREATE UNIQUE INDEX IF NOT EXISTS "user_auth_id" ON "users" ("auth_id");
CREATE INDEX IF NOT EXISTS "user_status" ON "users" ("status");

-- Create "user_profiles" table
CREATE TABLE IF NOT EXISTS "user_profiles" (
  "id" character varying NOT NULL,
  "user_id" character varying NOT NULL,
  "first_name" character varying NULL,
  "last_name" character varying NULL,
  "birth_date" timestamptz NULL,
  "gender" character varying NULL,
  "city" character varying NULL,
  "address" character varying NULL,
  "photo_url" character varying NULL,
  "sociodemographic_data" jsonb NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_user_profile_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE
);
CREATE UNIQUE INDEX IF NOT EXISTS "user_profile_user_id" ON "user_profiles" ("user_id");

-- Create "auth_providers" table
CREATE TABLE IF NOT EXISTS "auth_providers" (
  "id" character varying NOT NULL,
  "provider_type" character varying NOT NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "config" jsonb NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
CREATE UNIQUE INDEX IF NOT EXISTS "auth_provider_type" ON "auth_providers" ("provider_type");

-- Create "auth_sessions" table
CREATE TABLE IF NOT EXISTS "auth_sessions" (
  "id" character varying NOT NULL,
  "user_id" character varying NOT NULL,
  "auth_provider_id" character varying NULL,
  "token_jti" character varying NOT NULL,
  "ip_address" character varying NULL,
  "user_agent" character varying NULL,
  "expires_at" timestamptz NOT NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_auth_session_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE,
  CONSTRAINT "fk_auth_session_provider" FOREIGN KEY ("auth_provider_id") REFERENCES "auth_providers" ("id") ON DELETE SET NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS "auth_session_token_jti" ON "auth_sessions" ("token_jti");
CREATE INDEX IF NOT EXISTS "auth_session_is_active" ON "auth_sessions" ("is_active");

-- Create "security_events" table
CREATE TABLE IF NOT EXISTS "security_events" (
  "id" character varying NOT NULL,
  "user_id" character varying NULL,
  "event_type" character varying NOT NULL,
  "ip_address" character varying NULL,
  "details" jsonb NULL,
  "created_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_security_event_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS "security_event_type" ON "security_events" ("event_type");
CREATE INDEX IF NOT EXISTS "security_event_created_at" ON "security_events" ("created_at");

-- Create "roles" table
CREATE TABLE IF NOT EXISTS "roles" (
  "id" character varying NOT NULL,
  "name" character varying NOT NULL,
  "description" character varying NULL,
  "is_platform_level" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
CREATE UNIQUE INDEX IF NOT EXISTS "role_name" ON "roles" ("name");
CREATE INDEX IF NOT EXISTS "role_is_platform_level" ON "roles" ("is_platform_level");

-- Create "permissions" table
CREATE TABLE IF NOT EXISTS "permissions" (
  "id" character varying NOT NULL,
  "action" character varying NOT NULL,
  "description" character varying NULL,
  "created_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
CREATE UNIQUE INDEX IF NOT EXISTS "permission_action" ON "permissions" ("action");

-- Create "role_permissions" join table
CREATE TABLE IF NOT EXISTS "role_permissions" (
  "role_id" character varying NOT NULL,
  "permission_id" character varying NOT NULL,
  PRIMARY KEY ("role_id", "permission_id"),
  CONSTRAINT "fk_role_permissions_role" FOREIGN KEY ("role_id") REFERENCES "roles" ("id") ON DELETE CASCADE,
  CONSTRAINT "fk_role_permissions_permission" FOREIGN KEY ("permission_id") REFERENCES "permissions" ("id") ON DELETE CASCADE
);

-- Create "user_platform_roles" join table
CREATE TABLE IF NOT EXISTS "user_platform_roles" (
  "user_id" character varying NOT NULL,
  "role_id" character varying NOT NULL,
  PRIMARY KEY ("user_id", "role_id"),
  CONSTRAINT "fk_user_platform_roles_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE,
  CONSTRAINT "fk_user_platform_roles_role" FOREIGN KEY ("role_id") REFERENCES "roles" ("id") ON DELETE CASCADE
);
