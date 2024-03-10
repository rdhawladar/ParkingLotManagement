-- -------------------------------------------------------------
-- TablePlus 5.9.0(538)
--
-- https://tableplus.com/
--
-- Database: parkingapp
-- Generation Time: 2024-03-10 15:26:13.6670
-- -------------------------------------------------------------


DROP TABLE IF EXISTS "public"."maintenance";
-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS maintenance_id_seq;

-- Table Definition
CREATE TABLE "public"."maintenance" (
    "id" int4 NOT NULL DEFAULT nextval('maintenance_id_seq'::regclass),
    "parking_lot_id" int4 NOT NULL,
    "parking_slot_id" int4,
    "start_time" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "end_time" timestamp,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "public"."parkinglots";
-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS parkinglots_id_seq;

-- Table Definition
CREATE TABLE "public"."parkinglots" (
    "id" int4 NOT NULL DEFAULT nextval('parkinglots_id_seq'::regclass),
    "name" varchar(255) NOT NULL,
    "total_spaces" int4 NOT NULL,
    "manager_id" int4,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "public"."parkingsessions";
-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS parkingsessions_id_seq;

-- Table Definition
CREATE TABLE "public"."parkingsessions" (
    "id" int4 NOT NULL DEFAULT nextval('parkingsessions_id_seq'::regclass),
    "vehicle_id" int4 NOT NULL,
    "parking_lot_id" int4 NOT NULL,
    "parking_slot_id" int4,
    "parked_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "unparked_at" timestamp,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "public"."parkingslots";
-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS parkingslots_id_seq;

-- Table Definition
CREATE TABLE "public"."parkingslots" (
    "id" int4 NOT NULL DEFAULT nextval('parkingslots_id_seq'::regclass),
    "parking_lot_id" int4 NOT NULL,
    "slot_number" int4 NOT NULL,
    "is_available" bool DEFAULT true,
    "is_under_maintenance" bool DEFAULT false,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "public"."schema_migrations";
-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Table Definition
CREATE TABLE "public"."schema_migrations" (
    "version" int8 NOT NULL,
    "dirty" bool NOT NULL,
    PRIMARY KEY ("version")
);

DROP TABLE IF EXISTS "public"."users";
-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS users_id_seq;
DROP TYPE IF EXISTS "public"."user_role";
CREATE TYPE "public"."user_role" AS ENUM ('user', 'manager');

-- Table Definition
CREATE TABLE "public"."users" (
    "id" int4 NOT NULL DEFAULT nextval('users_id_seq'::regclass),
    "name" varchar(255) NOT NULL,
    "role" "public"."user_role" NOT NULL DEFAULT 'user'::user_role,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "public"."vehicles";
-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS vehicles_id_seq;

-- Table Definition
CREATE TABLE "public"."vehicles" (
    "id" int4 NOT NULL DEFAULT nextval('vehicles_id_seq'::regclass),
    "license_plate" varchar(255) NOT NULL,
    "owner_id" int4,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);

INSERT INTO "public"."parkinglots" ("id", "name", "total_spaces", "manager_id", "created_at", "updated_at") VALUES
(8, 'LOT 1', 50, 1, '2024-03-10 09:09:58.876787', '2024-03-10 09:09:58.876787');

INSERT INTO "public"."parkingsessions" ("id", "vehicle_id", "parking_lot_id", "parking_slot_id", "parked_at", "unparked_at", "created_at", "updated_at") VALUES
(4, 1, 8, 1, '2024-03-10 07:26:34.287646', '2024-03-10 09:52:29.377818', '2024-03-10 09:26:34.287646', '2024-03-10 09:26:34.287646'),
(5, 1, 8, 1, '2024-03-10 09:58:06.522583', '2024-03-10 09:58:55.518989', '2024-03-10 09:58:06.522583', '2024-03-10 09:58:06.522583'),
(6, 1, 8, 1, '2024-03-10 08:04:14.247615', NULL, '2024-03-10 10:04:14.247615', '2024-03-10 10:04:14.247615'),
(9, 1, 8, 2, '2024-03-10 11:37:51.502448', NULL, '2024-03-10 11:37:51.502448', '2024-03-10 11:37:51.502448'),
(10, 1, 8, 3, '2024-03-10 11:37:54.785449', NULL, '2024-03-10 11:37:54.785449', '2024-03-10 11:37:54.785449');

INSERT INTO "public"."parkingslots" ("id", "parking_lot_id", "slot_number", "is_available", "is_under_maintenance", "created_at", "updated_at") VALUES
(1, 8, 1, 'f', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(2, 8, 2, 'f', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(3, 8, 3, 'f', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(4, 8, 4, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(5, 8, 5, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(6, 8, 6, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(7, 8, 7, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(8, 8, 8, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(9, 8, 9, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(10, 8, 10, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(11, 8, 11, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(12, 8, 12, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(13, 8, 13, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(14, 8, 14, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(15, 8, 15, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(16, 8, 16, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(17, 8, 17, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(18, 8, 18, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(19, 8, 19, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(20, 8, 20, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(21, 8, 21, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(22, 8, 22, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(23, 8, 23, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(24, 8, 24, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(25, 8, 25, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(26, 8, 26, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(27, 8, 27, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(28, 8, 28, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(29, 8, 29, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(30, 8, 30, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(31, 8, 31, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(32, 8, 32, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(33, 8, 33, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(34, 8, 34, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(35, 8, 35, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(36, 8, 36, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(37, 8, 37, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(38, 8, 38, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(39, 8, 39, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(40, 8, 40, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(41, 8, 41, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(42, 8, 42, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(43, 8, 43, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(44, 8, 44, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(45, 8, 45, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(46, 8, 46, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(47, 8, 47, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(48, 8, 48, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(49, 8, 49, 't', 'f', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793'),
(50, 8, 50, 't', 't', '2024-03-10 09:09:58.878793', '2024-03-10 09:09:58.878793');

INSERT INTO "public"."schema_migrations" ("version", "dirty") VALUES
(1, 'f');

INSERT INTO "public"."users" ("id", "name", "role", "created_at", "updated_at") VALUES
(1, 'Parking Manager', 'manager', '2024-03-10 07:46:27.182757', '2024-03-10 07:46:27.182757'),
(2, 'User', 'user', '2024-03-10 07:47:01.30247', '2024-03-10 07:47:01.30247');

INSERT INTO "public"."vehicles" ("id", "license_plate", "owner_id", "created_at", "updated_at") VALUES
(1, '12345', 2, '2024-03-10 09:25:09.917158', '2024-03-10 09:25:09.917158');

ALTER TABLE "public"."maintenance" ADD FOREIGN KEY ("parking_lot_id") REFERENCES "public"."parkinglots"("id") ON DELETE CASCADE;
ALTER TABLE "public"."maintenance" ADD FOREIGN KEY ("parking_slot_id") REFERENCES "public"."parkingslots"("id") ON DELETE CASCADE;
ALTER TABLE "public"."parkinglots" ADD FOREIGN KEY ("manager_id") REFERENCES "public"."users"("id") ON DELETE SET NULL;
ALTER TABLE "public"."parkingsessions" ADD FOREIGN KEY ("vehicle_id") REFERENCES "public"."vehicles"("id") ON DELETE CASCADE;
ALTER TABLE "public"."parkingsessions" ADD FOREIGN KEY ("parking_slot_id") REFERENCES "public"."parkingslots"("id") ON DELETE SET NULL;
ALTER TABLE "public"."parkingsessions" ADD FOREIGN KEY ("parking_lot_id") REFERENCES "public"."parkinglots"("id") ON DELETE CASCADE;
ALTER TABLE "public"."parkingslots" ADD FOREIGN KEY ("parking_lot_id") REFERENCES "public"."parkinglots"("id") ON DELETE CASCADE;
ALTER TABLE "public"."vehicles" ADD FOREIGN KEY ("owner_id") REFERENCES "public"."users"("id") ON DELETE CASCADE;
