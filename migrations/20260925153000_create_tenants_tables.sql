-- Create "tenants" table
CREATE TABLE IF NOT EXISTS "tenants" (
  "id" character varying NOT NULL,
  "name" character varying NOT NULL,
  "branding" jsonb NULL,
  "domains" jsonb NULL,
  "limits" jsonb NULL,
  "features" jsonb NULL,
  "status" character varying NOT NULL DEFAULT 'ACTIVE',
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);

-- Create index "tenant_status" to table: "tenants"
CREATE INDEX IF NOT EXISTS "tenant_status" ON "tenants" ("status");

-- Create "tenant_settings" table
CREATE TABLE IF NOT EXISTS "tenant_settings" (
  "id" character varying NOT NULL,
  "tenant_id" character varying NOT NULL,
  "timezone" character varying NOT NULL DEFAULT 'UTC',
  "language" character varying NOT NULL DEFAULT 'es',
  "date_format" character varying NOT NULL DEFAULT 'YYYY-MM-DD',
  "institutional_email" character varying NOT NULL DEFAULT '',
  "policies" jsonb NULL,
  "feature_flags" jsonb NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_tenant_settings_tenant" FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id") ON DELETE CASCADE
);

-- Create index "tenant_settings_tenant_id" to table: "tenant_settings"
CREATE UNIQUE INDEX IF NOT EXISTS "tenant_settings_tenant_id" ON "tenant_settings" ("tenant_id");
