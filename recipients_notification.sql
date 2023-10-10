-- Adminer 4.8.1 PostgreSQL 16.0 (Debian 16.0-1.pgdg120+1) dump

DROP TABLE IF EXISTS "notification_contents";
CREATE TABLE "public"."notification_contents" (
    "recipient_id" bigint NOT NULL,
    "notification_id" bigint NOT NULL,
    "message_content" character varying(100) NOT NULL,
    CONSTRAINT "notification_contents_pkey" PRIMARY KEY ("recipient_id", "notification_id")
) WITH (oids = false);

INSERT INTO "notification_contents" ("recipient_id", "notification_id", "message_content") VALUES
(1,	1,	'Срочное уведомление о задолжности');

DROP TABLE IF EXISTS "notifications";
DROP SEQUENCE IF EXISTS notifications_notification_id_seq;
CREATE SEQUENCE notifications_notification_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1;

CREATE TABLE "public"."notifications" (
    "notification_id" bigint DEFAULT nextval('notifications_notification_id_seq') NOT NULL,
    "status" character varying(50) NOT NULL,
    "creation_date" date NOT NULL,
    "formation_date" date,
    "completion_date" date,
    "moderator_id" bigint NOT NULL,
    "customer_id" bigint NOT NULL,
    "notification_type" character varying(50) NOT NULL,
    CONSTRAINT "notifications_pkey" PRIMARY KEY ("notification_id")
) WITH (oids = false);

INSERT INTO "notifications" ("notification_id", "status", "creation_date", "formation_date", "completion_date", "moderator_id", "customer_id", "notification_type") VALUES
(1,	'Введен',	'2023-10-09',	'2023-11-09',	NULL,	2,	1,	'Срочное сообщение');

DROP TABLE IF EXISTS "recipients";
CREATE TABLE "public"."recipients" (
    "recipient_id" bigint NOT NULL,
    "fio" character varying(100) NOT NULL,
    "image_url" character varying(100) NOT NULL,
    "email" character varying(50) NOT NULL,
    "age" bigint NOT NULL,
    "adress" character varying(100) NOT NULL,
    "delivered" boolean,
    CONSTRAINT "recipients_pkey" PRIMARY KEY ("recipient_id")
) WITH (oids = false);

INSERT INTO "recipients" ("recipient_id", "fio", "image_url", "email", "age", "adress", "delivered") VALUES
(1,	'Олег Орлов Никитович',	'http://localhost:8080/image/men1.jpg',	'OlegO@mail.ru',	27,	'Москва, ул. Измайловская, д.13, кв.54',	NULL);

DROP TABLE IF EXISTS "users";
DROP SEQUENCE IF EXISTS users_user_id_seq;
CREATE SEQUENCE users_user_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1;

CREATE TABLE "public"."users" (
    "user_id" bigint DEFAULT nextval('users_user_id_seq') NOT NULL,
    "login" character varying(30) NOT NULL,
    "password" character varying(30) NOT NULL,
    "name" character varying(50) NOT NULL,
    "moderator" boolean NOT NULL,
    CONSTRAINT "users_pkey" PRIMARY KEY ("user_id")
) WITH (oids = false);

INSERT INTO "users" ("user_id", "login", "password", "name", "moderator") VALUES
(1,	'user 1',	'pass 1',	'Пользователь',	'f'),
(2,	'user 2',	'pass 2',	'Модератор',	't');

ALTER TABLE ONLY "public"."notification_contents" ADD CONSTRAINT "fk_notification_contents_notification" FOREIGN KEY (notification_id) REFERENCES notifications(notification_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."notification_contents" ADD CONSTRAINT "fk_notification_contents_recipient" FOREIGN KEY (recipient_id) REFERENCES recipients(recipient_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."notifications" ADD CONSTRAINT "fk_notifications_customer" FOREIGN KEY (customer_id) REFERENCES users(user_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."notifications" ADD CONSTRAINT "fk_notifications_moderator" FOREIGN KEY (moderator_id) REFERENCES users(user_id) NOT DEFERRABLE;

-- 2023-10-10 14:29:23.077514+00
