CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "balance" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "acc_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "from_acc_id" bigint NOT NULL,
  "to_acc_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "entries" ADD FOREIGN KEY ("acc_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_acc_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("to_acc_id") REFERENCES "accounts" ("id");

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "entries" ("acc_id");

CREATE INDEX ON "transfers" ("from_acc_id");

CREATE INDEX ON "transfers" ("to_acc_id");

CREATE INDEX ON "transfers" ("from_acc_id", "to_acc_id");