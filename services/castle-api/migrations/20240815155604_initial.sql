-- Create "rooms" table
CREATE TABLE "rooms" (
  "id" character varying NOT NULL,
  "name" character varying NOT NULL,
  "version" bigint NOT NULL DEFAULT 1,
  "type" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "users" table
CREATE TABLE "users" (
  "id" character varying NOT NULL,
  "first_name" character varying NOT NULL,
  "last_name" character varying NOT NULL,
  "nickname" character varying NOT NULL,
  "email" character varying NOT NULL,
  "password" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "users_email_key" to table: "users"
CREATE UNIQUE INDEX "users_email_key" ON "users" ("email");
-- Create index "users_nickname_key" to table: "users"
CREATE UNIQUE INDEX "users_nickname_key" ON "users" ("nickname");
-- Create "friendships" table
CREATE TABLE "friendships" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "user_id" character varying NOT NULL,
  "friend_id" character varying NOT NULL,
  "room_id" character varying NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "friendships_rooms_room" FOREIGN KEY ("room_id") REFERENCES "rooms" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "friendships_users_friend" FOREIGN KEY ("friend_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "friendships_users_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "friendship_user_id_friend_id" to table: "friendships"
CREATE UNIQUE INDEX "friendship_user_id_friend_id" ON "friendships" ("user_id", "friend_id");
-- Create "files" table
CREATE TABLE "files" (
  "id" character varying NOT NULL,
  "name" character varying NOT NULL,
  "mime_type" character varying NOT NULL,
  "disk" character varying NOT NULL,
  "size" bigint NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "messages" table
CREATE TABLE "messages" (
  "id" character varying NOT NULL,
  "content" text NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "message_reply_to" character varying NULL,
  "room_messages" character varying NULL,
  "user_messages" character varying NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "messages_messages_reply_to" FOREIGN KEY ("message_reply_to") REFERENCES "messages" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "messages_rooms_messages" FOREIGN KEY ("room_messages") REFERENCES "rooms" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "messages_users_messages" FOREIGN KEY ("user_messages") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Create index "messages_message_reply_to_key" to table: "messages"
CREATE UNIQUE INDEX "messages_message_reply_to_key" ON "messages" ("message_reply_to");
-- Create "message_attachments" table
CREATE TABLE "message_attachments" (
  "id" character varying NOT NULL,
  "type" character varying NOT NULL,
  "order" bigint NOT NULL,
  "attached_at" timestamptz NOT NULL,
  "file_message_attachment" character varying NOT NULL,
  "message_attachments" character varying NOT NULL,
  "room_message_attachments" character varying NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "message_attachments_files_message_attachment" FOREIGN KEY ("file_message_attachment") REFERENCES "files" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "message_attachments_messages_attachments" FOREIGN KEY ("message_attachments") REFERENCES "messages" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "message_attachments_rooms_message_attachments" FOREIGN KEY ("room_message_attachments") REFERENCES "rooms" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "message_attachments_file_message_attachment_key" to table: "message_attachments"
CREATE UNIQUE INDEX "message_attachments_file_message_attachment_key" ON "message_attachments" ("file_message_attachment");
-- Create "message_links" table
CREATE TABLE "message_links" (
  "id" character varying NOT NULL,
  "url" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "message_links" character varying NOT NULL,
  "room_message_links" character varying NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "message_links_messages_links" FOREIGN KEY ("message_links") REFERENCES "messages" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "message_links_rooms_message_links" FOREIGN KEY ("room_message_links") REFERENCES "rooms" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "message_voices" table
CREATE TABLE "message_voices" (
  "id" character varying NOT NULL,
  "length" bigint NOT NULL,
  "attached_at" timestamptz NOT NULL,
  "file_message_voice" character varying NOT NULL,
  "message_voice" character varying NOT NULL,
  "room_message_voices" character varying NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "message_voices_files_message_voice" FOREIGN KEY ("file_message_voice") REFERENCES "files" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "message_voices_messages_voice" FOREIGN KEY ("message_voice") REFERENCES "messages" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "message_voices_rooms_message_voices" FOREIGN KEY ("room_message_voices") REFERENCES "rooms" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "message_voices_file_message_voice_key" to table: "message_voices"
CREATE UNIQUE INDEX "message_voices_file_message_voice_key" ON "message_voices" ("file_message_voice");
-- Create index "message_voices_message_voice_key" to table: "message_voices"
CREATE UNIQUE INDEX "message_voices_message_voice_key" ON "message_voices" ("message_voice");
-- Create "room_members" table
CREATE TABLE "room_members" (
  "id" character varying NOT NULL,
  "joined_at" timestamptz NOT NULL,
  "user_id" character varying NOT NULL,
  "room_id" character varying NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "room_members_rooms_room" FOREIGN KEY ("room_id") REFERENCES "rooms" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "room_members_users_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "roommember_room_id_user_id" to table: "room_members"
CREATE UNIQUE INDEX "roommember_room_id_user_id" ON "room_members" ("room_id", "user_id");
