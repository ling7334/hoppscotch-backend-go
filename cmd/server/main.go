package main

import (
	"net/http"
	"os"
	"time"

	"graph"
	mw "middleware"
	"model"
	"rest"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	defaultPort = "8080"
)

var (
	dsn string
)

func init() {
	if _, err := os.Stat(".env"); err == nil {
		godotenv.Load(".env")
	}
	dsn = os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal().Msg("DATABASE_URL is not set")
	}
}

func initDB(db *gorm.DB) {

	// Create the database. This is a one-time step.
	// Comment out if running multiple times - You may see an error otherwise
	// db.Exec("CREATE SCHEMA public")
	// db.Exec("USE public")

	db.Exec("CREATE TYPE Team_Member_Role AS ENUM('OWNER', 'VIEWER', 'EDITOR');")
	db.Exec("CREATE TYPE Req_Type AS ENUM('REST', 'GQL');")

	// Migration to create tables for Order and Item schema
	db.AutoMigrate(

		&model.User{},
		&model.UserCollection{},
		&model.UserEnvironment{},
		&model.UserHistory{},
		&model.UserRequest{},
		&model.UserSetting{},
		&model.VerificationToken{},
		&model.Team{},
		&model.TeamCollection{},
		&model.TeamEnvironment{},
		&model.TeamInvitation{},
		&model.TeamMember{},
		&model.TeamRequest{},
		&model.Account{},
		&model.InfraConfig{},
		&model.InvitedUser{},
		&model.Shortcode{},
	)
}

func main() {
	var err error

	// dataSourceName := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect database")
	}

	// initDB(db)

	http.HandleFunc("/ping", rest.Ping)

	http.Handle("/v1/auth/", mw.LogMiddleware(mw.DBMiddleware(db, rest.ServeMux("/v1/auth/"))))

	http.Handle("/v1/team-collection/", mw.LogMiddleware(mw.OperatorMiddleware(db, mw.DBMiddleware(db, rest.TeamServeMux("/v1/team-collection/")))))

	// gqlgen config
	c := graph.Config{Resolvers: &graph.Resolver{
		DB: db,
	}}
	// gqlgen Directives
	c.Directives.IsAdmin = graph.IsAdmin
	c.Directives.IsLogin = graph.IsLogin

	srv := handler.New(graph.NewExecutableSchema(c))

	// Websocket support
	srv.AddTransport(&transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	// srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New(1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})

	// Complexity Limit
	srv.Use(extension.FixedComplexityLimit(50))

	// playground
	if os.Getenv("PRODUCTION") == "false" {
		http.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	}
	http.Handle("/graphql", mw.LogMiddleware(mw.JwtMiddleware(mw.OperatorMiddleware(db, srv))))

	log.Info().Msgf("listen on :%s", defaultPort)
	log.Fatal().Err(http.ListenAndServe(":"+defaultPort, nil)).Msg("something went wrong")
}
