package utils

import (
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"
	"time"

	"dkhalife.com/journey/internal/models"
	"dkhalife.com/journey/internal/repos"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	jwt "github.com/appleboy/gin-jwt/v2"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		log.Printf("IP:%s UA:%q Route:%s Status:%d\n",
			c.ClientIP(), c.Request.UserAgent(), c.Request.URL.Path, c.Writer.Status())
	}
}

type signIn struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func NewAuthMiddleware(userRepo *repos.UserRepository) (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "My Journey",
		Key:         []byte("SampleSecretKey"),
		Timeout:     time.Duration(24) * time.Hour,
		MaxRefresh:  time.Duration(24*7) * time.Hour,
		IdentityKey: IdentityKey,
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var req signIn
			if err := c.ShouldBindJSON(&req); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}

			user, err := userRepo.FindByEmail(c.Request.Context(), req.Email)
			if err != nil || user.Disabled {
				return nil, jwt.ErrFailedAuthentication
			}

			err = Matches(user.Password, req.Password)
			if err != nil {
				if err != bcrypt.ErrMismatchedHashAndPassword {
					log.Fatalf("found unknown error when matches password", "err", err)
				}
				return nil, jwt.ErrFailedAuthentication
			}

			return user, nil
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if u, ok := data.(*models.User); ok {
				return jwt.MapClaims{
					IdentityKey: fmt.Sprintf("%d", u.ID),
					"type":      "user",
					"scopes": []string{
						"history:read",
						"sensor:read",
						"sensor:write",
						"token:write",
						"user:read",
						"user:write",
					},
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)

			id, ok := claims[IdentityKey].(string)
			if !ok {
				log.Fatalf("failed to extract ID from claims")
				return nil
			}
			userID, err := strconv.Atoi(id)
			if err != nil {
				return nil
			}

			scopesRaw, ok := claims["scopes"].([]interface{})
			if !ok {
				return nil
			}

			var scopes []models.ApiTokenScope
			for _, scope := range scopesRaw {
				if s, ok := scope.(string); ok {
					scopes = append(scopes, models.ApiTokenScope(s))
				}
			}

			return &models.SignedInIdentity{
				UserID: userID,
				Type:   models.IdentityTypeUser,
				Scopes: scopes,
			}
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*models.SignedInIdentity); ok {
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"error": message,
			})
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.JSON(http.StatusOK, gin.H{
				"token":      token,
				"expiration": expire,
			})
		},
		RefreshResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.JSON(http.StatusOK, gin.H{
				"token":      token,
				"expiration": expire,
			})
		},
		TokenLookup:   "header: Authorization",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
}

func ScopeMiddleware(requiredScope models.ApiTokenScope) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentIdentity := CurrentIdentity(c)

		if slices.Contains(currentIdentity.Scopes, requiredScope) {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "Missing required scope: " + requiredScope,
		})
	}
}
