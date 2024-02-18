package findFilm

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
	Title string `json:"title" validate:"required"`
	//DayPrice int    `json:"dayPrice" validate:"required"`
}

type Response struct {
	resp.Response
	Id      string `json:"id,omitempty"`
	DayCost int    `json:"dayCost,omitempty"`
}

type FilmFinder interface {
	GetFilm(ctx context.Context, title string) (context.Context, string, int, error)
}

func New(log *slog.Logger, finder FilmFinder) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var op = "handlers/film/find/findFilm/New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(request.Context())),
		)
		var req Request

		if err := render.DecodeJSON(request.Body, &req); err != nil {
			log.Error("Failed to decode request body ", slog.String("error", err.Error()))
			render.JSON(writer, request, resp.Error("Failed to decode request"))

			return
		}

		log.Info("request body decoded ", slog.Any("request", req))
		if err := validator.New().Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)
			log.Error("Invalid request", slog.String("error", err.Error()))
			render.JSON(writer, request, resp.ValidationError(validateErr))
			return
		}

		ctx := context.Background()
		_, id, dayCost, err := finder.GetFilm(ctx, req.Title)
		if err != nil {
			log.Error("failed to get film", slog.String("error", err.Error()))

			render.JSON(writer, request, resp.Error("failed to get film"))
			return
		}
		log.Info("film was taken", slog.String("get film with id", id))

		render.JSON(writer, request, Response{
			Response: resp.Ok(),
			Id:       id,
			DayCost:  dayCost,
		})
	}
}
