drop trigger if exists start_cleanup_index_files on anidb_index_files;
drop function if exists cleanup_anidb_index_files;
drop function if exists manage_anidb_index_files_limit;
drop table anidb_index_files;
