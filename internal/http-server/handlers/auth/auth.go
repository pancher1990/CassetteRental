package auth

import (
	resp "CassetteRental/internal/lib/api/response"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"time"
)

type Request struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	resp.Response
	Token string `json:"token,omitempty"`
}

type CustomerGetter interface {
	GetCustomerIdByCredentials(email string, password string) (string, error)
}

type Hasher interface {
	Hash(password string) (string, error)
}

type TokenManager interface {
	NewJWT(userId string, ttl time.Duration) (string, error)
}

func New(log *slog.Logger, s CustomerGetter, m TokenManager, h Hasher) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handlers/auth/auth/New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(request.Context())),
		)
		var req Request

		err := render.DecodeJSON(request.Body, &req)

		if err != nil {
			log.Error("Failed to decode request body ", slog.String("error", err.Error()))
			render.JSON(writer, request, resp.Error("Failed to decode request"))
			return
		}

		log.Info("request body decoded ", slog.Any("request", req))
		if err = validator.New().Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)
			log.Error("Invalid request", slog.String("error", err.Error()))
			render.JSON(writer, request, resp.ValidationError(validateErr))
			return
		}

		passwordHash, err := h.Hash(req.Password)
		if err != nil {
			log.Error("error with hash", slog.String("error", err.Error()))
			render.JSON(writer, request, resp.Error("failed to 	hash password"))
			return
		}
		customerId, err := s.GetCustomerIdByCredentials(req.Email, passwordHash)
		if err != nil {
			log.Error("failed get credentials", slog.String("error", err.Error()))

			render.JSON(writer, request, resp.Error("failed to log in "))
			return
		}

		ttl := 3 * time.Minute
		token, err := m.NewJWT(customerId, ttl)
		if err != nil {
			log.Error("failed to authorization authorization", slog.String("error", err.Error()))
			render.JSON(writer, request, resp.Error(err.Error()))
			return
		}
		log.Info("auth complete", slog.String("user authorized with id", customerId))
		render.JSON(writer, request, Response{
			Response: resp.Ok(),
			Token:    token,
		})
	}
}
