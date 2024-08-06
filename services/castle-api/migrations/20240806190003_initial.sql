-- Create "files" table
CREATE TABLE "files" ("id" character varying NOT NULL, "name" character varying NOT NULL, "file_name" character varying NOT NULL, "mime_type" character varying NOT NULL, "disk" character varying NOT NULL, "size" bigint NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, PRIMARY KEY ("id"));
-- Create "users" table
CREATE TABLE "users" ("id" character varying NOT NULL, "first_name" character varying NOT NULL, "last_name" character varying NOT NULL, "nickname" character varying NOT NULL, "email" character varying NOT NULL, "hashed_password" bytea NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, PRIMARY KEY ("id"));
-- Create index "users_email_key" to table: "users"
CREATE UNIQUE INDEX "users_email_key" ON "users" ("email");
-- Create index "users_nickname_key" to table: "users"
CREATE UNIQUE INDEX "users_nickname_key" ON "users" ("nickname");
-- Create "friendships" table
CREATE TABLE "friendships" ("id" character varying NOT NULL, "created_at" timestamptz NOT NULL, "user_id" character varying NOT NULL, "friend_id" character varying NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "friendships_users_friend" FOREIGN KEY ("friend_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT "friendships_users_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- Create index "friendship_user_id_friend_id" to table: "friendships"
CREATE UNIQUE INDEX "friendship_user_id_friend_id" ON "friendships" ("user_id", "friend_id");
-- Create "rooms" table
CREATE TABLE "rooms" ("id" character varying NOT NULL, "name" character varying NOT NULL, "version" bigint NOT NULL DEFAULT 0, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, PRIMARY KEY ("id"));
-- Create "messages" table
CREATE TABLE "messages" ("id" character varying NOT NULL, "content" text NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, "room_messages" character varying NULL, "user_messages" character varying NULL, PRIMARY KEY ("id"), CONSTRAINT "messages_rooms_messages" FOREIGN KEY ("room_messages") REFERENCES "rooms" ("id") ON UPDATE NO ACTION ON DELETE SET NULL, CONSTRAINT "messages_users_messages" FOREIGN KEY ("user_messages") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE SET NULL);
-- Create "room_members" table
CREATE TABLE "room_members" ("id" character varying NOT NULL, "joined_at" timestamptz NOT NULL, "user_id" character varying NOT NULL, "room_id" character varying NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "room_members_rooms_room" FOREIGN KEY ("room_id") REFERENCES "rooms" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "room_members_users_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create index "roommember_user_id_room_id" to table: "room_members"
CREATE UNIQUE INDEX "roommember_user_id_room_id" ON "room_members" ("user_id", "room_id");
