create table index_files
(
    id uuid not null,
    name text not null,
    hash text not null,
    created_at timestamptz default now() not null,
    updated_at timestamptz default now() not null
);

create unique index index_files_hash_uindex
    on index_files (hash);

create unique index index_files_id_uindex
    on index_files (id);

create unique index index_files_name_uindex
    on index_files (name);

alter table index_files
    add constraint index_files_pk
        primary key (id);
