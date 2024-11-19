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
  "is_upload" bool,
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

CREATE TABLE auth_details (
    "auth_id" int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "description" TEXT,
    "vault_path" TEXT,
    "date_added" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "added_by" TEXT
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

CREATE TABLE logs (
    "id" int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "timestamp" TIMESTAMPTZ NOT NULL,
    "message" TEXT NOT NULL,
    "stack_trace" TEXT,
    "action_id" INT,
    "additional_info" TEXT
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
    p_is_upload bool,
    p_connection_id INT,
    p_project_id INT,
    p_date_added TIMESTAMPTZ DEFAULT Now(),
    p_added_by TEXT DEFAULT CURRENT_USER
) AS $$
BEGIN
    INSERT INTO public.actions (local_path, remote_path, bucket, is_upload, connection_id, project_id, date_added, added_by)
    VALUES (p_local_path, p_remote_path, p_bucket, p_is_upload, p_connection_id, p_project_id, p_date_added, p_added_by);
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

CREATE OR REPLACE PROCEDURE usp_insert_auth_detail(
    p_description TEXT,
    p_vault_path TEXT,
    p_date_added TIMESTAMPTZ DEFAULT Now(),
    p_added_by TEXT DEFAULT CURRENT_USER
) AS $$
BEGIN
    INSERT INTO public.auth_details (description, vault_path, date_added, added_by)
    VALUES (p_description, p_vault_path, p_date_added, p_added_by);
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

-- DROP FUNCTION public.usp_collect_schedule_actions(int4);

CREATE OR REPLACE FUNCTION public.usp_collect_schedule_actions(schedule_id_input integer)
 RETURNS TABLE(
    schedule_id integer, 
    schedule_description text, 
    schedule_date_added timestamp with time zone, 
    schedule_added_by text, 
    project_id integer, 
    project_name text, 
    project_description text, 
    project_date_added timestamp with time zone, 
    project_added_by text, 
    action_id integer, 
    action_local_path text, 
    action_remote_path text, 
    action_bucket text, 
    action_upload bool,
    action_date_added timestamp with time zone, 
    action_added_by text, 
    connection_id integer, 
    connection_date_added timestamp with time zone, 
    connection_added_by text, 
    protocol_id integer, 
    protocol_name text, 
    protocol_description text, 
    protocol_date_added timestamp with time zone, 
    protocol_added_by text, 
    config_id integer, 
    config_key text, 
    config_value text, 
    config_date_added timestamp with time zone, 
    config_added_by text, 
    auth_id integer, 
    auth_description text, 
    auth_vault_path text, 
    auth_date_added timestamp with time zone, 
    auth_added_by text
)
LANGUAGE plpgsql
AS $function$
BEGIN
    RETURN QUERY
    SELECT 
        s.schedule_id,
        s.description,
        s.date_added,
        s.added_by,
        p.project_id,
        p.name,
        p.description,
        p.date_added,
        p.added_by,
        a.action_id,
        a.local_path,
        a.remote_path,
        a.bucket,
		a.is_upload,
        a.date_added,
        a.added_by,
        c.connection_id,
        c.date_added,
        c.added_by,
        pr.protocol_id,
        pr.name,
        pr.description,
        pr.date_added,
        pr.added_by,
        cc.config_id,
        cc.key,
        cc.value,
        cc.date_added,
        cc.added_by,
        ad.auth_id,
        ad.description,
        ad.vault_path,
        ad.date_added,
        ad.added_by
    FROM schedule s
    JOIN schedule_actions sa ON s.schedule_id = sa.schedule_id
    JOIN actions a ON sa.action_id = a.action_id
    JOIN projects p ON a.project_id = p.project_id
    JOIN connections c ON a.connection_id = c.connection_id
    JOIN protocols pr ON c.protocol_id = pr.protocol_id
    JOIN auth_details ad ON c.auth_id = ad.auth_id
    LEFT JOIN connection_config cc ON c.connection_id = cc.connection_id
    WHERE s.schedule_id = schedule_id_input;
END;
$function$
;
