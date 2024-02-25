package setBalance

import (
	resp "CassetteRental/internal/lib/api/response"
	"CassetteRental/internal/storage"
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type Request struct {
	Balance int `json:"balance" validate:"required"`
}

type Response struct {
	resp.Response
}

type CustomerBalanceSetter interface {
	SetCustomerBalance(ctx context.Context, id string, balance int) (context.Context, error)
}

func New(log *slog.Logger, setter CustomerBalanceSetter) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handlers/customer/setBalance/setBalance/New"

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
			log.Error("Invalid request", slog.String("error", err.Error()))
			resp.BadRequest(writer, "Invalid request")

			return
		}

		id := chi.URLParam(request, "customerId")
		if id == "" {
			log.Info("id is empty")
			resp.BadRequest(writer, "Invalid url id is empty")
			return
		}

		ctx := context.Background()
		_, err = setter.SetCustomerBalance(ctx, id, req.Balance)
		if errors.Is(err, storage.ErrCustomerNotFound) {
			log.Error("failed to set balance", slog.String("error", err.Error()))
			resp.StatusNotFound(writer, "customer not found")
			return
		}
		if err != nil {
			log.Error("failed to set balance", slog.String("error", err.Error()))
			resp.InternalServerError(writer, "failed to set balance")
			return
		}
		log.Info("balance is changed", slog.String("set balance customer with id", id))

		render.JSON(writer, request, Response{
			Response: resp.Ok(),
		})
	}
}
