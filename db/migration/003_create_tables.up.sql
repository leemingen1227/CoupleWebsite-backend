DROP TABLE IF EXISTS pairs CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS verify_emails CASCADE;
DROP TABLE IF EXISTS invitations CASCADE;
DROP TABLE IF EXISTS user_pairs CASCADE;
DROP TABLE IF EXISTS sessions CASCADE;

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

/*
    grant_table_privilege
*/
/*for normal tables */
ALTER TABLE pairs OWNER TO couple_admin;
ALTER TABLE users OWNER TO couple_admin;
ALTER TABLE verify_emails OWNER TO couple_admin;
ALTER TABLE invitations OWNER TO couple_admin;
ALTER TABLE user_pairs OWNER TO couple_admin;

GRANT SELECT, INSERT, UPDATE, DELETE, REFERENCES ON TABLE pairs to couple_user;
GRANT USAGE ON SEQUENCE pairs_id_seq TO couple_user;
GRANT SELECT, INSERT, UPDATE, DELETE, REFERENCES ON TABLE users to couple_user;
GRANT SELECT, INSERT, UPDATE, DELETE, REFERENCES ON TABLE verify_emails to couple_user;
GRANT USAGE ON SEQUENCE verify_emails_id_seq TO couple_user;
GRANT SELECT, INSERT, UPDATE, DELETE, REFERENCES ON TABLE invitations to couple_user;
GRANT USAGE ON SEQUENCE invitations_id_seq TO couple_user;
GRANT SELECT, INSERT, UPDATE, DELETE, REFERENCES ON TABLE user_pairs to couple_user;


GRANT SELECT ON TABLE pairs to couple_readonly;
GRANT SELECT ON TABLE users to couple_readonly;
GRANT SELECT ON TABLE verify_emails to couple_readonly;
GRANT SELECT ON TABLE invitations to couple_readonly;
GRANT SELECT ON TABLE user_pairs to couple_readonly;
