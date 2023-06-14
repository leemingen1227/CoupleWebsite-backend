/*
    There are three different role in the database "couple_db"
    couple_admin is the owner of the database. This role should only be used when DBA is changing 
    the schema.
    couple_user is the role used by golang executable, while acting like a normal user for the website.
    Only CRUD privilege is granted, truncate table should not be granted.
    couple_readonly is used during debugging, which only have select privilege.
*/

CREATE ROLE couple_admin LOGIN PASSWORD 'admin_password' NOSUPERUSER INHERIT NOCREATEDB NOCREATEROLE NOREPLICATION;
CREATE ROLE couple_user LOGIN PASSWORD 'user_password' NOSUPERUSER INHERIT NOCREATEDB NOCREATEROLE NOREPLICATION;
CREATE ROLE couple_readonly LOGIN PASSWORD 'readonly_password' NOSUPERUSER INHERIT NOCREATEDB NOCREATEROLE NOREPLICATION;



/* setup db */
CREATE DATABASE couple_db with ENCODING = 'UTF8' LC_COLLATE = 'en_US.UTF-8' LC_CTYPE = 'en_US.UTF-8' CONNECTION LIMIT = -1 template=template0;

ALTER DATABASE couple_db OWNER TO couple_admin;

ALTER DATABASE couple_db SET timezone TO 'UTC';

REVOKE USAGE ON SCHEMA public FROM PUBLIC;
REVOKE CREATE ON SCHEMA public FROM PUBLIC;

GRANT USAGE ON SCHEMA public to couple_admin;
GRANT CREATE ON SCHEMA public to couple_admin;

GRANT USAGE ON SCHEMA public to couple_user;

GRANT USAGE ON SCHEMA public to couple_readonly;

/*create table*/
--the script will remove all the table in the database
\connect couple_db
DROP TABLE IF EXISTS pairs CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS verify_emails CASCADE;
DROP TABLE IF EXISTS invitations CASCADE;
DROP TABLE IF EXISTS user_pairs CASCADE;

CREATE TABLE pairs
(
    id BIGSERIAL PRIMARY KEY,
    create_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT current_timestamp,
    update_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT current_timestamp,
    start_date DATE
);


create table users
(
    id uuid,
    email character varying(200) not null,
    password_digest character varying(1000) not null,
    name character varying(255) not null,
    is_email_verified boolean not null default false,
    pair_id bigint,
    create_time timestamp without time zone not null default current_timestamp,
    update_time timestamp without time zone not null default current_timestamp,

    CONSTRAINT "users_pk" PRIMARY KEY (id),
    CONSTRAINT "users_u1" UNIQUE (email),
    CONSTRAINT "users_pair_id_fk" FOREIGN KEY (pair_id) REFERENCES pairs (id) ON DELETE CASCADE ON UPDATE CASCADE
);

create table verify_emails 
(
    id bigserial not null,
    user_id uuid not null,
    email character varying(200) not null,
    secret_code varchar not null,
    is_used boolean not null default false,
    create_time timestamp without time zone not null default current_timestamp,
    expired_time timestamp without time zone not null default (now() + interval '15 minutes'),
    CONSTRAINT "verify_emails_pk" PRIMARY KEY (id)
);
ALTER TABLE verify_emails ADD FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE;

create table invitations
(
    id bigserial not null,
    inviter_id uuid not null,
    invitee_email character varying(200) not null,
    invitation_token varchar not null,
    is_accepted boolean not null default false,
    create_time timestamp without time zone not null default current_timestamp,
    CONSTRAINT "invitations_pk" PRIMARY KEY (id)
);
ALTER TABLE invitations ADD FOREIGN KEY (inviter_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE;

CREATE TABLE user_pairs
(
    pair_id bigint NOT NULL,
    user_id uuid NOT NULL,

    CONSTRAINT user_pairs_pk PRIMARY KEY (pair_id, user_id),
    CONSTRAINT user_pairs_user_id_fk FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT user_pairs_pair_id_fk FOREIGN KEY (pair_id) REFERENCES pairs (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE sessions 
(
  id uuid PRIMARY KEY, 
  email varchar NOT NULL,
  refresh_token varchar NOT NULL,
  user_agent varchar NOT NULL,
  client_ip varchar NOT NULL,
  is_blocked boolean NOT NULL DEFAULT false,
  expires_at timestamptz NOT NULL, 
  created_at timestamptz NOT NULL DEFAULT (now())
);
ALTER TABLE sessions ADD FOREIGN KEY (email) REFERENCES users (email) ON DELETE CASCADE ON UPDATE CASCADE;

CREATE TABLE blog
(
    id uuid PRIMARY KEY,
    user_id uuid NOT NULL,
    pair_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    picture VARCHAR(255) NOT NULL, -- This field stores the identifier of the picture in S3
    create_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT current_timestamp,
    update_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT current_timestamp,
    
    CONSTRAINT blog_pair_id_fk FOREIGN KEY (pair_id) REFERENCES pairs (id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT blog_user_id_fk FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE
);


\connect couple_db
/*
    grant_table_privilege
*/
/*for normal tables */
ALTER TABLE pairs OWNER TO couple_admin;
ALTER TABLE users OWNER TO couple_admin;
ALTER TABLE verify_emails OWNER TO couple_admin;
ALTER TABLE invitations OWNER TO couple_admin;
ALTER TABLE user_pairs OWNER TO couple_admin;
ALTER TABLE sessions OWNER TO couple_admin;
ALTER TABLE blog OWNER TO couple_admin;


GRANT SELECT, INSERT, UPDATE, DELETE, REFERENCES ON TABLE pairs to couple_user;
GRANT USAGE ON SEQUENCE pairs_id_seq TO couple_user;
GRANT SELECT, INSERT, UPDATE, DELETE, REFERENCES ON TABLE users to couple_user;
GRANT SELECT, INSERT, UPDATE, DELETE, REFERENCES ON TABLE verify_emails to couple_user;
GRANT USAGE ON SEQUENCE verify_emails_id_seq TO couple_user;
GRANT SELECT, INSERT, UPDATE, DELETE, REFERENCES ON TABLE invitations to couple_user;
GRANT USAGE ON SEQUENCE invitations_id_seq TO couple_user;
GRANT SELECT, INSERT, UPDATE, DELETE, REFERENCES ON TABLE user_pairs to couple_user;
GRANT SELECT, INSERT, UPDATE, DELETE, REFERENCES ON TABLE sessions to couple_user;
GRANT SELECT, INSERT, UPDATE, DELETE, REFERENCES ON TABLE blog to couple_user;


GRANT SELECT ON TABLE pairs to couple_readonly;
GRANT SELECT ON TABLE users to couple_readonly;
GRANT SELECT ON TABLE verify_emails to couple_readonly;
GRANT SELECT ON TABLE invitations to couple_readonly;
GRANT SELECT ON TABLE user_pairs to couple_readonly;
GRANT SELECT ON TABLE sessions to couple_readonly;
GRANT SELECT ON TABLE blog to couple_readonly;
