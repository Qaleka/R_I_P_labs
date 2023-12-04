-- Adminer 4.8.1 PostgreSQL 16.0 (Debian 16.0-1.pgdg120+1) dump

DROP TABLE IF EXISTS "notification_contents";
CREATE TABLE "public"."notification_contents" (
    "recipient_id" uuid DEFAULT gen_random_uuid() NOT NULL,
    "notification_id" uuid DEFAULT gen_random_uuid() NOT NULL,
    CONSTRAINT "notification_contents_pkey" PRIMARY KEY ("recipient_id", "notification_id")
) WITH (oids = false);

INSERT INTO "notification_contents" ("recipient_id", "notification_id") VALUES
('4bea0842-bcb8-416e-9a63-d89a63e978ca',	'c4fdc129-ed48-48df-a262-6be92a3acb12'),
('b9778018-9c13-46fd-b785-4a803dc8be0b',	'1565f159-b6a2-4108-bc29-456dd05a8ac4'),
('18ab9f76-7648-49d2-857d-75ffddf13bea',	'02f5e742-0221-484d-a796-c74f4691e693'),
('b9778018-9c13-46fd-b785-4a803dc8be0b',	'200b2366-36d5-49b2-9770-85c1628c20f0'),
('b9778018-9c13-46fd-b785-4a803dc8be0b',	'cf7bf391-a53c-4995-ab17-82acca3a6bd1');

DROP TABLE IF EXISTS "notifications";
CREATE TABLE "public"."notifications" (
    "uuid" uuid DEFAULT gen_random_uuid() NOT NULL,
    "status" character varying(20) NOT NULL,
    "creation_date" timestamp NOT NULL,
    "formation_date" timestamp,
    "completion_date" timestamp,
    "moderator_id" uuid,
    "customer_id" uuid NOT NULL,
    "notification_type" character varying(50) NOT NULL,
    CONSTRAINT "notifications_pkey" PRIMARY KEY ("uuid")
) WITH (oids = false);

INSERT INTO "notifications" ("uuid", "status", "creation_date", "formation_date", "completion_date", "moderator_id", "customer_id", "notification_type") VALUES
('c4fdc129-ed48-48df-a262-6be92a3acb12',	'Введен',	'2023-10-09 00:00:00',	'2023-11-09 00:00:00',	'2023-12-09 00:00:00',	'87d54d58-1e24-4cca-9c83-bd2523902729',	'2d217868-ab6d-41fe-9b34-7809083a2e8a',	'Срочное сообщение'),
('3c0a95f9-2e96-421b-903d-3bb9262cd77e',	'В работе',	'2023-09-24 00:00:00',	'2023-09-27 00:00:00',	'2023-09-28 00:00:00',	'87d54d58-1e24-4cca-9c83-bd2523902729',	'2d217868-ab6d-41fe-9b34-7809083a2e8a',	'Еженедельное уведомление'),
('a10adc9b-4957-4769-a029-71b76c764e49',	'удалён',	'2023-10-29 00:00:00',	NULL,	NULL,	NULL,	'2d217868-ab6d-41fe-9b34-7809083a2e8a',	''),
('1565f159-b6a2-4108-bc29-456dd05a8ac4',	'отклонён',	'2023-09-15 00:00:00',	'2023-09-15 00:00:00',	'2023-09-16 00:00:00',	'87d54d58-1e24-4cca-9c83-bd2523902729',	'2d217868-ab6d-41fe-9b34-7809083a2e8a',	'Уведомление о задолжности'),
('02f5e742-0221-484d-a796-c74f4691e693',	'завершён',	'2023-10-09 00:00:00',	'2023-10-15 00:00:00',	'2023-10-20 00:00:00',	'87d54d58-1e24-4cca-9c83-bd2523902729',	'2d217868-ab6d-41fe-9b34-7809083a2e8a',	'Электронное напоминание'),
('200b2366-36d5-49b2-9770-85c1628c20f0',	'сформирован',	'2023-10-30 00:00:00',	'2023-10-30 00:00:00',	NULL,	NULL,	'2d217868-ab6d-41fe-9b34-7809083a2e8a',	''),
('cf7bf391-a53c-4995-ab17-82acca3a6bd1',	'черновик',	'2023-11-02 00:00:00',	NULL,	NULL,	NULL,	'2d217868-ab6d-41fe-9b34-7809083a2e8a',	'');

DROP TABLE IF EXISTS "recipients";
CREATE TABLE "public"."recipients" (
    "uuid" uuid DEFAULT gen_random_uuid() NOT NULL,
    "fio" character varying(100) NOT NULL,
    "image_url" character varying(100),
    "email" character varying(75) NOT NULL,
    "age" bigint NOT NULL,
    "adress" character varying(100) NOT NULL,
    "is_deleted" boolean DEFAULT false NOT NULL,
    CONSTRAINT "recipients_pkey" PRIMARY KEY ("uuid")
) WITH (oids = false);

INSERT INTO "recipients" ("uuid", "fio", "image_url", "email", "age", "adress", "is_deleted") VALUES
('4bea0842-bcb8-416e-9a63-d89a63e978ca',	'Олег Орлов Никитович',	'localhost:9000/images/men1.jpg',	'OlegO@mail.ru',	27,	'Москва, ул. Измайловская, д.13, кв.54',	't'),
('18ab9f76-7648-49d2-857d-75ffddf13bea',	'Василий Гречко Валентинович',	'localhost:9000/images/men2.jpg',	'Grechko_101@mail.ru',	31,	'Москва, ул. Тверская, д.25, кв.145',	'f'),
('b9778018-9c13-46fd-b785-4a803dc8be0b',	'Александр Лейко Кириллович',	'localhost:9000/images/men3.jpg',	'Alek221@mail.ru',	37,	'Москва, ул. Изюмская, д.15, кв.89',	'f'),
('365f46c8-b498-47b9-92d3-97319ff62711',	'Андрей Отрис Даниллович',	'localhost:9000/images/men1.jpg',	'Andr1@mail.ru',	32,	'Москва, ул. Изюмская, д.15, кв.79',	'f'),
('9b8914a6-c599-450d-893d-b8ebb766dd07',	'Кирилл Лейка Кириллович',	'localhost:9000/images/men4.jpg',	'KriLeik@mail.ru',	30,	'Москва, ул. Бутовская, д.15, кв.79',	'f'),
('b06a0b5a-6ede-4636-a97b-83976dc10575',	'Андрей Ермолин Данилович',	'localhost:9000/images/men5.jpg',	'andrErm@gmail.com',	33,	'Москва, ул. Клинская, д.22, кв.12',	'f');

DROP TABLE IF EXISTS "users";
CREATE TABLE "public"."users" (
    "uuid" uuid DEFAULT gen_random_uuid() NOT NULL,
    "login" character varying(30) NOT NULL,
    "password" character varying(30) NOT NULL,
    "name" character varying(50) NOT NULL,
    "moderator" boolean NOT NULL,
    CONSTRAINT "users_pkey" PRIMARY KEY ("uuid")
) WITH (oids = false);

INSERT INTO "users" ("uuid", "login", "password", "name", "moderator") VALUES
('2d217868-ab6d-41fe-9b34-7809083a2e8a',	'user 1',	'pass 1',	'Пользователь',	'f'),
('87d54d58-1e24-4cca-9c83-bd2523902729',	'user 2',	'pass 2',	'Модератор',	't');

ALTER TABLE ONLY "public"."notification_contents" ADD CONSTRAINT "fk_notification_contents_notification" FOREIGN KEY (notification_id) REFERENCES notifications(uuid) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."notification_contents" ADD CONSTRAINT "fk_notification_contents_recipient" FOREIGN KEY (recipient_id) REFERENCES recipients(uuid) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."notifications" ADD CONSTRAINT "fk_notifications_customer" FOREIGN KEY (customer_id) REFERENCES users(uuid) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."notifications" ADD CONSTRAINT "fk_notifications_moderator" FOREIGN KEY (moderator_id) REFERENCES users(uuid) NOT DEFERRABLE;

-- 2023-12-04 21:08:38.415951+00
