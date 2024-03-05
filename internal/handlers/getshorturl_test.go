package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rinem/url-shortener-go/store"
	storemock "github.com/rinem/url-shortener-go/store/mock"
	"github.com/stretchr/testify/assert"
)

func TestGetShortURL(t *testing.T) {
	testCases := []struct {
		name                      string
		url                       string
		expectedStatusCode        int
		expectedGetShortURLParams string
		getShortURLResult         *store.ShortURL
	}{
		{
			name:                      "success",
			url:                       "/123",
			expectedStatusCode:        http.StatusMovedPermanently,
			expectedGetShortURLParams: "123",
			getShortURLResult: &store.ShortURL{
				Destination: "http://google.com",
			},
		},
		{
			name:                      "fail with 404",
			url:                       "/123",
			expectedStatusCode:        http.StatusNotFound,
			expectedGetShortURLParams: "123",
			getShortURLResult:         nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			assert := assert.New(t)
			shortURLStoreMock := &storemock.MockShortURLStore{}

			shortURLStoreMock.On("GetShortURLBySlug", tc.expectedGetShortURLParams).Return(tc.getShortURLResult, nil)

			handler := NewGetShortURLHandler(GetShortURLHandlerParams{
				ShortURLStore: shortURLStoreMock,
			})

			request := httptest.NewRequest("GET", tc.url, nil)
			responseRecorder := httptest.NewRecorder()

			handler.ServeHTTP(responseRecorder, request)

			response := responseRecorder.Result()
			defer response.Body.Close()

			assert.Equal(tc.expectedStatusCode, response.StatusCode)

			shortURLStoreMock.AssertExpectations(t)
		})
	}

}
