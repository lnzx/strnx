CREATE TABLE "public"."daily" (
    "id" SERIAL PRIMARY KEY,
    "address" VARCHAR(41) not null,
    "earnings" REAL not null default 0,
    "date" DATE not null default CURRENT_DATE
);
CREATE UNIQUE INDEX "daily_addr_date" on "public"."daily"("address", "date");

CREATE TABLE "public"."wallet" (
      "id" SERIAL PRIMARY KEY,
      "name" VARCHAR(24) not null,
      "address" VARCHAR(41) not null,
      "nodes" INT2 [] not null DEFAULT '{0,0}',
      "balance" REAL not null default 0,
      "daily" REAL not null default 0,
      "group" VARCHAR(8) DEFAULT '-'
);
CREATE UNIQUE INDEX "wallet_addr" on "public"."wallet"("address");

CREATE TABLE public.earn (
    node_id VARCHAR(36) PRIMARY KEY,
    earning real NOT NULL DEFAULT 0,
    status VARCHAR(9) NOT NULL,
    isp VARCHAR(64),
    country VARCHAR(32),
    city VARCHAR(32),
    region VARCHAR(32),
    created date
);

CREATE TABLE public.node (
    id smallserial PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    ip VARCHAR(15) NOT NULL,
    bandwidth integer NOT NULL DEFAULT 0,
    traffic VARCHAR(15) DEFAULT '',
    price real NOT NULL DEFAULT 0,
    renew VARCHAR(10) DEFAULT '',
    state VARCHAR(12) DEFAULT '',
    node_id VARCHAR(36) DEFAULT '',
    type VARCHAR(8) DEFAULT '',
    cpu integer NOT NULL DEFAULT 0,
    ram VARCHAR(8) DEFAULT '',
    disk VARCHAR(32) DEFAULT '',
    pool_id integer NOT NULL DEFAULT 0
);
CREATE UNIQUE INDEX "node_ip" on "public"."node"("ip");

CREATE TABLE public.pool (
    id smallserial PRIMARY KEY,
    traffic integer NOT NULL DEFAULT 0,
    used integer NULL DEFAULT 0
);
