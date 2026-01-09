create schema mailer;

create table mailer.mails (
                              id character varying(255) not null,
                              recipient character varying(255) not null,
                              text text not null,
                              subject character varying(255) not null,
                              status character varying(255) not null,
                              created_at timestamp with time zone not null,
                              idempotency_key character varying(255) not null,
                              constraint mails_pkey primary key (id)
);

create index IF not exists mails_index_0 on mailer.mails using btree (idempotency_key) TABLESPACE pg_default;

create schema content;

create table content.folders (
                                 id uuid not null,
                                 name character varying(255) not null,
                                 display_name character varying(255) not null,
                                 constraint folders_pkey primary key (id)
);

create table content.content_types (
                                       name character varying(255) not null,
                                       constraint content_types_pkey primary key (name)
);

create table content.content (
                                 id uuid not null,
                                 user_id uuid not null,
                                 display_name character varying(255) not null,
                                 text character varying(255) null,
                                 media_url character varying(255) not null,
                                 type character varying(255) not null,
                                 deleted_at timestamp with time zone null,
                                 created_at timestamp with time zone not null,
                                 constraint content_pkey primary key (id),
                                 constraint content_type_fkey foreign KEY (type) references content.content_types (name)
);

create index IF not exists content_index_0 on content.content using btree (user_id, created_at desc) TABLESPACE pg_default;

create table content.folders_contents (
                                          folder_id uuid not null,
                                          content_id uuid not null,
                                          created_at timestamp with time zone null,
                                          constraint folders_contents_pkey primary key (folder_id, content_id),
                                          constraint folders_contents_content_id_fkey foreign KEY (content_id) references content.content (id) on update CASCADE on delete CASCADE,
                                          constraint folders_contents_folder_id_fkey foreign KEY (folder_id) references content.folders (id) on update CASCADE on delete CASCADE
);

create schema authorization_service;

create table authorization_service.roles (
                                             id uuid not null,
                                             name character varying(255) not null,
                                             level integer not null default 0,
                                             constraint roles_pkey primary key (id)
);

create table authorization_service.tokens (
                                              id uuid not null,
                                              user_id uuid not null,
                                              provider_name character varying(255) not null,
                                              token character varying(255) not null,
                                              expires_at timestamp with time zone not null,
                                              replaced_by_token uuid null,
                                              revoked_by_ip character varying(20) null,
                                              revoked_at timestamp with time zone null,
                                              created_at timestamp with time zone not null,
                                              constraint tokens_pkey primary key (id),
                                              constraint tokens_replaced_by_token_fkey foreign KEY (replaced_by_token) references authorization_service.tokens (id) on update CASCADE on delete CASCADE,
                                              constraint tokens_user_id_fkey foreign KEY (user_id) references authorization_service.users (id) on update CASCADE on delete CASCADE
);

create table authorization_service.users (
                                             id uuid not null,
                                             nickname character varying(50) not null,
                                             email character varying(255) not null,
                                             password_hash character varying(255) not null,
                                             code character varying(255) null,
                                             code_requested_at timestamp with time zone null,
                                             is_confirmed boolean not null default false,
                                             created_at timestamp with time zone not null,
                                             banned_before timestamp with time zone null,
                                             role_id uuid null,
                                             constraint users_pkey primary key (id),
                                             constraint users_role_id_fkey foreign KEY (role_id) references authorization_service.roles (id)
);

