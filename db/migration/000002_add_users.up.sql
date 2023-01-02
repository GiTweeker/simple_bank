CREATE TABLE "users" (
                         "username" varchar(100) PRIMARY KEY,
                         "hashed_password" varchar(2048) NOT NULL,
                         "full_name" varchar(255) NOT NULL,
                         "email" varchar UNIQUE NOT NULL,
                         "created_at" timestamptz DEFAULT (now()),
                         "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00Z'
);
ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

--CREATE unique INDEX ON "accounts" ("owner", "currency");
alter table "accounts" add constraint "owner_currency_key" unique ("owner", "currency")