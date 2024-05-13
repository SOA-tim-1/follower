package main

import (
	"context"
	"fmt"
	"follower/handler"
	"follower/model"
	follower "follower/proto"
	"follower/repo"
	"follower/service"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

func startServer(userHandler *handler.UserHandler) {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api/follower/get-user/{id}", userHandler.GetUserById).Methods("GET")
	router.HandleFunc("/api/follower/get-follows/{id}", userHandler.GetFollows).Methods("GET")
	router.HandleFunc("/api/follower/get-followers/{id}", userHandler.GetFollowers).Methods("GET")
	router.HandleFunc("/api/follower/get-suggestions/{id}", userHandler.GetSuggestionsForUser).Methods("GET")
	router.HandleFunc("/api/follower/follow-connection", userHandler.CreateFollowConnection).Methods("POST")
	router.HandleFunc("/api/follower/delete-follow-connection", userHandler.DeleteFollowConnection).Methods("DELETE")
	router.HandleFunc("/api/follower/check-following", userHandler.CheckIfFirstFollowSecond).Methods("GET")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	println("Server starting")
	log.Fatal(http.ListenAndServe(":8092",
		handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		)(router)))
}

func startServerGRPC(userHandlerGRPC *handler.UserHandlergRPC) {
	listener, err := net.Listen("tcp", ":8092")
	if err != nil {
		log.Fatalln(err)
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(listener)

	// Bootstrap gRPC server.
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	// Bootstrap gRPC service server and respond to request.
	//userHandler := handlers.UserHandler{}
	follower.RegisterUserServiceServer(grpcServer, userHandlerGRPC)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal("server error: ", err)
		}
	}()

	stopCh := make(chan os.Signal)
	signal.Notify(stopCh, syscall.SIGTERM)

	<-stopCh

	grpcServer.Stop()
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

func initDatabase(store repo.UserRepository) {
	store.DropAll()

	var user1 = model.User{ID: 105}
	err := store.WriteUser(&user1)

	var user2 = model.User{ID: 106}
	err = store.WriteUser(&user2)

	var user3 = model.User{ID: 107}
	err = store.WriteUser(&user3)

	var user4 = model.User{ID: 108}
	err = store.WriteUser(&user4)

	var user5 = model.User{ID: 109}
	err = store.WriteUser(&user5)

	var user6 = model.User{ID: 110}
	err = store.WriteUser(&user6)

	var user7 = model.User{ID: 111}
	err = store.WriteUser(&user7)

	var user8 = model.User{ID: 112}
	err = store.WriteUser(&user8)

	var user9 = model.User{ID: 113}
	err = store.WriteUser(&user9)

	var user10 = model.User{ID: 114}
	err = store.WriteUser(&user10)

	var user11 = model.User{ID: 115}
	err = store.WriteUser(&user11)

	var user12 = model.User{ID: 116}
	err = store.WriteUser(&user12)

	var user13 = model.User{ID: 117}
	err = store.WriteUser(&user13)

	var user14 = model.User{ID: 118}
	err = store.WriteUser(&user14)

	var user15 = model.User{ID: 119}
	err = store.WriteUser(&user15)

	var user16 = model.User{ID: 120}
	err = store.WriteUser(&user16)

	err = store.CreateFollowConnection(user2.ID, user3.ID)
	err = store.CreateFollowConnection(user3.ID, user2.ID)
	err = store.CreateFollowConnection(user3.ID, user5.ID)
	err = store.CreateFollowConnection(user3.ID, user6.ID)

	if err != nil {
		fmt.Println("Init database error")
	}

}

func main() {
	storeLogger := log.New(os.Stdout, "[person-store] ", log.LstdFlags)

	store, err := New(storeLogger)
	if err != nil {
		storeLogger.Fatal(err)
	}
	defer store.CloseDriverConnection(context.Background())
	store.CheckConnection()

	initDatabase(*store)

	userService := &service.UserService{UserRepo: store}
	//userHandler := &handler.UserHandler{UserService: userService}
	userHandlerGRPC := &handler.UserHandlergRPC{UserService: userService}

	startServerGRPC(userHandlerGRPC)

}
