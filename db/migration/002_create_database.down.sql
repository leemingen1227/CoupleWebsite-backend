REVOKE USAGE ON SCHEMA public FROM couple_readonly;
REVOKE USAGE ON SCHEMA public FROM couple_user;
REVOKE CREATE ON SCHEMA public FROM couple_admin;
REVOKE USAGE ON SCHEMA public FROM couple_admin;
ALTER DATABASE couple_db OWNER TO postgres;
DROP DATABASE IF EXISTS couple_db;