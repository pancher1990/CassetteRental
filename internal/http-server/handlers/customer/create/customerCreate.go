package customerCreate

import (
	resp "CassetteRental/internal/lib/api/response"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type Request struct {
	Name string `json:"name" validate:"required"`
	//IsActive bool   `json:"isActive,omitempty" validate:"boolean"`
	Balance int `json:"balance,omitempty"`
}

type Response struct {
	resp.Response
	Id string `json:"id,omitempty"`
}

type CustomerSaver interface {
	AddNewCustomer(name string, isActive bool, balance int) (string, error)
}

func New(log *slog.Logger, saver CustomerSaver) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handlers/customer/create/customerCreate/New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(request.Context())),
		)
		var req Request

		err := render.DecodeJSON(request.Body, &req)

		if err != nil {
			log.Error("Failed to decode request body ", err.Error())
			render.JSON(writer, request, resp.Error("Failed to decode request"))
			return
		}

		log.Info("request body decoded ", slog.Any("request", req))
		if err = validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("Invalid request", err.Error())
			render.JSON(writer, request, resp.ValidationError(validateErr))
			return
		}

		id, err := saver.AddNewCustomer(req.Name, true, req.Balance)
		if err != nil {
			log.Error("failed to add customer", err.Error())

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
