# gator

A CLI blog aggregator from [boot.dev](https://www.boot.dev/courses/build-blog-aggregator-golang) built using Go.

## Installation

Ensure go is installed and execute the following command in your terminal:

```bash
go mod download
```

Ensure postgres is installed and create a `gator` database.

```sql
CREATE DATABASE gator;
```

Ensure `goose` is installed, then navigate to `sql/schema` and migrate the database.

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
goose postgres <postgres_connection_string> up
```

Create a json file in your home directory `~/.gatorconfig.json` containing the following your postgres connection url. Username can be left blank as you'll need to `register` the user in the db first later.

```json
{
  "db_url": "connection_string_goes_here",
  "current_user_name": "username_goes_here"
}
```

And you should be able to run `go run . <command>` to use the CLI. Alternatively, you can run `go build` to compile the binary and run using that instead, i.e. `<binary_name> <command>`

## Usage

To use the CLI, run `go run . <command>`, or `<binary_name> <command>` after compiling. Refer below for the list of available commands.

- login: `go run . login <username>` to login once registered.
- register: `go run . register <username>` to register.
- reset: `go run . reset` will reset the database.
- users: `go run . users` to list all users.
- agg: : `go run . agg <time>` to fetch the feeds, runs continuously in intervals set by `<time>` (5s, 3m, 1h, etc).
- addfeed: `go run . addfeed <feed_name> <feed_url>` to add feed to db.
- feeds: `go run . feeds` to list all available feeds in db.
- follow: `go run . follow <feed_url>` to follow the feed.
- following: `go run . following` to show list of followed feeds.
- unfollow: `go run . unfollow <feed_url>` to unfollow the feed.
- browse: `go run . browse <limit>` to show latest `<limit>` number of posts.
- help: `go run . help` to show the list of available commands.
