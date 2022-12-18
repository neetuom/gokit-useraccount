package main

import (
	"context"
	"database/sql"

	//PostgreSQL library
	_ "github.com/lib/pq"

	// MySQL library
     _ "github.com/go-sql-driver/mysql"
	"flag"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"net/http"
	"os"
	"os/signal"
	"syscall"

	"gokit-useraccount.com/account"
)



//Syntax 1 - PostgreSQL database
// const dbsource = "postgresql://postgres:postgres@localhost:5432/gokitexample?sslmode=disable"

//Syntax 1 - MySQl database
// const dbsource = "root:A******@10@tcp(localhost:3306)/golangdb"



// Syntax 2 - part A
const (
    username = "root"
    password = "A******@10"
    hostname = "localhost:3306"
    dbname   = "golangdb"
)

// Syntax 2 - part B
//a small function that will return us this DSN when the database name is passed as a parameter.
func dsn(dbName string) string {
     return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}



func main() {

     var httpAddr = flag.String("http",":8088","http listen address")

     var logger log.Logger
     	{
     		logger = log.NewLogfmtLogger(os.Stderr)
     		logger = log.NewSyncLogger(logger)
     		logger = log.With(logger,
     			"service", "account",
     			"time:", log.DefaultTimestampUTC,
     			"caller", log.DefaultCaller,
     		)
     	}

     level.Info(logger).Log("msg", "service started")

     defer level.Info(logger).Log("msg", "service ended")


     //Syntax 1 - To open DB connection
     /*
     var db *sql.DB
     	{
     		var err error

     		db, err = sql.Open("mysql", dbsource)  // db, err = sql.Open("postgres", dbsource)   - PostgreSQL

     		if err != nil {
     			level.Error(logger).Log("exit", err)
     			os.Exit(-1)
     		}
     	}
     */


     //Syntax 2 - To open DB connection
        db, err := sql.Open("mysql", dsn(dbname))
        if err != nil {
             			level.Error(logger).Log("exit", err)
             			os.Exit(-1)
        }


     	flag.Parse()

     	ctx := context.Background()


     	var srv account.Service
     	{
     		repository := account.NewRepo(db, logger)
     		srv = account.NewService(repository, logger)
     	}


     	errs := make(chan error)

     	go func() {
     		c := make(chan os.Signal, 1)
     		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
     		errs <- fmt.Errorf("%s", <-c)
     	}()


     	endpoints := account.MakeEndpoints(srv)

     	go func() {
     		fmt.Println("listening on port", *httpAddr)
     		handler := account.NewHTTPServer(ctx, endpoints)
     		errs <- http.ListenAndServe(*httpAddr, handler)
     	}()


        // defer the close till after the main function has finished executing
        defer db.Close()

     	level.Error(logger).Log("exit", <-errs)
}