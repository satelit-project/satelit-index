--
-- PostgreSQL database dump
--

-- Dumped from database version 12beta3
-- Dumped by pg_dump version 12beta3

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: manage_anidb_index_files_limit(integer); Type: FUNCTION; Schema: public; Owner: test
--

CREATE FUNCTION public.manage_anidb_index_files_limit(_limit integer) RETURNS void
    LANGUAGE plpgsql
    AS $_$
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
$_$;


ALTER FUNCTION public.manage_anidb_index_files_limit(_limit integer) OWNER TO test;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: anidb_index_files; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.anidb_index_files (
    id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    name text NOT NULL,
    hash text NOT NULL
);


ALTER TABLE public.anidb_index_files OWNER TO test;

--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.schema_migration (
    version character varying(14) NOT NULL
);


ALTER TABLE public.schema_migration OWNER TO test;

--
-- Name: anidb_index_files anidb_index_files_pk; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.anidb_index_files
    ADD CONSTRAINT anidb_index_files_pk PRIMARY KEY (id);


--
-- Name: anidb_index_files_hash_uindex; Type: INDEX; Schema: public; Owner: test
--

CREATE UNIQUE INDEX anidb_index_files_hash_uindex ON public.anidb_index_files USING btree (hash);


--
-- Name: anidb_index_files_id_uindex; Type: INDEX; Schema: public; Owner: test
--

CREATE UNIQUE INDEX anidb_index_files_id_uindex ON public.anidb_index_files USING btree (id);


--
-- Name: anidb_index_files_name_uindex; Type: INDEX; Schema: public; Owner: test
--

CREATE UNIQUE INDEX anidb_index_files_name_uindex ON public.anidb_index_files USING btree (name);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- PostgreSQL database dump complete
--

