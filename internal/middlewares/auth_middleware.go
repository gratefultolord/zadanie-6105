package middlewares

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"zadanie-6105/pkg/utils"

	"gorm.io/gorm"
)

type contextKey string

const (
	userContextKey         = contextKey("username")
	organizationContextKey = contextKey("organizationID")
)

func AuthMiddleware(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var username string
			var organizationID string

			if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
				// Read the body into bytes
				bodyBytes, err := io.ReadAll(r.Body)
				if err != nil {
					log.Println("Failed to read request body:", err)
					utils.RespondWithError(w, http.StatusBadRequest, "Invalid request format")
					return
				}

				// Restore the body so the next handler can read it
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				// Define a struct for the authentication fields
				type AuthRequest struct {
					CreatorUsername string `json:"creatorUsername"`
					AuthorID        string `json:"authorId"`
					AuthorType      string `json:"authorType"`
				}

				// Decode the body bytes into the struct
				var authReq AuthRequest
				if err := json.Unmarshal(bodyBytes, &authReq); err != nil {
					log.Println("Failed to parse request body:", err)
					utils.RespondWithError(w, http.StatusBadRequest, "Invalid request format")
					return
				}

				// Handle authentication based on authorType
				if authReq.CreatorUsername != "" {
					username = authReq.CreatorUsername
				} else if authReq.AuthorID != "" && authReq.AuthorType != "" {
					if authReq.AuthorType == "User" {
						// Fetch username from the database using authorId
						var err error
						username, err = fetchUsernameByEmployeeID(r.Context(), db, authReq.AuthorID)
						if err != nil {
							log.Println("Failed to fetch username by authorId:", err)
							utils.RespondWithError(w, http.StatusUnauthorized, "User not authenticated")
							return
						}
					} else if authReq.AuthorType == "Organization" {
						// Fetch organization ID to verify existence
						exists, err := isOrganizationExists(r.Context(), db, authReq.AuthorID)
						if err != nil || !exists {
							log.Println("Failed to verify organizationId:", err)
							utils.RespondWithError(w, http.StatusUnauthorized, "Organization not authenticated")
							return
						}
						organizationID = authReq.AuthorID
					} else {
						log.Println("Invalid authorType provided")
						utils.RespondWithError(w, http.StatusBadRequest, "Invalid authorType")
						return
					}
				} else {
					log.Println("No creatorUsername or authorId/authorType found in request body")
					utils.RespondWithError(w, http.StatusUnauthorized, "User not authenticated")
					return
				}
			} else {
				username = r.URL.Query().Get("username")
				authorID := r.URL.Query().Get("authorId")
				authorType := r.URL.Query().Get("authorType")

				if username == "" && authorID != "" && authorType != "" {
					if authorType == "User" {
						var err error
						username, err = fetchUsernameByEmployeeID(r.Context(), db, authorID)
						if err != nil {
							log.Println("Failed to fetch username by authorId:", err)
							utils.RespondWithError(w, http.StatusUnauthorized, "User not authenticated")
							return
						}
					} else if authorType == "Organization" {
						exists, err := isOrganizationExists(r.Context(), db, authorID)
						if err != nil || !exists {
							log.Println("Failed to verify organizationId:", err)
							utils.RespondWithError(w, http.StatusUnauthorized, "Organization not authenticated")
							return
						}
						organizationID = authorID
					} else {
						log.Println("Invalid authorType provided")
						utils.RespondWithError(w, http.StatusBadRequest, "Invalid authorType")
						return
					}
				}
			}

			if username == "" && organizationID == "" {
				log.Println("No authentication information provided")
				utils.RespondWithError(w, http.StatusUnauthorized, "User not authenticated")
				return
			}

			// Add the username or organizationID to the context
			ctx := r.Context()
			if username != "" {
				ctx = context.WithValue(ctx, userContextKey, username)
			}
			if organizationID != "" {
				ctx = context.WithValue(ctx, organizationContextKey, organizationID)
			}

			// Proceed to the next handler with the updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// fetchUsernameByEmployeeID retrieves the username associated with the given employee ID
func fetchUsernameByEmployeeID(ctx context.Context, db *gorm.DB, employeeID string) (string, error) {
	var username string
	err := db.WithContext(ctx).
		Table("employee").
		Select("username").
		Where("id = ?", employeeID).
		Scan(&username).Error
	if err != nil {
		log.Println("Database error in fetchUsernameByEmployeeID:", err)
		return "", err
	}
	if username == "" {
		log.Println("No user found with employee ID:", employeeID)
		return "", gorm.ErrRecordNotFound
	}
	log.Printf("Username fetched for employee ID %s: %s", employeeID, username)
	return username, nil
}

// isOrganizationExists checks if an organization with the given ID exists
func isOrganizationExists(ctx context.Context, db *gorm.DB, organizationID string) (bool, error) {
	var count int64
	err := db.WithContext(ctx).
		Table("organization").
		Where("id = ?", organizationID).
		Count(&count).Error
	if err != nil {
		log.Println("Database error in isOrganizationExists:", err)
		return false, err
	}
	return count > 0, nil
}

func GetUsernameFromContext(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(userContextKey).(string)
	return username, ok
}

func GetOrganizationIDFromContext(ctx context.Context) (string, bool) {
	orgID, ok := ctx.Value(organizationContextKey).(string)
	return orgID, ok
}
