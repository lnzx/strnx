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
      "nodes" INT2 not null default 0,
      "balance" REAL not null default 0,
      "daily" REAL not null default 0
);
