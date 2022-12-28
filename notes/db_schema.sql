\connect bank

CREATE TABLE "accounts" (
                            "id" bigserial PRIMARY KEY,
                            "owner" varchar NOT NULL,
                            "balance" decimal NOT NULL,
                            "currency" varchar(3) NOT NULL,
                            "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "transfer" (
                            "id" bigserial PRIMARY KEY,
                            "from_account_id" bigint NOT NULL,
                            "to_account_id" bigint NOT NULL,
                            "created_at" timestamptz DEFAULT (now()),
                            "amount" decimal NOT NULL
);

CREATE TABLE "entries" (
                           "id" bigserial PRIMARY KEY,
                           "account_id" bigint NOT NULL,
                           "amount" decimal NOT NULL,
                           "created_at" timestamptz DEFAULT (now())
);

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "transfer" ("from_account_id");

CREATE INDEX ON "transfer" ("to_account_id");

CREATE INDEX ON "transfer" ("from_account_id", "to_account_id");

CREATE INDEX ON "entries" ("account_id");

COMMENT ON COLUMN "transfer"."amount" IS 'must be positive';

COMMENT ON COLUMN "entries"."amount" IS 'can be negative';

ALTER TABLE "transfer" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfer" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");
