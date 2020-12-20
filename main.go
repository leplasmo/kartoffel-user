// kartoffel-user/main.go

package main

import (
	"log"

	pb "github.com/leplasmo/kartoffel-user/proto/user"
	"github.com/micro/go-micro/v2"
)

const schema = `
	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(36) NOT NULL,
		name VARCHAR(125) NOT NULL,
		email VARCHAR(225) NOT NULL UNIQUE,
		password VARCHAR(225) NOT NULL,
		company VARCHAR(125),
		PRIMARY KEY (id)
	);
`

func main() {

	// Create the database connection
	db, err := NewConnection()
	if err != nil {
		log.Panic(err)
	}

	defer db.Close()

	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}

	// Run the schema to create the table if it does not exist
	db.MustExec(schema)

	repo := NewPostgresRepository(db)

	tokenService := &TokenService{repo}

	// Create a new service
	service := micro.NewService(
		micro.Name("kartoffel.user"),
		micro.Version("0.1.0"),
	)

	// Init parses the micro command line flags
	service.Init()

	// Register handler
	if err := pb.RegisterUserServiceHandler(service.Server(), &handler{repo, tokenService}); err != nil {
		log.Panic(err)
	}

	// Run the server
	if err := service.Run(); err != nil {
		log.Panic(err)
	}
}
