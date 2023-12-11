CREATE TABLE "deposits" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigserial NOT NULL,
  "amount" bigint NOT NULL,
  "user" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "withdraws" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigserial NOT NULL,
  "amount" bigint NOT NULL,
  "user" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

ALTER TABLE "deposits" ADD FOREIGN KEY ("user") REFERENCES "users" ("username");

ALTER TABLE "withdraws" ADD FOREIGN KEY ("user") REFERENCES "users" ("username");
