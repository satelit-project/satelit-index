-- name: CountIndexFiles :one
-- Returns number of index files with given hash.
select count(*) from anidb_index_files
where hash = $1;

-- name: AddIndexFile :exec
-- Adds new index file with given name and hash or does nothing if index file already exists.
insert into anidb_index_files (hash)
values ($1)
on conflict do nothing;

-- name: LatestIndexFile :one
-- Returns most recent index file record.
select * from anidb_index_files
order by created_at desc
limit 1;
