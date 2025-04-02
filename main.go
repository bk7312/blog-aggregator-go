package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/bk7312/blog-aggregator-go/internal/config"
	"github.com/bk7312/blog-aggregator-go/internal/database"
	_ "github.com/lib/pq"
)

func main() {

	if len(os.Args) < 2 {
		log.Fatal("insufficient arguments provided")
	}

	s := state{
		cfg: new(config.Config),
		db:  new(database.Queries),
	}
	*s.cfg = config.Read()

	db, err := sql.Open("postgres", s.cfg.DbUrl)
	if err != nil {
		log.Fatal("fail to read from db url")
	}

	s.db = database.New(db)

	c := commands{
		handlers: make(map[string]func(*state, command) error),
	}

	c.register("login", handleLogin)
	c.register("register", handleRegister)
	c.register("reset", handleReset)
	c.register("users", handleUsers)
	c.register("agg", handleAgg)
	c.register("addfeed", middlewareLoggedIn(handleAddFeed))
	c.register("feeds", handleFeeds)
	c.register("follow", middlewareLoggedIn(handleFollow))
	c.register("following", middlewareLoggedIn(handleFollowing))
	c.register("unfollow", middlewareLoggedIn(handleUnfollow))
	c.register("browse", middlewareLoggedIn(handleBrowse))
	c.register("help", passCmds(handleHelp, &c))

	// fmt.Println(len(os.Args), os.Args)
	// 2 [/var/.../exe/blog-aggregator-go, cmdname, ...args]

	comm := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	err = c.run(&s, comm)
	if err != nil {
		log.Fatal(err)
	}

}
