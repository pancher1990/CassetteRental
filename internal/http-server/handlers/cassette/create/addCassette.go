package addCassette

import (
	resp "CassetteRental/internal/lib/api/response"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"strings"
)

type Request struct {
	Id    string `json:"id" validate:"required"`
	Count int    `json:"count" validate:"required"`
}

type Response struct {
	resp.Response
	Ids []string `json:"ids,omitempty"`
}

type CassetteSaver interface {
	AddNewCassette(filmId string) (string, error)
}

func New(log *slog.Logger, saver CassetteSaver) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handlers/cassette/create/addCassette/New"

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
		var ids []string
		for i := 0; i < req.Count; i++ {
			id, err := saver.AddNewCassette(req.Id)
			if err != nil {
				log.Error("failed to add cassette", slog.String("error", err.Error()))

				render.JSON(writer, request, resp.Error("failed to add cassette"))
				return
			}
			ids = append(ids, id)
		}

		log.Info("cassette added", slog.String("create cassette with ids", strings.Join(ids, ", ")))

		render.JSON(writer, request, Response{
			Response: resp.Ok(),
			Ids:      ids,
		})
	}
}
