package makeRent_test

import (
	makeRent "CassetteRental/internal/http-server/handlers/rent/create"
	"CassetteRental/internal/http-server/handlers/rent/create/mocks"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		log *slog.Logger
	}
	tests := []struct {
		name      string
		args      args
		want      http.HandlerFunc
		title     string
		rentDays  int
		respError string
	}{
		{
			name:     "simple positive scenario",
			title:    "Test Film",
			rentDays: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rentMaker := mocks.NewRentMaker(t)
			ctx := context.Background()
			customerId := "14430e91-6a41-4735-a892-f70cede9060d"
			ctx = context.WithValue(ctx, "customerId", customerId)
			filmId := "a46cc6da-3cec-48cf-9959-2ed4ef1c38fa"
			dayPrice := 1
			balance := 10
			rentCost := 5
			cassetteId := "3bb3d807-3790-49c1-adf8-6b7eb3b1cf88"
			orderId := "ead24610-b0e4-4d86-8e49-ae75eb9fd516"
			rentId := "25c33017-3ccb-420e-b07f-aae5b577c040"

			rentMaker.On("GetFilm", mock.Anything, tt.title).
				Return(ctx, filmId, dayPrice, nil)
			rentMaker.On("GetCustomerBalance", mock.Anything, customerId).
				Return(ctx, balance, nil)
			rentMaker.On("FindAvailableCassette", mock.Anything, filmId).
				Return(ctx, cassetteId, nil)
			rentMaker.On("SetCassetteStatus", mock.Anything, cassetteId, false).
				Return(ctx, nil)
			rentMaker.On("CreateOrder", mock.Anything, customerId).
				Return(ctx, orderId, nil)
			rentMaker.On("CreateCassetteInOrder", mock.Anything, cassetteId, orderId, rentCost).
				Return(ctx, nil)
			rentMaker.On("CreateRent", mock.Anything, customerId, cassetteId, tt.rentDays).
				Return(ctx, rentId, nil)
			rentMaker.On("SetCustomerBalance", mock.Anything, customerId, balance-rentCost).
				Return(ctx, nil)

			log := slog.New(slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug},
			))

			handler := makeRent.New(log, rentMaker)
			input := fmt.Sprintf(`{"title":"%s" , "rentDays":%d}`, tt.title, tt.rentDays)

			req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/rent/create", bytes.NewBuffer([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			require.Equal(t, rr.Code, http.StatusOK)
			body := rr.Body.String()

			var resp makeRent.Response
			require.NoError(t, json.Unmarshal([]byte(body), &resp))
			require.Equal(t, tt.respError, resp.Error)
			require.Equal(t, rentId, resp.RentId)
			require.Equal(t, orderId, resp.OrderId)

		})
	}
}
