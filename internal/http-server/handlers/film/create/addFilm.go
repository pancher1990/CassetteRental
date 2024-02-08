package addFilm

import (
	resp "CassetteRental/internal/lib/api/response"
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

		id, err := saver.AddNewFilm(req.Title, req.DayPrice)
		if err != nil {
			log.Error("failed to add film", err.Error())

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
