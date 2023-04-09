package api

import (
	"fmt"

	db "userMicroService/db/sqlc"
	"userMicroService/token"
	"userMicroService/util"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// Creates a new HTTP server and setup routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setUpRouter()
	return server, nil
}

func (server *Server) setUpRouter() {
	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

	g1 := router.Group("/api/user")
	g1.POST("/", server.createUser)
	g1.POST("/login", server.logInUser)
	g1.POST("/forgotpassword", server.forgotPassword)

	authRoutes := router.Group("/api/user").Use(authMiddleware(server.tokenMaker))
	authRoutes.DELETE("/deleteUser/:username", server.deleteUser)
	authRoutes.PATCH("/resetPassword", server.resetPassword)

	authRoutes.GET("/:username", server.getUser)

	server.router = router
}

// Runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
