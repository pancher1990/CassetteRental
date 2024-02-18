package setBalance

import (
	resp "CassetteRental/internal/lib/api/response"
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

		id := chi.URLParam(request, "id")
		if id == "" {
			log.Info("id is empty")

			render.JSON(writer, request, resp.Error("invalid request"))

			return
		}

		ctx := context.Background()
		_, err = setter.SetCustomerBalance(ctx, id, req.Balance)
		if err != nil {
			log.Error("failed to set balance", slog.String("error", err.Error()))

			render.JSON(writer, request, resp.Error("failed to set balance"))
			return
		}
		log.Info("balance is changed", slog.String("set balance customer with id", id))

		render.JSON(writer, request, Response{
			Response: resp.Ok(),
		})
	}
}
