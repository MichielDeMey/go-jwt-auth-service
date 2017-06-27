INSERT INTO "public"."users"("email", "password") VALUES('michiel@demey.io', crypt('test123', gen_salt('bf', 8))) RETURNING "id", "email", "password";
INSERT INTO "public"."users"("email", "password") VALUES('gilles@demey.io', crypt('test456', gen_salt('bf', 8))) RETURNING "id", "email", "password";
