# SQL files

These create the table structure and procedures required by the application. The compose file mounts this folder directly to /docker-entrypoint-initdb.d, so any scripts contained within will be run when the postgres container is initialized.
