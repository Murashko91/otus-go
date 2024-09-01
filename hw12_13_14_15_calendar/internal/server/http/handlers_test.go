package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

const (
	contentTypeJSON = "application/json"
	userID          = 1
)

func TestEventHTTPApi(t *testing.T) {
	logg := logger.New("info")
	db := memorystorage.New()
	if err := db.Connect(); err != nil {
		t.Errorf("storage connect error %s", err.Error())
	}
	calendarApp := app.New(logg, db)
	h := Handler{
		app: calendarApp,
	}
	server := httptest.NewServer(http.HandlerFunc(h.eventHandler))
	urlServer, _ := url.Parse(server.URL)
	urlServer.Path = "event"

	t.Run("events API", func(t *testing.T) {
		// insert events
		params := urlServer.Query()
		params.Add(app.UserIDKey, "1")
		urlServer.RawQuery = params.Encode()

		for i := 0; i < 30; i++ {
			requestJSON := storage.Event{
				Title:     fmt.Sprintf("test %d", i),
				Descr:     fmt.Sprintf("test %d", i),
				StartDate: time.Now(),
				EndDate:   time.Now(),
				UserID:    1,
			}

			reqBody, err := json.Marshal(requestJSON)
			require.NoError(t, err, "not expected body marshal error")

			resp, err := doHTTPCall(urlServer.String(), http.MethodPost, bytes.NewReader(reqBody))
			require.NoError(t, err, "not expected request error")
			require.Equal(t, resp.StatusCode, http.StatusOK)

			resBody, err := io.ReadAll(resp.Body)
			defer resp.Body.Close()
			require.NoError(t, err, "not expected read body error")

			var resJSON storage.Event
			err = json.Unmarshal(resBody, &resJSON)
			require.NoError(t, err, "response unmarshal error")
			require.Equal(t, resJSON.UserID, userID, "not expected userId")
			require.Equal(t, resJSON.Title, requestJSON.Title, "not expected Title")
			require.NotNil(t, resJSON.ID, "event id not populated")
		}

		// update events

		for i := 1; i < 31; i++ {
			requestJSON := storage.Event{
				ID:        i,
				Title:     fmt.Sprintf("test updated %d", i),
				Descr:     fmt.Sprintf("test updated %d", i),
				StartDate: time.Now().Add(time.Second).Add(time.Hour * 24 * time.Duration(i)),
				EndDate:   time.Now().Add(time.Hour * 24 * time.Duration(i)),
				UserID:    1,
			}

			reqBody, err := json.Marshal(requestJSON)
			require.NoError(t, err, "not expected body marshal error")
			resp, err := doHTTPCall(urlServer.String(), http.MethodPatch, bytes.NewReader(reqBody))
			require.NoError(t, err, "not expected request error")
			require.Equal(t, resp.StatusCode, http.StatusOK)

			defer resp.Body.Close()

			require.NoError(t, err, "not expected request error")

			resBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err, "not expected read body error")

			var resJSON storage.Event
			err = json.Unmarshal(resBody, &resJSON)
			require.NoError(t, err, "response unmarshal error")
			require.Equal(t, resJSON.UserID, userID, "not expected userId")
			require.Equal(t, resJSON.Title, requestJSON.Title, "not expected Title")
			require.NotNil(t, resJSON.ID, "event id not populated")
		}

		// get events

		requests := make([]string, 0, 3)

		params.Add("start_date", time.Now().Add(time.Hour*3).Format(time.RFC3339))
		params.Add("duration", "day")
		urlServer.RawQuery = params.Encode()
		urlServer.Path = "event"

		requests = append(requests, urlServer.String())
		params.Set("duration", "week")
		urlServer.RawQuery = params.Encode()
		requests = append(requests, urlServer.String())

		params.Set("duration", "month")
		urlServer.RawQuery = params.Encode()

		requests = append(requests, urlServer.String())

		for _, urlReq := range requests {
			resp, err := doHTTPCall(urlReq, http.MethodGet, nil)
			require.NoError(t, err, "not expected error")
			var resJSON []storage.Event

			resBody, err := io.ReadAll(resp.Body)
			defer resp.Body.Close()
			require.NoError(t, err, "not expected read body error")
			require.Equal(t, resp.StatusCode, http.StatusOK)

			err = json.Unmarshal(resBody, &resJSON)
			require.NoError(t, err, "response unmarshal error")

			rParams, _ := url.Parse(urlReq)
			switch rParams.Query().Get("duration") {
			case "day":
				require.Equal(t, len(resJSON), 1)
			case "week":
				require.Equal(t, len(resJSON), 7)
			case "month":
				require.Equal(t, len(resJSON), 30)
			}
		}

		// test delete events

		testDeleteEvents(t, urlServer)
	})
}

func testDeleteEvents(t *testing.T, urlServer *url.URL) {
	t.Helper()
	for i := 1; i < 31; i++ {
		requestJSON := storage.Event{
			ID:     i,
			UserID: 1,
		}

		reqBody, err := json.Marshal(requestJSON)
		require.NoError(t, err, "not expected body marshal error")
		resp, err := doHTTPCall(urlServer.String(), http.MethodDelete, bytes.NewReader(reqBody))

		require.NoError(t, err, "not expected request error")

		defer resp.Body.Close()

		require.Equal(t, resp.StatusCode, http.StatusOK)
	}
}

func doHTTPCall(url string, method string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(context.Background(), method, url, body)
	if err != nil {
		return nil, err
	}

	// set content-type header to JSON
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// create HTTP client and execute request
	client := &http.Client{}

	return client.Do(req)
}
