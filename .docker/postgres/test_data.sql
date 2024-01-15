-- Sequences
CREATE SEQUENCE IF NOT EXISTS elements_id_seq;

-- Table Definition
CREATE TABLE "public"."elements" (
    "id" int8 NOT NULL DEFAULT nextval('elements_id_seq'::regclass),
    "type" text,
    "x" int8,
    "y" int8,
    "width" int8,
    "height" int8,
    "value_from" text,
    "font" text,
    "font_size" int8,
    "pdf_template_id" int8,
    PRIMARY KEY ("id")
);

-- Sequences
CREATE SEQUENCE IF NOT EXISTS pdf_templates_id_seq;

-- Table Definition
CREATE TABLE "public"."pdf_templates" (
    "id" int8 NOT NULL DEFAULT nextval('pdf_templates_id_seq'::regclass),
    "updated_at" timestamptz,
    "created_at" timestamptz,
    "name" text,
    PRIMARY KEY ("id")
);


-- Table Definition
CREATE TABLE "public"."users" (
    "email" varchar NOT NULL,
    "password" bytea NOT NULL,
    PRIMARY KEY ("email")
);

INSERT INTO "public"."elements" ("id","type","x","y","width","height","value_from","font","font_size","pdf_template_id") VALUES (1,'text',10,10,0,0,'value','Arial',20,1);

INSERT INTO "public"."pdf_templates" ("id","updated_at","created_at","name") VALUES (1,'2024-01-15 19:43:39.695769+00','2024-01-15 19:43:39.695769+00','Test Template');

INSERT INTO "public"."users" ("email","password") VALUES ('test@test.com','\x2432612431302446493930656630382e644c55615a38644f4474677275647139744264744c6d34385a62656e7352794f53375257334d626175564169');

