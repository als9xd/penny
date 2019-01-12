package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/jmoiron/sqlx"

	"github.com/als9xd/penny/server/models"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

const signingKey = "Change Me!"
const identityKey = "id"

func AuthMiddleware(db *sqlx.DB) *jwt.GinJWTMiddleware {

	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte(signingKey),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.Profile); ok {
				return jwt.MapClaims{
					identityKey: v.Id,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &models.Profile{
				Id: int(claims["id"].(float64)),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			profile := models.Profile{}
			err := db.Get(&profile, `
			SELECT * FROM profile WHERE username = $1;
			`, loginVals.Username)
			if err == sql.ErrNoRows {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": gin.H{
						"code":    http.StatusUnauthorized,
						"message": fmt.Sprintf("No profile with username '%s'", loginVals.Username),
					},
				})
				return nil, jwt.ErrFailedAuthentication
			}
			if err != nil {
				log.Fatal(err)
			}
			err = bcrypt.CompareHashAndPassword([]byte(profile.Password), []byte(loginVals.Password))
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}
			return &profile, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*models.Profile); ok && v.Id != 0 {
				var profileExists bool
				err := db.Get(&profileExists, "SELECT 1 FROM profile WHERE id = $1;", v.Id)
				if err != nil && err != sql.ErrNoRows {
					log.Fatal(err)
				}

				return profileExists
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
	return authMiddleware
}
