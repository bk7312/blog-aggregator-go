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

And you should be able to run `go run . <command>` to use the CLI.

Alternatively, you can install using:

```bash
go install github.com/bk7312/blog-aggregator-go
```

And run directly using `blog-aggregator-go <command>` instead.

## Usage

To use the CLI, run `go run . <command>`, or `blog-aggregator-go <command>` after installing. Refer below for the list of available commands.

- login: `login <username>` to login once registered.
- register: `register <username>` to register.
- reset: `reset` will reset the database.
- users: `users` to list all users.
- agg: : `agg <time>` to fetch the feeds, runs continuously in intervals set by `<time>` (5s, 3m, 1h, etc).
- addfeed: `addfeed <feed_name> <feed_url>` to add feed to db.
- feeds: `feeds` to list all available feeds in db.
- follow: `follow <feed_url>` to follow the feed.
- following: `following` to show list of followed feeds.
- unfollow: `unfollow <feed_url>` to unfollow the feed.
- browse: `browse <limit>` to show latest `<limit>` number of posts.
- help: `help` to show the list of available commands.
