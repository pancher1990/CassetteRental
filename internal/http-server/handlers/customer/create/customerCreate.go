package customerCreate

import (
	resp "CassetteRental/internal/lib/api/response"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type Request struct {
	Name string `json:"name" validate:"required"`
	//IsActive bool   `json:"isActive,omitempty" validate:"boolean"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	resp.Response
	Id string `json:"id,omitempty"`
}

type CustomerSaver interface {
	AddNewCustomer(name string, isActive bool, email string, hashPassword string) (string, error)
}

type Hasher interface {
	Hash(password string) (string, error)
}

func New(log *slog.Logger, saver CustomerSaver, hasher Hasher) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handlers/customer/create/customerCreate/New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(request.Context())),
		)
		var req Request

		err := render.DecodeJSON(request.Body, &req)

		if err != nil {
			log.Error("Failed to decode request body ", slog.String("error", err.Error()))
			resp.BadRequest(writer, "Failed to decode request body")

			return
		}

		log.Info("request body decoded ", slog.Any("request", req))
		if err = validator.New().Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)
			log.Error("Invalid request ", slog.String("error", err.Error()))
			resp.BadRequest(writer, "Invalid request")

			return
		}
		hashPassword, err := hasher.Hash(req.Password)
		if err != nil {
			log.Error("Failed to add customer", slog.String("error", err.Error()))
			render.JSON(writer, request, resp.Error("Error with hash"))
		}
		id, err := saver.AddNewCustomer(req.Name, true, req.Email, hashPassword)
		if err != nil {
			log.Error("failed to add customer", slog.String("error", err.Error()))

			render.JSON(writer, request, resp.Error("failed to add customer"))
			return
		}
		log.Info("customer added", slog.String("create customer with id", id))

		render.JSON(writer, request, Response{
			Response: resp.Ok(),
			Id:       id,
		})
	}
}
