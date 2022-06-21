package v1

import (
	"github.com/labstack/echo/v4/middleware"
	"github.com/pmoieni/auth/api/v1/handlers"
)

func (s *Server) initRoutes() {
	s.Handler.Use(middleware.Secure())
	// TODO: add CSRF protection
	// s.Handler.Use(middleware.CSRF())
	s.Handler.Use(middleware.BodyLimit("32M"))
	// TODO: implement Redis for in memory store - https://echo.labstack.com/middleware/rate-limiter/
	// TODO: add timeout functionality - https://echo.labstack.com/middleware/timeout/
	s.Handler.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	v1Group := s.Handler.Group("/api/v1")

	accountGroup := v1Group.Group("/account")

	passwordResetGroup := accountGroup.Group("/password")

	passwordResetGroup.POST("/token", handlers.SendPasswordResetToken)
	passwordResetGroup.POST("/reset", handlers.ResetPassword)

	authGroup := v1Group.Group("/auth")

	directAuthGroup := authGroup.Group("/direct")

	directAuthGroup.POST("/register", handlers.Register)
	directAuthGroup.POST("/login", handlers.Login)
	directAuthGroup.GET("/token", handlers.RefreshToken)

	// OIDCAuthGroup := authGroup.Group("/oauth")

	// googleOIDCGroup := OIDCAuthGroup.Group("/google")
	// twitterOIDCGroup := OIDCAuthGroup.Group("/twitter")
	// facebookOIDCGroup := OIDCAuthGroup.Group("/facebook")
}
