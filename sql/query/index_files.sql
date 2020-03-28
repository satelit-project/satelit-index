-- name: CountIndexFiles :one
-- Returns number of index files with given hash.
select count(*) from anidb_index_files
where hash = $1;

-- name: AddIndexFile :exec
-- Adds new index file with given hash and remote path or does nothing if index file already exists.
insert into anidb_index_files (hash, file_path)
values ($1, $2)
on conflict do nothing;

-- name: LatestIndexFile :one
-- Returns most recent index file record.
select * from anidb_index_files
order by created_at desc
limit 1;

-- name: IndexFileByHash :one
-- Returns index file record with specified hash.
select * from anidb_index_files
where hash = $1;
