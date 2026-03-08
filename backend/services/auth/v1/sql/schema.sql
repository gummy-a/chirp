--
-- PostgreSQL database dump
--

\restrict RkYaik88Z154zCqt45IQcYPOV70msYaOHhUIL3MhBpbNUCUJL8wpILqTbOrEAzO

-- Dumped from database version 18.1 (Debian 18.1-1.pgdg13+2)
-- Dumped by pg_dump version 18.1 (Debian 18.1-1.pgdg13+2)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
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
-- Name: accounts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.accounts (
    id uuid DEFAULT uuidv7() NOT NULL,
    email text NOT NULL,
    password_hash text NOT NULL,
    password_algorithm text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


--
-- Name: temporary_accounts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.temporary_accounts (
    id uuid DEFAULT uuidv7() NOT NULL,
    email text NOT NULL,
    password_hash text NOT NULL,
    number_code integer NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: accounts accounts_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.accounts
    ADD CONSTRAINT accounts_email_key UNIQUE (email);


--
-- Name: accounts accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.accounts
    ADD CONSTRAINT accounts_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: temporary_accounts temporary_accounts_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.temporary_accounts
    ADD CONSTRAINT temporary_accounts_email_key UNIQUE (email);


--
-- Name: temporary_accounts temporary_accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.temporary_accounts
    ADD CONSTRAINT temporary_accounts_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

\unrestrict RkYaik88Z154zCqt45IQcYPOV70msYaOHhUIL3MhBpbNUCUJL8wpILqTbOrEAzO

