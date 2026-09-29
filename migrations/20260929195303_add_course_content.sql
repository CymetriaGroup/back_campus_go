-- Modify "course_modules" table
ALTER TABLE "course_modules" ADD COLUMN "created_at" timestamptz NOT NULL, ADD COLUMN "updated_at" timestamptz NOT NULL, ADD COLUMN "title" character varying NOT NULL, ADD COLUMN "sequence_order" bigint NOT NULL, ADD COLUMN "version_id" character varying NOT NULL, ADD CONSTRAINT "course_modules_course_versions_modules" FOREIGN KEY ("version_id") REFERENCES "course_versions" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Create index "coursemodules_version_id_sequence_order" to table: "course_modules"
CREATE UNIQUE INDEX "coursemodules_version_id_sequence_order" ON "course_modules" ("version_id", "sequence_order");
-- Create "lessons" table
CREATE TABLE "lessons" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "title" character varying NOT NULL,
  "sequence_order" bigint NOT NULL,
  "module_id" character varying NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "lessons_course_modules_lessons" FOREIGN KEY ("module_id") REFERENCES "course_modules" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "lessons_module_id_sequence_order" to table: "lessons"
CREATE UNIQUE INDEX "lessons_module_id_sequence_order" ON "lessons" ("module_id", "sequence_order");
-- Create "activities" table
CREATE TABLE "activities" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "title" character varying NOT NULL,
  "type" character varying NOT NULL,
  "is_required" boolean NOT NULL DEFAULT true,
  "sequence_order" bigint NOT NULL,
  "lesson_id" character varying NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "activities_lessons_activities" FOREIGN KEY ("lesson_id") REFERENCES "lessons" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "activities_lesson_id_sequence_order" to table: "activities"
CREATE UNIQUE INDEX "activities_lesson_id_sequence_order" ON "activities" ("lesson_id", "sequence_order");
-- Create "course_resources" table
CREATE TABLE "course_resources" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "type" character varying NOT NULL,
  "name" character varying NOT NULL,
  "url_storage_key" character varying NOT NULL,
  "mime_type" character varying NOT NULL,
  "size_bytes" bigint NOT NULL,
  "position" bigint NOT NULL,
  "metadata" jsonb NULL,
  "activity_id" character varying NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "course_resources_activities_resources" FOREIGN KEY ("activity_id") REFERENCES "activities" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "courseresources_activity_id_position" to table: "course_resources"
CREATE UNIQUE INDEX "courseresources_activity_id_position" ON "course_resources" ("activity_id", "position");
