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

-- functions

/*
 * Creates trigger that will run on each insert into `anidb_index_files`
 * to cleanup old rows.
 *
 * Arguments:
 *  - _limit: how many of newer rows to keep
 */
create or replace function manage_anidb_index_files_limit(_limit int) returns void as
$$
begin
    -- creates a trigger function to cleanup most old rows
    execute format($q$
            create or replace function cleanup_anidb_index_files() returns trigger as
            $qq$
            begin
                delete
                from anidb_index_files
                where id not in (
                    select id
                    from anidb_index_files
                    order by updated_at desc
                    limit %s
                );

                return null;
            end;
            $qq$ language plpgsql;
        $q$, _limit);

    -- creates an insert trigger to run cleanup function
    execute $q$
        drop trigger if exists start_cleanup_index_files on anidb_index_files;
        create trigger start_cleanup_index_files
            after insert
            on anidb_index_files
            for each statement
        execute procedure cleanup_anidb_index_files();
    $q$;

end;
$$ language plpgsql;
