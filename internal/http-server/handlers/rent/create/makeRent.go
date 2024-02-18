package makeRent

import (
	resp "CassetteRental/internal/lib/api/response"
	"context"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type Request struct {
	Title    string `json:"title" validate:"required"`
	RentDays int    `json:"rentDays" validate:"required"`
}

type Response struct {
	resp.Response
	OrderId string `json:"orderId,omitempty"`
	RentId  string `json:"rentId,omitempty"`
}

type RentMaker interface {
	GetFilm(ctx context.Context, title string) (context.Context, string, int, error)
	GetCustomerBalance(ctx context.Context, id string) (context.Context, int, error)
	FindAvailableCassette(ctx context.Context, filmId string) (context.Context, string, error)
	SetCassetteStatus(ctx context.Context, id string, available bool) (context.Context, error)
	CreateOrder(ctx context.Context, customerId string) (context.Context, string, error)
	CreateCassetteInOrder(ctx context.Context, cassetteId string, orderId string, rentCost int) (context.Context, error)
	CreateRent(ctx context.Context, customerId string, cassetteId string, rentDays int) (context.Context, string, error)
	SetCustomerBalance(ctx context.Context, id string, balance int) (context.Context, error)
}

func New(log *slog.Logger, rentMaker RentMaker) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handlers/rent/create/makeRent/New"

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

		customerId, ok := request.Context().Value("customerId").(string)
		if ok == false {
			log.Error("Failed to take customerId")
			render.JSON(writer, request, resp.Error("Failed to take customerId"))
		}

		ctx := context.Background()
		ctx = context.WithValue(ctx, "returnTransaction", true)

		ctx, filmId, dayPrice, err := rentMaker.GetFilm(ctx, req.Title)
		if err != nil {
			log.Error("failed to create rent ", slog.String("error", err.Error()))
			render.JSON(writer, request, resp.Error("failed to create rent"))
			return
		}

		ctx, balance, err := rentMaker.GetCustomerBalance(ctx, customerId)
		if err != nil {
			log.Error("failed to create rent ", slog.String("error", err.Error()))
			render.JSON(writer, request, resp.Error("failed to create rent"))
			return
		}

		rentCost := req.RentDays * dayPrice
		if rentCost > balance {
			log.Error("failed to create rent insufficient funds on the balance", slog.String("error", err.Error()))
			render.JSON(writer, request, resp.Error("failed to create rent insufficient funds on the balance"))
			return
		}

		ctx, cassetteId, err := rentMaker.FindAvailableCassette(ctx, filmId)
		if err != nil {
			log.Error("failed to create rent ", slog.String("error", err.Error()))
			render.JSON(writer, request, resp.Error("failed to create rent"))
			return
		}
		ctx, err = rentMaker.SetCassetteStatus(ctx, cassetteId, false)
		if err != nil {
			log.Error("failed to create rent ", slog.String("error", err.Error()))
			render.JSON(writer, request, resp.Error("failed to create rent"))
			return
		}

		ctx, orderId, err := rentMaker.CreateOrder(ctx, customerId)
		if err != nil {
			log.Error("failed to create rent ", slog.String("error", err.Error()))
			render.JSON(writer, request, resp.Error("failed to create rent. Failed to create Order"))
			return
		}

		ctx, err = rentMaker.CreateCassetteInOrder(ctx, cassetteId, orderId, rentCost)
		if err != nil {
			log.Error("failed to create rent ", slog.String("error", err.Error()))
			render.JSON(writer, request, resp.Error("failed to create rent. Failed to create Order"))
			return
		}

		ctx, rentId, err := rentMaker.CreateRent(ctx, customerId, cassetteId, req.RentDays)
		if err != nil {
			log.Error("failed to create rent ", slog.String("error", err.Error()))
			render.JSON(writer, request, resp.Error("failed to create rent."))
			return
		}

		ctx = context.WithValue(ctx, "returnTransaction", nil)

		if _, err := rentMaker.SetCustomerBalance(ctx, customerId, balance-rentCost); err != nil {
			log.Error("failed to create rent ", slog.String("error", err.Error()))
			render.JSON(writer, request, resp.Error("failed to create rent. Error in change customer balance"))
			return
		}

		log.Info("rent created",
			slog.String("create rent with id", rentId),
			slog.String("create rent with id", orderId))

		render.JSON(writer, request, Response{
			Response: resp.Ok(),
			RentId:   rentId,
			OrderId:  orderId,
		})
	}
}
