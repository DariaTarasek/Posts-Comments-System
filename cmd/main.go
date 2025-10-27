package main

import (
	"OzonTestTask/internal/config"
	"OzonTestTask/internal/graphql/generated"
	"OzonTestTask/internal/graphql/resolvers"
	"OzonTestTask/internal/service/comment"
	"OzonTestTask/internal/service/post"
	in_memory "OzonTestTask/internal/storage/in-memory"
	"OzonTestTask/internal/storage/postgreSQL"
	"OzonTestTask/internal/subscription"
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Используются переменные окружения")
	}
	conf := config.GetConfig()
	fmt.Printf("Выбрано хранилище %s\n", conf.StorageType)
	var postService *post.PostService
	var commentService *comment.CommentService
	var subService subscription.Subscription

	if conf.StorageType == config.PostgresStorage {
		// подключение к БД
		db, err := postgreSQL.NewDBConnection(conf.PostgresDSN)
		if err != nil {
			log.Fatalf("не удалось подключиться к БД: %v", err)
		}
		fmt.Println("Подключено хранилище postgres")
		defer db.Close()

		// подключение к БД с pgx
		pool, err := subscription.NewPGXPool(conf.PostgresDSN)
		if err != nil {
			log.Fatalf("не удалось создать pgx pool: %v", err)
		}
		defer pool.Close()

		storage := postgreSQL.NewStorage(db)
		subService = subscription.NewPostgresSubscription(pool)
		postService = post.NewPostService(storage)
		commentService = comment.NewCommentService(storage, subService)

	} else if conf.StorageType == config.InMemoryStorage {
		subService = subscription.NewInMemorySubscription()
		inMemoryStorage := in_memory.NewInMemoryStorage()
		postService = post.NewPostService(inMemoryStorage)
		commentService = comment.NewCommentService(inMemoryStorage, subService)
		fmt.Println("Подключено in-memory хранилище")
	} else {
		log.Fatalf("неизвестный тип хранилища")
	}

	resolver := &resolvers.Resolver{
		PostService:         postService,
		CommentService:      commentService,
		SubscriptionService: subService,
	}

	server := handler.New(generated.NewExecutableSchema(generated.Config{
		Resolvers: resolver,
	}))
	server.AddTransport(transport.POST{})
	server.AddTransport(transport.GET{})
	server.AddTransport(transport.Websocket{})

	http.Handle("/", playground.Handler("GraphQL Playground", "/graphql"))
	http.Handle("/graphql", server)

	port := ":8080"

	httpServer := &http.Server{
		Addr:    port,
		Handler: nil,
	}

	go func() {
		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
		fmt.Printf("Сервер запущен на %s\n", port)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	httpServer.Shutdown(ctx)
	subService.Close()
}
