-- +goose Up

/*
* Stores info about all index files for Anidb anime titles.
*/
create table anidb_index_files
(
    id         uuid        default uuid_generate_v4() not null,
    hash       text                                   not null,
    url        text                                   not null,
    created_at timestamptz default now()              not null,
    updated_at timestamptz default now()              not null
);

-- hash should be unique
create unique index anidb_index_files_hash_uindex
    on anidb_index_files (hash);

-- as well as id
create unique index anidb_index_files_id_uindex
    on anidb_index_files (id);

-- make `id` a primary key
alter table anidb_index_files
    add constraint anidb_index_files_pk
        primary key (id);

-- autoupdate updated_at
select manage_updated_at('anidb_index_files');

-- +goose Down

select unmanage_updated_at('anidb_index_files');
drop table anidb_index_files;
