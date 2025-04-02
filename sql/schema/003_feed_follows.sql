-- +goose Up
CREATE TABLE feed_follows(
	"id" uuid PRIMARY KEY,
	"created_at" timestamp NOT NULL,
	"updated_at" timestamp NOT NULL,
	"feed_id" uuid NOT NULL,
	"user_id" uuid NOT NULL,
	FOREIGN KEY ("user_id") REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY ("feed_id") REFERENCES feeds(id) ON DELETE CASCADE,
	CONSTRAINT unique_user_feed UNIQUE ("feed_id", "user_id")
);


-- +goose Down
DROP TABLE feed_follows;