CREATE TABLE "sessions" (
                         "id" uuid PRIMARY KEY,
                         "username" varchar(100) NOT NULL,
                         "refresh_token" varchar(2048) NOT NULL,
                         "user_agent" varchar(255)  NOT NULL,
                         "client_ip" varchar(255) not null ,
                         "is_blocked" boolean not null default false,
                         "created_at" timestamptz DEFAULT (now()),
                         "expires_at" timestamptz not null
);

alter table "sessions" add foreign key ("username") references "users" ("username")

