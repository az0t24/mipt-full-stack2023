package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	handler "scoreboardpro/internal/handler/http"
	"scoreboardpro/internal/repository"
	"scoreboardpro/internal/repository/sqlite"
	"scoreboardpro/internal/service"
	"scoreboardpro/pkg/auth"
	"scoreboardpro/pkg/oauth"
	"time"
)

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	//logFile, err := os.OpenFile("log/backend.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	//if err != nil {
	//	log.Fatalf("error opening file: %v\n", err)
	//}
	//defer logFile.Close()
	//mw := io.MultiWriter(os.Stdout, logFile)
	//log.SetOutput(mw)

	db_uri, ok := os.LookupEnv("DB_URI")
	if !ok {
		log.Println("cannot get DB_URI from ENV")
		db_uri = "test.db"
	}

	db, err := sqlite.NewSQLiteDB(db_uri)
	if err != nil {
		log.Panicf("Failed to initialize database: %s\n", err.Error())
	} else {
		log.Println("Database is initialized")
	}

	repo := repository.NewRepository(db)

	serv := service.NewService(repo)

	signingKey, ok := os.LookupEnv("AUTH_SIGNING_KEY")
	if !ok {
		log.Println("cannot get AUTH_SIGNING_KEY from ENV")
		signingKey = "siuefui4nfweu"
	}
	authManager := auth.NewAuthManager([]byte(signingKey))

	// TODO: move these secrets to ENV, then LookupEnv
	client_id := "8772607b1d804d8a9795c7524de39147"
	client_secret := "87d52e5245e74f04b7159e7db24e484f"
	oauthManager := oauth.NewOAuthManager(client_id, client_secret)

	h := handler.NewHandler(serv, authManager, oauthManager)

	srv := &http.Server{
		Addr: ":8082",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      h.NewRouter(), // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)

}
