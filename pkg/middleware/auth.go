package middleware

import (
	"errors"
	"example/project-management-system/internal/utils/response"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
)

const (
	missingJWTErrorMessage = "Requires authentication"
	invalidJWTErrorMessage = "Bad credentials"
	internalServerErrorMessage = "Internal Server Error"
)

func ValidateJWT(audience, domain, env string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// For testing
		if env == "test" {
			next.ServeHTTP(w, r)
			return
		}

		issuerURL, err := url.Parse(domain)
		if err != nil {
			log.Fatalf("Failed to parse the issuer url: %v", err)
		}

		provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

		jwtValidator, err := validator.New(
			provider.KeyFunc,
			validator.HS256,
			issuerURL.String(),
			[]string{audience},
		)
		if err != nil {
			// log.Fatalf("Failed to set up the jwt validator")
			log.Printf("Failed to set up the jwt validator: %v", err)
			response.WriteJson(w, http.StatusInternalServerError, map[string]string{"message": "JWT validator setup failed"})
			return
		}

    if authHeaderParts := strings.Fields(r.Header.Get("Authorization")); len(authHeaderParts) > 0 && strings.ToLower(authHeaderParts[0]) != "bearer" {
			errorMessage := map[string]string{"message": invalidJWTErrorMessage}
			if err := response.WriteJson(w, http.StatusUnauthorized, errorMessage); err != nil {
				log.Printf("Failed to write error message: %v", err)
			}
			return
		}

		errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("Encountered error while validating JWT: %v", err)
			if errors.Is(err, jwtmiddleware.ErrJWTMissing) {
				errorMessage := map[string]string{"message": missingJWTErrorMessage}
				if err := response.WriteJson(w, http.StatusUnauthorized, errorMessage); err != nil {
					log.Printf("Failed to write error message: %v", err)
				}
				return
			}
			if errors.Is(err, jwtmiddleware.ErrJWTInvalid) {
				errorMessage := map[string]string{"message": invalidJWTErrorMessage}
				if err := response.WriteJson(w, http.StatusUnauthorized, errorMessage); err != nil {
					log.Printf("Failed to write error message: %v", err)
				}
				return
			}
			errorMessage := map[string]string{"message": internalServerErrorMessage}
			response.WriteJson(w, http.StatusInternalServerError, errorMessage)
		}

		middleware := jwtmiddleware.New(
			jwtValidator.ValidateToken,
			jwtmiddleware.WithErrorHandler(errorHandler),
		)

		middleware.CheckJWT(next).ServeHTTP(w, r)
	})
}