CREATE TABLE
    public.daily (
         id bigint NOT NULL DEFAULT unique_rowid(),
         address character varying(41) NOT NULL,
         earnings real NOT NULL DEFAULT 0.0,
         date character varying(10) NOT NULL,
         PRIMARY KEY (id)
);

ALTER TABLE public.daily
    ADD CONSTRAINT daily_address_date_unique UNIQUE (address, date);

CREATE TABLE
    public.wallet (
          id bigint NOT NULL DEFAULT unique_rowid(),
          name character varying(24) NOT NULL,
          address character varying(41) NOT NULL,
          nodes smallint NOT NULL DEFAULT 0,
          balance real NOT NULL DEFAULT 0.0,
          daily real NOT NULL DEFAULT 0.0,
          PRIMARY KEY (id)
);