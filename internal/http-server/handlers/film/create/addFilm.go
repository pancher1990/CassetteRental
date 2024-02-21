package addFilm

import (
	resp "CassetteRental/internal/lib/api/response"
	"CassetteRental/internal/storage"
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
	DayPrice int    `json:"dayPrice" validate:"required"`
}

type Response struct {
	resp.Response
	Id string `json:"id,omitempty"`
}

type FilmSaver interface {
	AddNewFilm(name string, dayPrice int) (string, error)
	GetFilm(ctx context.Context, title string) (context.Context, string, int, error)
}

func New(log *slog.Logger, saver FilmSaver) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handlers/film/create/addFilm/New"

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
		ctx := context.Background()
		if err := checkUnique(ctx, req.Title, saver); err != nil {
			log.Error("failed to add film need unique title", err.Error())
			render.JSON(writer, request, resp.Error("failed to add film, title already exist"))
			return
		}
		id, err := saver.AddNewFilm(req.Title, req.DayPrice)
		if err != nil {
			log.Error("failed to add film", slog.String("error", err.Error()))

			render.JSON(writer, request, resp.Error("failed to add film"))
			return
		}
		log.Info("film added", slog.String("create film with id", id))

		render.JSON(writer, request, Response{
			Response: resp.Ok(),
			Id:       id,
		})
	}
}

func checkUnique(ctx context.Context, title string, filmGetter FilmSaver) error {
	_, _, _, err := filmGetter.GetFilm(ctx, title)
	if errors.Is(err, storage.ErrFilmNotFound) {
		return nil
	} else {
		return errors.New("film title already exist")
	}
}
