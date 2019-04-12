package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/99designs/gqlgen/handler"
	"github.com/go-chi/chi"
	"github.com/odysseyhack/planet-society/protocol/cryptography"
	"github.com/odysseyhack/planet-society/protocol/database"
	"github.com/odysseyhack/planet-society/protocol/protocol"
	"github.com/odysseyhack/planet-society/protocol/transport"
	"github.com/odysseyhack/planet-society/protocol/utils"
	log "github.com/sirupsen/logrus"
)

const (
	dbFile = "hackaton_db.db"
)

func main() {
	utils.ConfigureLogger()
	log.Infoln("responder application")

	log.Infoln("creating temporary directory")
	dir, err := ioutil.TempDir("", "responder")
	if err != nil {
		log.Fatalln("failed to create temporary dir:", err)
	}
	log.Infoln("using temporary directory:", dir)

	defer func() {
		log.Infoln("removing temporary directory:", dir)
		if err := os.RemoveAll(dir); err != nil {
			log.Warningln("failed to clean temporary directory:", err)
		}
	}()

	log.Infoln("generating one shot keychain")
	keychain, err := cryptography.OneShotKeychain()
	if err != nil {
		log.Fatalln("failed to generate keychain:", err)
	}
	log.Infoln("keychain main public key:", keychain.MainPublicKey.String())
	log.Infoln("keychain signature key:  ", keychain.SignaturePublicKey.String())

	if err := utils.WriteKeyToDir(keychain.MainPublicKey); err != nil {
		log.Fatalln("failed to write public key:", err)
	}

	defer func() {
		log.Infoln("removing public key")
		if err := utils.CleanKey(); err != nil {
			log.Warningln("cleaning key failed:", err)
		}
	}()

	filePath := filepath.Join(dir, dbFile)
	log.Infoln("using database:", filePath)
	db, err := database.LoadDatabase(filePath, keychain)
	if err != nil {
		log.Fatalln("failed to load database:", err)
	}

	defer func() {
		log.Infoln("closing database:", dir)
		if err := db.Close(); err != nil {
			log.Warningln("failed to close database:", err)
		}
	}()

	log.Infoln("filling database with items")
	generator := newGenerator()
	if err := generator.generate(db); err != nil {
		log.Fatalln("failed to fill db with data:", err)
	}

	router := chi.NewRouter()
	router.Use(Middleware(db))
	resolver := protocol.NewResolver(db)
	router.Handle("/", handler.Playground("GraphQL playground", "/query"))
	router.Handle("/query", handler.GraphQL(protocol.NewExecutableSchema(protocol.Config{Resolvers: resolver})))
	go func() {
		log.Warningln(http.ListenAndServe(":8088", router))
	}()

	proto := protocol.NewProtocol(&AlwaysAcceptPlugin{})
	go proto.Loop()

	ws := transport.NewWebsocket(proto.Connections)
	log.Infoln("starting listener at: :15000")
	if err := ws.Listen(":15000"); err != nil {
		log.Warningln("transport listener failed: %s", err)
	}
}

func Middleware(db *database.Database) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// generate random key
			key := cryptography.RandomKey32()

			// put it in context
			ctx := context.WithValue(r.Context(), "transaction-key", key.String())
			// and call the next with our new context
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
