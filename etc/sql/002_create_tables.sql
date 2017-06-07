CREATE TABLE "public"."users" (
    "id" serial,
    "email" text NOT NULL,
    "password" text,
    PRIMARY KEY ("id"),
    UNIQUE ("email")
);
