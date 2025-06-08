package router

import (
	"fmt"
	"kmem/internal/config"
	"kmem/internal/db"
	"kmem/internal/models"
	"kmem/internal/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func checkUserPostBody(ctx *gin.Context, user *models.User) error {
	if err := ctx.Bind(&user); err != nil {
		models.ErrorResponse(
			http.StatusBadRequest,
			models.ErrInvalidInput,
			"invalid request format",
		).Send(ctx)

		return fmt.Errorf("failed to bind user body")
	}

	if len(user.Username) < 4 {
		models.ErrorResponse(
			http.StatusBadRequest,
			models.ErrInvalidInput,
			"user name must be at least 4 characters",
		).Send(ctx)

		return fmt.Errorf("short username")
	}

	if len(user.Password) < 8 {
		models.ErrorResponse(http.StatusBadRequest,
			models.ErrInvalidInput,
			"password must be at least 8 characters").Send(ctx)

		return fmt.Errorf("short password")
	}

	return nil
}

// signup handler
func signup(pg *db.Postgres) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user models.User

		// response to client handled already
		if err := checkUserPostBody(ctx, &user); err != nil {
			log.Println(err)
			return
		}

		if err := pg.InsertUser(user); err != nil {
			models.ErrorResponse(
				http.StatusInternalServerError,
				models.ErrDatabase,
				"failed to create user",
			).Send(ctx)

			return
		}

		models.SuccessResponse(nil).Send(ctx)
	}
}

// login handler
func login(pg *db.Postgres, conf *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user models.User

		if err := checkUserPostBody(ctx, &user); err != nil {
			log.Println(err)
			return
		}

		// check password
		dbuser, err := pg.QueryUser(user.Username)
		if err != nil {
			models.ErrorResponse(
				http.StatusUnauthorized,
				models.ErrUnauthorized,
				"incorrect username or password",
			).Send(ctx)

			return
		}

		if res := utils.CheckPasswordHash(dbuser.Password, user.Password); !res {
			models.ErrorResponse(
				http.StatusUnauthorized,
				models.ErrUnauthorized,
				"incorrect username or password",
			).Send(ctx)

			return
		}

		accessToken, err := utils.GenTokenString(conf.JwtSecretKey(), user.Username, utils.ACCESS_TOKEN_DUR)
		if err != nil {
			models.ErrorResponse(
				http.StatusInternalServerError,
				models.ErrInvalidToken,
				"failed to create token",
			).Send(ctx)

			return
		}

		refreshToken, err := utils.GenTokenString(conf.JwtSecretKey(), user.Username, utils.REFRESH_TOKEN_DUR)
		if err != nil {
			models.ErrorResponse(
				http.StatusInternalServerError,
				models.ErrInvalidToken,
				"failed to create token",
			).Send(ctx)

			return
		}

		ctx.SetCookie(utils.ACCESS_TOKEN_KEY, accessToken, int(utils.ACCESS_TOKEN_DUR), "/", "", false, true)
		ctx.SetCookie(utils.REFRESH_TOKEN_KEY, refreshToken, int(utils.REFRESH_TOKEN_DUR), "/", "", false, true)

		if err := pg.UpdateLastLogin(user.Username); err != nil {
			log.Printf("failed to update last login for user %s: %v", user.Username, err)
		}

		models.SuccessResponse(nil).Send(ctx)
	}
}

// logout handler
func logout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.SetCookie(utils.ACCESS_TOKEN_KEY, "", -1, "/", "", false, true)
		ctx.SetCookie(utils.REFRESH_TOKEN_KEY, "", -1, "/", "", false, true)

		models.SuccessResponse(nil).Send(ctx)
	}
}

func validateToken(ctx *gin.Context, cookieName, jwtSecret string) (string, error) {
	tokenStr, err := ctx.Cookie(cookieName)
	if err != nil {
		return "", fmt.Errorf("failed to find cookie: %v", err)
	}

	token, err := utils.ParseToken(jwtSecret, tokenStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %v", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("failed to get claims from token")
	}

	val, ok := claims["username"]
	if !ok {
		return "", fmt.Errorf("failed to get username from token")
	}

	username, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("invalid token claim")
	}

	return username, nil
}

func me() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		v, ok := ctx.Get(utils.USERNAME_KEY)
		if !ok {
			models.ErrorResponse(
				http.StatusBadRequest,
				models.ErrUnauthorized,
				"user unauthorized",
			).Send(ctx)

			return
		}

		username, ok := v.(string)
		if !ok {
			models.ErrorResponse(
				http.StatusBadRequest,
				models.ErrUnauthorized,
				"user unauthorized",
			).Send(ctx)

			return
		}

		if len(username) < 4 {
			models.ErrorResponse(
				http.StatusBadRequest,
				models.ErrInvalidInput,
				"user name must be at least 4 characters",
			).Send(ctx)

			return
		}

		models.SuccessResponse(nil).Send(ctx)
	}
}

// check access token & refresh token from cookies
// if access token is expired, refresh
func authMiddleware(conf *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		jwtSecret := conf.JwtSecretKey()

		username, err := validateToken(ctx, utils.ACCESS_TOKEN_KEY, jwtSecret)
		if err != nil {
			log.Printf("failed to validate access token: %v\n", err)
			log.Println("getting refresh token...")

			username, refreshErr := validateToken(ctx, utils.REFRESH_TOKEN_KEY, jwtSecret)
			if refreshErr != nil {
				models.ErrorResponse(
					http.StatusUnauthorized,
					models.ErrInvalidToken,
					"failed to validate token",
				).Send(ctx)

				ctx.Abort()

				return
			}

			// renew tokens
			accessToken, err := utils.GenTokenString(conf.JwtSecretKey(), username, utils.ACCESS_TOKEN_DUR)
			if err != nil {
				models.ErrorResponse(
					http.StatusInternalServerError,
					models.ErrInvalidToken,
					"failed to create token",
				).Send(ctx)

				ctx.Abort()

				return
			}

			refreshToken, err := utils.GenTokenString(conf.JwtSecretKey(), username, utils.REFRESH_TOKEN_DUR)
			if err != nil {
				models.ErrorResponse(
					http.StatusInternalServerError,
					models.ErrInvalidToken,
					"failed to create token",
				).Send(ctx)

				ctx.Abort()

				return
			}

			ctx.SetCookie(utils.ACCESS_TOKEN_KEY, accessToken, int(utils.ACCESS_TOKEN_DUR), "/", "", false, true)
			ctx.SetCookie(utils.REFRESH_TOKEN_KEY, refreshToken, int(utils.REFRESH_TOKEN_DUR), "/", "", false, true)
		}

		ctx.Set(utils.USERNAME_KEY, username)
		ctx.Next()
	}
}
