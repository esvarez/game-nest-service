package api

import (
	"flag"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/esvarez/game-nest-service/api/handler"
	"github.com/esvarez/game-nest-service/internal/config"
	"github.com/esvarez/game-nest-service/internal/storage"
	"github.com/esvarez/game-nest-service/service/game"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

func TestIntegration(t *testing.T) {
	handler := setUp()
	srv := httptest.NewServer(handler)
	defer srv.Close()
	// ids := make([]string, 0)
	t.Run("GameHandler", func(t *testing.T) {
		tests := map[string]struct {
			url            string
			id             string
			body           io.Reader
			expectedStatus int
			expectedBody   string
			httpMethod     string
		}{
			"01_CreateGame fail bad request": {
				url:            "/api/v1/game",
				body:           strings.NewReader(`{"name": "Test", "description": "Test"}`),
				expectedBody:   `{"status_code":400,"type":"invalid_entity","message":"Key: 'Game.MinPlayers' Error:Field validation for 'MinPlayers' failed on the 'required' tag\nKey: 'Game.MaxPlayers' Error:Field validation for 'MaxPlayers' failed on the 'required' tag\nKey: 'Game.Duration' Error:Field validation for 'Duration' failed on the 'required' tag: game item invalid entity"}`,
				httpMethod:     http.MethodPost,
				expectedStatus: http.StatusBadRequest,
			},
			"01_CreateGame success": {
				url:            "/api/v1/game",
				body:           strings.NewReader(`{"name":"Catan Test","description":"Catan is a board game","min_players":2,"max_players":4,"duration":30}`),
				httpMethod:     http.MethodPost,
				expectedStatus: http.StatusCreated,
			},
			"02_Get all games success": {
				url:            "/api/v1/game",
				expectedBody:   `"name":"Catan Test","min_players":2,"max_players":4,"description":"Catan is a board game","duration":30`,
				httpMethod:     http.MethodGet,
				expectedStatus: http.StatusOK,
			},
			"02_Get game success": {
				url:            "/api/v1/game/no-exist",
				expectedBody:   `"name":"Catan","min_players":2,"max_players":4,"description":"Catan is a board game","duration":30}`,
				id:             "/",
				httpMethod:     http.MethodGet,
				expectedStatus: http.StatusNotFound,
			},
			"03_Update game success": {
				url:            "/api/v1/game/{id}",
				id:             "1",
				body:           strings.NewReader(`{"name":"test"}`),
				httpMethod:     http.MethodPut,
				expectedStatus: http.StatusOK,
			},
			"04_Delete game success": {
				url:            "/api/v1/game/{id}",
				id:             "1",
				httpMethod:     http.MethodDelete,
				expectedStatus: http.StatusOK,
			},
		}

		keys := make([]string, 0)
		for key := range tests {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, name := range keys {
			tc := tests[name]
			t.Run(name, func(t *testing.T) {
				switch {
				case tc.httpMethod == http.MethodGet:
					assert.HTTPBodyContains(t, handler.ServeHTTP, tc.httpMethod, tc.url, nil, tc.expectedBody)
				default:
					req, err := http.NewRequest(tc.httpMethod, srv.URL+tc.url, tc.body)
					if err != nil {
						t.Fatal(err)
					}

					res, err := http.DefaultClient.Do(req)
					if err != nil {
						t.Fatal(err)
					}

					assert.Equal(t, tc.expectedStatus, res.StatusCode)
					if tc.expectedBody != "" {
						body, err := ioutil.ReadAll(res.Body)
						if err != nil {
							t.Fatal(err)
						}
						assert.Contains(t, tc.expectedBody, strings.Trim(string(body), "\n"))
						/*
							if res.StatusCode < 300 {
								var games []presenter.GameResponse
								json.Unmarshal(body, &games)
								ids = append(ids, games[0].ID)
							}
						*/
					}
				}
			})
		}
	})
}

func setUp() http.Handler {
	var pathFile string
	flag.StringVar(&pathFile, "public-config-file",
		"./test_file/config.yml", "Path to public config file")

	var (
		v      = validator.New()
		conf   = config.LoadConfiguration(pathFile, v)
		l      = config.CreateLogger(conf)
		r      = mux.NewRouter()
		client = storage.CreateDynamoClient(conf)

		gameClient = storage.NewGameClient(client, l)

		gameService = game.NewService(gameClient, l, v)

		gameHandler = handler.NewGameHandler(gameService, l)
	)
	r = r.PathPrefix("/api/v1").Subrouter()
	handler.MakeGameHandler(r, gameHandler)
	return r
}
