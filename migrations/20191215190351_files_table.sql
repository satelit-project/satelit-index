-- +goose Up

/*
* Stores info about all index files for Anidb anime titles.
*/
create table anidb_index_files
(
    id         uuid                      not null,
    created_at timestamptz default now() not null,
    updated_at timestamptz default now() not null,

    -- name of the stored index file
    name       text                      not null,

    -- it's hash
    hash       text                      not null
);

-- hash should be unique
create unique index anidb_index_files_hash_uindex
    on anidb_index_files (hash);

-- as well as id
create unique index anidb_index_files_id_uindex
    on anidb_index_files (id);

-- as well as the stored file name
create unique index anidb_index_files_name_uindex
    on anidb_index_files (name);

-- make `id` a primary key
alter table anidb_index_files
    add constraint anidb_index_files_pk
        primary key (id);

-- autoupdate updated_at
select manage_updated_at('anidb_index_files');

-- +goose Down

select unmanage_updated_at('anidb_index_files');
drop table anidb_index_files;
