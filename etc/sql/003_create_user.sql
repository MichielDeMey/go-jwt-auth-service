INSERT INTO "public"."users"("email", "password") VALUES('de.mey.michiel@gmail.com', crypt('test123', gen_salt('bf', 8))) RETURNING "id", "email", "password";
INSERT INTO "public"."users"("email", "password") VALUES('bino@vito.be', crypt('test456', gen_salt('bf', 8))) RETURNING "id", "email", "password";
