-- name: CountIndexFiles :one
-- Returns number of index files with given hash
select count(*) from anidb_index_files
where name = $1;

-- name: AddIndexFile :exec
-- Adds new index file with given name and hash
-- Does nothing if index file already exists
insert into anidb_index_files (name, hash)
values ($1, $2)
on conflict do nothing;
