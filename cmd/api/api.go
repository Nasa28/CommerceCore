package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Nasa28/CommerceCore/cmd/service/product"
	"github.com/Nasa28/CommerceCore/cmd/service/role"
	"github.com/Nasa28/CommerceCore/cmd/service/user"
	"github.com/gorilla/mux"
)

type APIServer struct {
	port string
	db   *sql.DB
}

func NewAPIServer(port string, db *sql.DB) *APIServer {
	return &APIServer{
		port: port,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()

	subRouter := router.PathPrefix("/api/v1").Subrouter()
	userStore := user.NewStore(s.db)
	userServiceHandler := user.NewHandler(userStore)
	userServiceHandler.RegisterRoutes(subRouter)

	productRepository := product.NewStore(s.db)
	productServiceHandler := product.NewProductHandler(productRepository, userStore)
	productServiceHandler.RegisterRoutes(subRouter)

	roleStore := role.NewRolesStore(s.db)
	roleService := role.NewRoleHandler(roleStore)
	roleService.RegisterRoutes(subRouter)
	log.Println("App listening on port:", s.port)
	return http.ListenAndServe(s.port, router)
}
