-- Create "course_categories" table
CREATE TABLE "course_categories" (
  "id" character varying NOT NULL,
  "tenant_id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "code" character varying NOT NULL,
  "name" character varying NOT NULL,
  "description" text NULL DEFAULT '',
  "parent_id" character varying NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "course_modules" table
CREATE TABLE "course_modules" (
  "id" character varying NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "course_templates" table
CREATE TABLE "course_templates" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "code" character varying NOT NULL,
  "title" character varying NOT NULL,
  "description" text NULL DEFAULT '',
  PRIMARY KEY ("id")
);
-- Create index "course_templates_code_key" to table: "course_templates"
CREATE UNIQUE INDEX "course_templates_code_key" ON "course_templates" ("code");
-- Create "course_versions" table
CREATE TABLE "course_versions" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "version_tag" character varying NOT NULL,
  "status" character varying NOT NULL,
  "estimated_hours" bigint NOT NULL DEFAULT 0,
  "template_id" character varying NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "course_versions_course_templates_versions" FOREIGN KEY ("template_id") REFERENCES "course_templates" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "syllabis" table
CREATE TABLE "syllabis" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "objectives" text NULL DEFAULT '',
  "entry_profile" text NULL DEFAULT '',
  "exit_profile" text NULL DEFAULT '',
  "methodology" character varying NULL DEFAULT '',
  "durations_hours" bigint NOT NULL DEFAULT 0,
  "version_id" character varying NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "syllabis_course_versions_syllabi" FOREIGN KEY ("version_id") REFERENCES "course_versions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
