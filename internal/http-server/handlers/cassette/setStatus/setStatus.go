package setStatus

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
	Available *bool `json:"available" validate:"required"`
}

type Response struct {
	resp.Response
}

type CassetteStatusSetter interface {
	SetCassetteStatus(ctx context.Context, id string, available bool) (context.Context, error)
	GetCassetteStatus(id string) (bool, error)
}

func New(log *slog.Logger, setter CassetteStatusSetter) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handlers/cassette/create/setStatus/New"

		log = slog.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(request.Context())))

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

		id := chi.URLParam(request, "cassetteId")
		if id == "" {
			log.Info("id is empty")
			resp.BadRequest(writer, "Invalid request")
			return
		}
		ctx := context.Background()

		_, err = setter.GetCassetteStatus(id)
		if errors.Is(err, storage.ErrCassetteNotFound) {
			log.Error("failed to set available status to cassette", slog.String("error", err.Error()))
			resp.StatusNotFound(writer, err.Error())
			return
		}
		if err != nil {
			log.Error("failed to set available status to cassette", slog.String("error", err.Error()))
			resp.InternalServerError(writer, err.Error())
		}

		_, err = setter.SetCassetteStatus(ctx, id, *req.Available)
		if err != nil {
			log.Error("failed to set available status to cassette", slog.String("error", err.Error()))
			resp.InternalServerError(writer, err.Error())
			return
		}

		log.Info(
			"available status is changed",
			slog.String("set available status to cassette with id", id),
		)

		render.JSON(writer, request, resp.Ok())
	}

}
