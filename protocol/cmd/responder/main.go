package main

import (
	"context"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/99designs/gqlgen/handler"
	"github.com/go-chi/chi"
	"github.com/odysseyhack/planet-society/protocol/cryptography"
	"github.com/odysseyhack/planet-society/protocol/database"
	"github.com/odysseyhack/planet-society/protocol/models"
	"github.com/odysseyhack/planet-society/protocol/protocol"
	"github.com/odysseyhack/planet-society/protocol/transport"
	"github.com/odysseyhack/planet-society/protocol/utils"
	"github.com/phob0s-pl/generator"
	log "github.com/sirupsen/logrus"
)

const (
	dbFile = "hackaton_db.db"
)

var (
	keychain *cryptography.Keychain
)

func main() {
	utils.ConfigureLogger()
	log.Infoln("creating temporary directory")
	dir, err := ioutil.TempDir("", "responder")
	if err != nil {
		log.Fatalln("failed to create temporary dir:", err)
	}
	log.Infoln("using temporary directory:", dir)
	defer cleanup(dir)

	if err := createKeychain(); err != nil {
		log.Fatalln("failed to generate keychain:", err)
	}

	db, err := createDatabase(dir)
	if err != nil {
		log.Fatalln("failed to create database:", err)
	}

	defer func() {
		log.Infoln("closing database:", dir)
		if err := db.Close(); err != nil {
			log.Warningln("failed to close database:", err)
		}
	}()

	if err := serve(db); err != nil {
		log.Warningf("%s", err)
	}
}

func cleanup(dir string) {
	log.Infoln("removing public key")
	if err := utils.CleanKey(); err != nil {
		log.Warningln("cleaning key failed:", err)
	}

	log.Infoln("removing temporary directory:", dir)
	if err := os.RemoveAll(dir); err != nil {
		log.Warningln("failed to clean temporary directory:", err)
	}
}

func serve(db *database.Database) error {
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
	return ws.Listen(":15000")
}

func Middleware(db *database.Database) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := getMetadata(r)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func permissionFromHeader(r *http.Request) *models.Permission {
	return &models.Permission{
		Title:              r.Header.Get("title"),
		Description:        r.Header.Get("description"),
		RequesterSignature: r.Header.Get("signature"),
		ResponderSignature: sign(r.Header.Get("requester")),
		Expiration:         time.Now().Add(time.Hour * 120).Format(time.RFC3339),
		LawApplying:        "European Union",
	}
}

func getMetadata(r *http.Request) context.Context {
	var ctx = r.Context()
	transactionID := r.Header.Get("TransactionID")
	if transactionID != "" {
		log.Infoln("using transactionID:", transactionID)
		ctx = context.WithValue(ctx, "TransactionID", transactionID)
	} else {
		tid := cryptography.RandomKey32()
		ctx = context.WithValue(ctx, "TransactionID", tid.String())
	}
	k, _ := cryptography.Key32FromString(r.Header.Get("requester"))
	permission := permissionFromHeader(r)

	permission.RequesterPublicKey = models.Key32{Key: k}
	permission.TransactionID = transactionID

	checkPermissionType(r, permission)
	utils.AddPermission(permission)
	return ctx
}

func checkPermissionType(r *http.Request, permission *models.Permission) {
	if r.Header.Get("permission-type") == "digital telecommunication agreement" {
		permission.LegalReliationships = models.LegalReliationships{
			MyRights:       []string{"use telecommunication services until agreement expires"},
			TheirDuties:    []string{"provide high availability telecommunication service"},
			MyPowers:       []string{"cancel contract within 14 days from signing"},
			TheirLiability: []string{"liable for the consequences of the agreement termination"},
		}
	}
}

func sign(data string) string {
	signer := cryptography.NewSigner(keychain.SignaturePrivateKey, keychain.SignaturePublicKey)
	a, _ := signer.Sign([]byte(data))
	return hex.EncodeToString(a)
}

func createKeychain() (err error) {
	log.Infoln("generating one shot keychain")
	keychain, err = cryptography.OneShotKeychain()
	if err != nil {
		return err
	}

	if err := utils.WriteKeyToDir(keychain.MainPublicKey); err != nil {
		return err
	}

	log.Infoln("keychain main public key:", keychain.MainPublicKey.String())
	log.Infoln("keychain signature key:  ", keychain.SignaturePublicKey.String())
	return nil
}

func createDatabase(dir string) (*database.Database, error) {
	filePath := filepath.Join(dir, dbFile)
	log.Infoln("using database:", filePath)
	db, err := database.LoadDatabase(filePath, keychain)
	if err != nil {
		return nil, err
	}

	log.Infoln("filling database with items")
	dbGenerator := generator.NewGenerator()
	if err := dbGenerator.Generate(db); err != nil {
		return nil, err
	}

	return db, nil
}
