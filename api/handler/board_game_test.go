package handler

import (
	"testing"
)

//go:generate mockery --name UseCase --dir ./service/boardgame/service --outpkg mocks --output ./api/handler/mocks --case=underscore

func TestRouting(t *testing.T) {
	/*
		r := mux.NewRouter()
		r = r.PathPrefix("/api/v1").Subrouter()
		bgService := &mocks.UseCase{}
		l := &logrus.Logger{}
		handler := NewBoardGameHandler(bgService, l)
		MakeGameHandler(r, handler)
		srv := httptest.NewServer(r)
		defer srv.Close()
		tests := map[string]struct {
			url            string
			body           io.Reader
			expectedStatus int
			expectedBody   string
			httpMethod     string
			mockSetup      func(bgService *mocks.UseCase)
		}{
			"get games": {
				url:        "/api/v1/game",
				httpMethod: "GET",
			},
		}

		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				fmt.Println(srv.URL + test.url)
				req, err := http.NewRequest(test.httpMethod, srv.URL+test.url, nil)
				if err != nil {
					t.Fatal(err)
				}
				res, err := http.DefaultClient.Do(req)
				if err != nil {
					t.Fatal(err)
				}
				defer res.Body.Close()
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, 200, res.StatusCode)
				assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
				assert.Equal(t, "", string(body))
			})
		}
	*/
}
