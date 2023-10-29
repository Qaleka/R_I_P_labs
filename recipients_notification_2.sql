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
('18ab9f76-7648-49d2-857d-75ffddf13bea',	'02f5e742-0221-484d-a796-c74f4691e693');

DROP TABLE IF EXISTS "notifications";
CREATE TABLE "public"."notifications" (
    "uuid" uuid DEFAULT gen_random_uuid() NOT NULL,
    "status" character varying(20) NOT NULL,
    "creation_date" date NOT NULL,
    "formation_date" date,
    "completion_date" date,
    "moderator_id" uuid,
    "customer_id" uuid NOT NULL,
    "notification_type" character varying(50) NOT NULL,
    CONSTRAINT "notifications_pkey" PRIMARY KEY ("uuid")
) WITH (oids = false);

INSERT INTO "notifications" ("uuid", "status", "creation_date", "formation_date", "completion_date", "moderator_id", "customer_id", "notification_type") VALUES
('c4fdc129-ed48-48df-a262-6be92a3acb12',	'Введен',	'2023-10-09',	'2023-11-09',	'2023-12-09',	'87d54d58-1e24-4cca-9c83-bd2523902729',	'2d217868-ab6d-41fe-9b34-7809083a2e8a',	'Срочное сообщение'),
('3c0a95f9-2e96-421b-903d-3bb9262cd77e',	'В работе',	'2023-09-24',	'2023-09-27',	'2023-09-28',	'87d54d58-1e24-4cca-9c83-bd2523902729',	'2d217868-ab6d-41fe-9b34-7809083a2e8a',	'Еженедельное уведомление'),
('02f5e742-0221-484d-a796-c74f4691e693',	'Завершен',	'2023-10-09',	'2023-10-15',	'2023-10-20',	'87d54d58-1e24-4cca-9c83-bd2523902729',	'2d217868-ab6d-41fe-9b34-7809083a2e8a',	'Электронное напоминание'),
('1565f159-b6a2-4108-bc29-456dd05a8ac4',	'Отменен',	'2023-09-15',	'2023-09-15',	'2023-09-16',	'87d54d58-1e24-4cca-9c83-bd2523902729',	'2d217868-ab6d-41fe-9b34-7809083a2e8a',	'Уведомление о задолжности'),
('03d54ad2-0b82-4d49-9ca0-e67add804f4d',	'Удален',	'2023-10-08',	'2023-10-09',	'2023-10-10',	'87d54d58-1e24-4cca-9c83-bd2523902729',	'2d217868-ab6d-41fe-9b34-7809083a2e8a',	'Уведомление о доступе на сайт'),
('6899a458-f86f-4cad-95ee-43c0e646d0ec',	'черновик',	'2023-10-29',	NULL,	NULL,	NULL,	'2d217868-ab6d-41fe-9b34-7809083a2e8a',	'');

DROP TABLE IF EXISTS "recipients";
CREATE TABLE "public"."recipients" (
    "uuid" uuid DEFAULT gen_random_uuid() NOT NULL,
    "fio" character varying(100) NOT NULL,
    "image_url" character varying(100) NOT NULL,
    "email" character varying(75) NOT NULL,
    "age" bigint NOT NULL,
    "adress" character varying(100) NOT NULL,
    "is_deleted" boolean DEFAULT false NOT NULL,
    CONSTRAINT "recipients_pkey" PRIMARY KEY ("uuid")
) WITH (oids = false);

INSERT INTO "recipients" ("uuid", "fio", "image_url", "email", "age", "adress", "is_deleted") VALUES
('4bea0842-bcb8-416e-9a63-d89a63e978ca',	'Олег Орлов Никитович',	'http://localhost:9000/image/men1.jpg',	'OlegO@mail.ru',	27,	'Москва, ул. Измайловская, д.13, кв.54',	't'),
('b9778018-9c13-46fd-b785-4a803dc8be0b',	'Александр Лейко Кириллович',	'http://localhost:9000/image/men3.jpg',	'Alek221@mail.ru',	37,	'Москва, ул. Изюмская, д.15, кв.89',	'f'),
('18ab9f76-7648-49d2-857d-75ffddf13bea',	'Василий Гречко Валентинович',	'http://localhost:9000/image/men2.jpg',	'Grechko_101@mail.ru',	31,	'Москва, ул. Тверская, д.25, кв.145',	'f');

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

-- 2023-10-29 18:49:26.064055+00
