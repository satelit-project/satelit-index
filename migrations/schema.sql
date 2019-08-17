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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: index_files; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.index_files (
    id uuid NOT NULL,
    name text NOT NULL,
    hash text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.index_files OWNER TO test;

--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.schema_migration (
    version character varying(14) NOT NULL
);


ALTER TABLE public.schema_migration OWNER TO test;

--
-- Name: index_files index_files_pk; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.index_files
    ADD CONSTRAINT index_files_pk PRIMARY KEY (id);


--
-- Name: index_files_hash_uindex; Type: INDEX; Schema: public; Owner: test
--

CREATE UNIQUE INDEX index_files_hash_uindex ON public.index_files USING btree (hash);


--
-- Name: index_files_id_uindex; Type: INDEX; Schema: public; Owner: test
--

CREATE UNIQUE INDEX index_files_id_uindex ON public.index_files USING btree (id);


--
-- Name: index_files_name_uindex; Type: INDEX; Schema: public; Owner: test
--

CREATE UNIQUE INDEX index_files_name_uindex ON public.index_files USING btree (name);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: test
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- PostgreSQL database dump complete
--

