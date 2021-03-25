package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"home24-technical-test/internal/http/controller"
	userAdapter "home24-technical-test/internal/user/adapter"
	"home24-technical-test/pkg/data"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
)

// Server represents http server
type Server struct {
	userController         *controller.UserController
	getUserAdapter         userAdapter.GetUserAdapter
	getLoginSessionAdapter userAdapter.GetLoginSessionAdapter
}

func (s *Server) compileRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTION"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Access-Token", "X-Requested-With"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(cors.Handler)

	// Prometheus handler
	//

	// Add routes
	//
	r.HandleFunc("/v1/login", s.userController.Login)
	r.Route("/v1", func(r chi.Router) {
		r.Use(s.authorizedOnly(s.getUserAdapter, s.getLoginSessionAdapter))

		r.Post("/logout", s.userController.Logout)
		r.Get("/session", s.userController.GetLoginSession)

		r.Group(func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.Put("/password", s.userController.ChangePassword)
			})
		})

	})

	return r
}

// Serve encapsulate process to listen and serve
func (s *Server) Serve(exposingPort string) {
	// Compile all the routes
	r := s.compileRouter()

	// Run the server
	log.Printf("About to listen on %s. Go to http://127.0.0.1:%s", exposingPort, exposingPort)
	srv := http.Server{Addr: fmt.Sprintf(":%s", exposingPort), Handler: r}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

// NewServer create new http server
func NewServer(
	getUserAdapter userAdapter.GetUserAdapter,
	getLoginSessionAdapter userAdapter.GetLoginSessionAdapter,
	loginAdapter userAdapter.LoginAdapter,
	logoutAdapter userAdapter.LogoutAdapter,
	changePasswordAdapter userAdapter.ChangePasswordAdapter,
	dataManager *data.Manager,
) *Server {
	userController := controller.NewUserController(getLoginSessionAdapter, loginAdapter, logoutAdapter, changePasswordAdapter, dataManager)

	return &Server{
		userController:         userController,
		getUserAdapter:         getUserAdapter,
		getLoginSessionAdapter: getLoginSessionAdapter,
	}
}
