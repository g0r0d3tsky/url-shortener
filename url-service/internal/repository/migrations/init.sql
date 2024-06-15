CREATE TABLE "data" (
  "id" uuid PRIMARY KEY,
  "longUrl" varchar(3000),
  "shortUrl" varchar,
  "expiresAt" timestamp without time zone
);

CREATE TABLE "keys" (
  "id" uuid PRIMARY KEY,
  "key" serial,
  "encode" varchar,
  "urlID" uuid
);

ALTER TABLE "keys" ADD FOREIGN KEY ("urlID") REFERENCES "data" ("id");
