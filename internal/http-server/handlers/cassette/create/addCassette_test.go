package addCassette_test

import (
	"CassetteRental/internal/http-server/handlers/cassette/create"
	"CassetteRental/internal/http-server/handlers/cassette/create/mocks"
	"bytes"
	"encoding/json"
	"fmt"
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
		name        string
		args        args
		want        http.HandlerFunc
		count       int
		id          string
		expectedIds []string
		respError   string
	}{
		{
			name:        "simple positive scenario",
			id:          "a46cc6da-3cec-48cf-9959-2ed4ef1c38fa",
			count:       1,
			expectedIds: []string{"3bb3d807-3790-49c1-adf8-6b7eb3b1cf88"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cassetteSaver := mocks.NewCassetteSaver(t)
			cassetteSaver.On("AddNewCassette", tt.id).
				Return("3bb3d807-3790-49c1-adf8-6b7eb3b1cf88", nil).Once()
			cassetteSaver.On("GetFilmById", tt.id).
				Return("какой-то id", 123, nil).Once()
			log := slog.New(slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug},
			))

			handler := addCassette.New(log, cassetteSaver)
			input := fmt.Sprintf(`{"count":%d, "id":"%s"}`, tt.count, tt.id)

			req, err := http.NewRequest(http.MethodPost, "/cassette/add", bytes.NewBuffer([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			require.Equal(t, rr.Code, http.StatusOK)
			body := rr.Body.String()

			var resp addCassette.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))
			require.Equal(t, tt.respError, resp.Error)

			require.Equal(t, tt.expectedIds, resp.Ids, "The response Ids do not match the expected ones")
		})
	}
}
