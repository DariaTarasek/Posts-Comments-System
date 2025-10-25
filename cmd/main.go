package main

import (
	"OzonTestTask/internal/config"
	"OzonTestTask/internal/graphql/generated"
	"OzonTestTask/internal/graphql/resolvers"
	"OzonTestTask/internal/service/comment"
	"OzonTestTask/internal/service/post"
	in_memory "OzonTestTask/internal/storage/in-memory"
	"OzonTestTask/internal/storage/postgreSQL"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Используются переменные окружения")
	}
	conf := config.GetConfig()
	fmt.Printf("Выбрано хранилище %s\n", conf.StorageType)
	var postService *post.PostService
	var commentService *comment.CommentService

	if conf.StorageType == config.PostgresStorage {
		db, err := postgreSQL.NewDBConnection(conf.PostgresDSN)
		if err != nil {
			log.Fatalf("не удалось подключиться к БД: %v", err)
		}
		fmt.Println("Подключено хранилище postgres")
		defer db.Close()
		storage := postgreSQL.NewStorage(db)
		postService = post.NewPostService(storage)
		commentService = comment.NewCommentService(storage)
	} else if conf.StorageType == config.InMemoryStorage {
		inMemoryStorage := in_memory.NewInMemoryStorage()
		postService = post.NewPostService(inMemoryStorage)
		commentService = comment.NewCommentService(inMemoryStorage)
		fmt.Println("Подключено in-memory хранилище")
	} else {
		log.Fatalf("неизвестный тип хранилища")
	}

	resolver := &resolvers.Resolver{
		PostService:    postService,
		CommentService: commentService,
	}

	server := handler.New(generated.NewExecutableSchema(generated.Config{
		Resolvers: resolver,
	}))
	server.AddTransport(transport.POST{})
	server.AddTransport(transport.GET{})

	http.Handle("/", playground.Handler("GraphQL Playground", "/graphql"))
	http.Handle("/graphql", server)

	port := ":8080"
	fmt.Printf("Сервер запущен на порту %s/\n", port)
	log.Fatal(http.ListenAndServe(port, nil))

}
