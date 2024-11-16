CREATE TABLE "projects" (
  "project_id" int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  "name" text,
  "description" text,
  "date_added" timestamptz,
  "added_by" text
);

CREATE TABLE "connections" (
  "connection_id" int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  "protocol_id" int,
  "auth_id" int,
  "date_added" timestamptz,
  "added_by" text
);

CREATE TABLE "protocols" (
  "protocol_id" int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  "name" text,
  "description" text,
  "date_added" timestamptz,
  "added_by" text
);

CREATE TABLE "actions" (
  "action_id" int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  "local_path" text,
  "remote_path" text,
  "bucket" text,
  "connection_id" int,
  "project_id" int,
  "date_added" timestamptz,
  "added_by" text
);

CREATE TABLE "connection_config" (
  "config_id" int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  "key" text,
  "value" text,
  "connection_id" int,
  "date_added" timestamptz,
  "added_by" text
);

CREATE TABLE "auth_details" (
  "auth_id" int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  "description" text,
  "username" text,
  "password" text,
  "private_key" text,
  "access_key" text,
  "secret_key" text,
  "date_added" timestamptz,
  "added_by" text
);

CREATE TABLE "schedule" (
  "schedule_id" int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  "description" text,
  "date_added" timestamptz,
  "added_by" text
);

CREATE TABLE "schedule_actions" (
  "schedule_id" int,
  "action_id" int
);

ALTER TABLE "schedule_actions" ADD CONSTRAINT "pk_schedule_actions" PRIMARY KEY ("schedule_id", "action_id");

ALTER TABLE "connections" ADD FOREIGN KEY ("protocol_id") REFERENCES "protocols" ("protocol_id");

ALTER TABLE "connections" ADD FOREIGN KEY ("auth_id") REFERENCES "auth_details" ("auth_id");

ALTER TABLE "actions" ADD FOREIGN KEY ("connection_id") REFERENCES "connections" ("connection_id");

ALTER TABLE "actions" ADD FOREIGN KEY ("project_id") REFERENCES "projects" ("project_id");

ALTER TABLE "connection_config" ADD FOREIGN KEY ("connection_id") REFERENCES "connections" ("connection_id");

ALTER TABLE "schedule_actions" ADD FOREIGN KEY ("schedule_id") REFERENCES "schedule" ("schedule_id");

ALTER TABLE "schedule_actions" ADD FOREIGN KEY ("action_id") REFERENCES "actions" ("action_id");


-- Procedure for inserting into projects
CREATE OR REPLACE PROCEDURE usp_insert_project(
    p_name TEXT,
    p_description TEXT,
    p_date_added TIMESTAMPTZ DEFAULT Now(),
    p_added_by TEXT DEFAULT CURRENT_USER
) AS $$
BEGIN
    INSERT INTO public.projects (name, description, date_added, added_by)
    VALUES (p_name, p_description, p_date_added, p_added_by);
END;
$$ LANGUAGE plpgsql;

-- Procedure for inserting into connections
CREATE OR REPLACE PROCEDURE usp_insert_connection(
    p_protocol_id INT,
    p_auth_id INT,
    p_date_added TIMESTAMPTZ DEFAULT Now(),
    p_added_by TEXT DEFAULT CURRENT_USER
) AS $$
BEGIN
    INSERT INTO public.connections (protocol_id, auth_id, date_added, added_by)
    VALUES (p_protocol_id, p_auth_id, p_date_added, p_added_by);
END;
$$ LANGUAGE plpgsql;

-- Procedure for inserting into protocols
CREATE OR REPLACE PROCEDURE usp_insert_protocol(
    p_name TEXT,
    p_description TEXT,
    p_date_added TIMESTAMPTZ DEFAULT Now(),
    p_added_by TEXT DEFAULT CURRENT_USER
) AS $$
BEGIN
    INSERT INTO public.protocols (name, description, date_added, added_by)
    VALUES (p_name, p_description, p_date_added, p_added_by);
END;
$$ LANGUAGE plpgsql;

-- Procedure for inserting into actions
CREATE OR REPLACE PROCEDURE usp_insert_action(
    p_local_path TEXT,
    p_remote_path TEXT,
    p_bucket TEXT,
    p_connection_id INT,
    p_project_id INT,
    p_date_added TIMESTAMPTZ DEFAULT Now(),
    p_added_by TEXT DEFAULT CURRENT_USER
) AS $$
BEGIN
    INSERT INTO public.actions (local_path, remote_path, bucket, connection_id, project_id, date_added, added_by)
    VALUES (p_local_path, p_remote_path, p_bucket, p_connection_id, p_project_id, p_date_added, p_added_by);
END;
$$ LANGUAGE plpgsql;

-- Procedure for inserting into connection_config
CREATE OR REPLACE PROCEDURE usp_insert_connection_config(
    p_key TEXT,
    p_value TEXT,
    p_connection_id INT,
    p_date_added TIMESTAMPTZ DEFAULT Now(),
    p_added_by TEXT DEFAULT CURRENT_USER
) AS $$
BEGIN
    INSERT INTO public.connection_config (key, value, connection_id, date_added, added_by)
    VALUES (p_key, p_value, p_connection_id, p_date_added, p_added_by);
END;
$$ LANGUAGE plpgsql;

-- Procedure for inserting into auth_details
CREATE OR REPLACE PROCEDURE usp_insert_auth_detail(
    p_description TEXT,
    p_username TEXT,
    p_password TEXT DEFAULT NULL,
    p_private_key TEXT DEFAULT NULL,
    p_access_key TEXT DEFAULT NULL,
    p_secret_key TEXT DEFAULT NULL,
    p_date_added TIMESTAMPTZ DEFAULT Now(),
    p_added_by TEXT DEFAULT CURRENT_USER
) AS $$
BEGIN
    INSERT INTO public.auth_details (description, username, password, private_key, access_key, secret_key, date_added, added_by)
    VALUES (p_description, p_username, p_password, p_private_key, p_access_key, p_secret_key, p_date_added, p_added_by);
END;
$$ LANGUAGE plpgsql;

-- Procedure for inserting into schedule
CREATE OR REPLACE PROCEDURE usp_insert_schedule(
    p_description TEXT,
    p_date_added TIMESTAMPTZ DEFAULT Now(),
    p_added_by TEXT DEFAULT CURRENT_USER
) AS $$
BEGIN
    INSERT INTO public.schedule (description, date_added, added_by)
    VALUES (p_description, p_date_added, p_added_by);
END;
$$ LANGUAGE plpgsql;

-- Procedure for inserting into schedule_actions
CREATE OR REPLACE PROCEDURE usp_insert_schedule_action(
    p_schedule_id INT,
    p_action_id INT
) AS $$
BEGIN
    INSERT INTO public.schedule_actions (schedule_id, action_id)
    VALUES (p_schedule_id, p_action_id);
END;
$$ LANGUAGE plpgsql;

