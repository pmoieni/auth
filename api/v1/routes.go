package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pmoieni/auth/api/v1/handlers"
	"github.com/pmoieni/auth/models"
	"github.com/pmoieni/auth/providers/direct"
	"github.com/pmoieni/auth/store"
)

type userInfoRes struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

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

	meGroup := accountGroup.Group("/me")
	meGroup.Use(direct.CheckAuth)
	// example route for getting user info
	meGroup.GET("/", func(c echo.Context) error {
		userEmail, ok := c.Get("user-email").(string)
		if !ok {
			return &models.ErrorResponse{Status: http.StatusUnauthorized, Message: http.StatusText(http.StatusUnauthorized)}
		}
		u := store.User{
			Email: userEmail,
		}
		userInfo, err := u.GetUser()
		if err != nil {
			return err
		}

		userInfoRes := userInfoRes{
			Username: userInfo.Username,
			Email:    userInfo.Email,
		}

		return c.JSON(http.StatusOK, userInfoRes)
	})

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
