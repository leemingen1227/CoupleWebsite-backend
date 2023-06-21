package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	db "github.com/leemingen1227/couple-server/db/sqlc"
	"github.com/leemingen1227/couple-server/token"
	"github.com/leemingen1227/couple-server/util"
)

const (
	authorizationHeaderKey = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker, store db.Store, config util.Config, redisClient *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := fmt.Errorf("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := fmt.Errorf("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			// check the session id in cookie
			sessionID_string, err := ctx.Cookie("session_id")
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
				return
			}

			// Try to get the session from Redis first
			sessionData, err := redisClient.Get(ctx, sessionID_string).Result()

			// If the session is not found in Redis, or if there is an error,
			// then get the session from PostgreSQL
			var session db.Session
			if err != nil || sessionData == "" {
				//parse the session id
				sessionID, err := uuid.Parse(sessionID_string)
				if err != nil {
					ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
					return
				}

				//retrieve the session
				session, err = store.GetSession(ctx, sessionID)
				if err != nil {
					ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
					return
				}

				// Convert the session to a JSON string
				sessionJSON, err := json.Marshal(session)
				if err != nil {
					ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
					return
				}

				// Store the session data in Redis for future requests
				err = redisClient.Set(ctx, sessionID_string, sessionJSON, 24*time.Hour).Err()
				if err != nil {
					ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
					return
				}
			}else{
				// Unmarshal session data from Redis
				
				err = json.Unmarshal([]byte(sessionData), &session)
				if err != nil {
					// handle error
					ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
				}
				log.Println("session from redis: ", session)
			}

				

			refreshToken_payload, err := tokenMaker.VerifyToken(session.RefreshToken)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
				return
			}

			//create new access token
			user, err := store.GetUser(ctx, refreshToken_payload.UserID)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
				return
			}

			newAccessToken, newPayload, err := tokenMaker.CreateToken(user.ID, config.AccessTokenDuration)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
				return
			}

			// Update the access token in the response header
			ctx.Header("Authorization", fmt.Sprintf("Bearer %s", newAccessToken))
			// Optionally, you can also update the access token in the response payload if needed

			ctx.Set(authorizationPayloadKey, newPayload)

			// Continue with the request
			ctx.Next()
			return

		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}