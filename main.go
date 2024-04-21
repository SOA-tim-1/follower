package main

import (
	"context"
	"fmt"
	"follower/model"
	"follower/repo"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func initializeNeo4j(ctx context.Context) (neo4j.DriverWithContext, error) {
	dbUri := "neo4j+s://15922f8a.databases.neo4j.io"
	dbUser := "neo4j"
	dbPassword := "b_YTsWDnWRARx3s9-Tu0xncPu4bg9ws___PiTURovS8"
	driver, err := neo4j.NewDriverWithContext(
		dbUri,
		neo4j.BasicAuth(dbUser, dbPassword, ""))
	if err != nil {
		return nil, err
	}
	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		driver.Close(ctx)
		return nil, err
	}
	return driver, nil
}

func startServer() {
	router := mux.NewRouter().StrictSlash(true)

	// router.HandleFunc("/user", userHandler.Create).Methods("POST")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	println("Server starting")
	log.Fatal(http.ListenAndServe(":8092",
		handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		)(router)))
}

func New(logger *log.Logger) (*repo.UserRepository, error) {
	uri := "neo4j+s://15922f8a.databases.neo4j.io"
	user := "neo4j"
	pass := "b_YTsWDnWRARx3s9-Tu0xncPu4bg9ws___PiTURovS8"
	auth := neo4j.BasicAuth(user, pass, "")

	driver, err := neo4j.NewDriverWithContext(uri, auth)
	if err != nil {
		logger.Panic(err)
		return nil, err
	}

	// Return repository with logger and DB session
	return &repo.UserRepository{
		Driver: driver,
		Logger: logger,
	}, nil
}

func main() {
	storeLogger := log.New(os.Stdout, "[person-store] ", log.LstdFlags)

	// NoSQL: Initialize Movie Repository store
	store, err := New(storeLogger)
	if err != nil {
		storeLogger.Fatal(err)
	}
	defer store.CloseDriverConnection(context.Background())
	store.CheckConnection()

	var user1 = model.User{ID: 101}
	err = store.WritePerson(&user1)

	var user2 = model.User{ID: 102}
	err = store.WritePerson(&user2)

	var user3 = model.User{ID: 103}
	err = store.WritePerson(&user3)

	var user4 = model.User{ID: 104}
	err = store.WritePerson(&user4)

	var user5 = model.User{ID: 105}
	err = store.WritePerson(&user5)

	var user6 = model.User{ID: 106}
	err = store.WritePerson(&user6)

	var user7 = model.User{ID: 107}
	err = store.WritePerson(&user7)

	user, err := store.FindById(101)
	if err != nil {
		fmt.Println("Error")
	}
	fmt.Println("Created user id is %d", user.ID)

	err = store.CreateFollowConnection(user1.ID, user2.ID)
	err = store.CreateFollowConnection(user1.ID, user3.ID)
	err = store.CreateFollowConnection(user2.ID, user4.ID)
	err = store.CreateFollowConnection(user2.ID, user5.ID)
	err = store.CreateFollowConnection(user3.ID, user7.ID)
	err = store.CreateFollowConnection(user3.ID, user7.ID)
	err = store.GetFollows(user1.ID)
	err = store.GetSuggestionsForUser(user1.ID)

	startServer()

}
