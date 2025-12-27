package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/jwtauth"

	"github.com/akarashov/gofermart/internal/models"
	"github.com/akarashov/gofermart/internal/services"
)

type mockOrderService struct {
	uploadFunc func(ctx context.Context, userID, orderNumber string) error
	getFunc    func(ctx context.Context, userID string) ([]models.OrderResponse, error)
}

func (m *mockOrderService) UploadOrder(ctx context.Context, userID, orderNumber string) error {
	if m.uploadFunc != nil {
		return m.uploadFunc(ctx, userID, orderNumber)
	}
	return nil
}

func (m *mockOrderService) GetUserOrders(ctx context.Context, userID string) ([]models.OrderResponse, error) {
	if m.getFunc != nil {
		return m.getFunc(ctx, userID)
	}
	return nil, nil
}

func TestGetUserOrdersHandler(t *testing.T) {
	// helper to create request with JWT context
	makeReqWithToken := func() *http.Request {
		ja := jwtauth.New("HS256", []byte("testkey"), nil)
		token, _, err := ja.Encode(map[string]interface{}{"sub": "user1"})
		if err != nil {
			t.Fatalf("encode token: %v", err)
		}
		ctx := jwtauth.NewContext(context.Background(), token, nil)
		req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
		return req
	}

	tests := []struct {
		name       string
		req        *http.Request
		svcRet     []models.OrderResponse
		svcErr     error
		wantStatus int
		wantBody   string
	}{
		{"NoAuth", httptest.NewRequest(http.MethodGet, "/", nil), nil, nil, http.StatusUnauthorized, ""},
		{"NoOrders", makeReqWithToken(), nil, services.ErrNoOrdersForUser, http.StatusNoContent, ""},
		{"Success", makeReqWithToken(), []models.OrderResponse{{Number: "79927398713", Status: models.OrderStatusNew, Accrual: float32(0), UploadedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)}}, nil, http.StatusOK, "79927398713"},
		{"InternalError", makeReqWithToken(), nil, errors.New("boom"), http.StatusInternalServerError, ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := &mockOrderService{getFunc: func(ctx context.Context, userID string) ([]models.OrderResponse, error) {
				return tc.svcRet, tc.svcErr
			}}

			handler := NewOrdersHandler(svc)
			rr := httptest.NewRecorder()

			req := tc.req
			handler.GetOrders(rr, req)

			if rr.Code != tc.wantStatus {
				t.Fatalf("want status %d, got %d, body: %s", tc.wantStatus, rr.Code, rr.Body.String())
			}

			if tc.wantBody != "" && !strings.Contains(rr.Body.String(), tc.wantBody) {
				t.Fatalf("want body to contain %q, got %q", tc.wantBody, rr.Body.String())
			}
		})
	}
}

func TestUploadOrderHandler(t *testing.T) {
	// helper to create request with JWT context
	makeReqWithToken := func(body string) *http.Request {
		ja := jwtauth.New("HS256", []byte("testkey"), nil)
		token, _, err := ja.Encode(map[string]interface{}{"sub": "user1"})
		if err != nil {
			t.Fatalf("encode token: %v", err)
		}
		ctx := jwtauth.NewContext(context.Background(), token, nil)
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body)).WithContext(ctx)
		req.Header.Set("Content-Type", "text/plain")
		return req
	}

	tests := []struct {
		name       string
		req        *http.Request
		svcErr     error
		wantStatus int
	}{
		{"BadContentType", httptest.NewRequest(http.MethodPost, "/", nil), nil, http.StatusBadRequest},
		{"NoAuth", func() *http.Request {
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("123"))
			r.Header.Set("Content-Type", "text/plain")
			return r
		}(), nil, http.StatusUnauthorized},
		{"EmptyBody", makeReqWithToken(""), nil, http.StatusBadRequest},
		{"Success", makeReqWithToken("79927398713"), nil, http.StatusAccepted},
		{"AlreadyByUser", makeReqWithToken("79927398713"), services.ErrOrderAlreadyUploadedByUser, http.StatusOK},
		{"AlreadyByAnother", makeReqWithToken("79927398713"), services.ErrOrderAlreadyUploadedByAnotherUser, http.StatusConflict},
		{"InvalidNumber", makeReqWithToken("abc"), services.ErrInvalidOrderNumber, http.StatusUnprocessableEntity},
		{"InternalError", makeReqWithToken("79927398713"), errors.New("boom"), http.StatusInternalServerError},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var svc mockOrderService
			svc.uploadFunc = func(ctx context.Context, userID, orderNumber string) error {
				return tc.svcErr
			}

			handler := NewOrdersHandler(&svc)
			rr := httptest.NewRecorder()

			req := tc.req
			handler.UploadOrder(rr, req)

			if rr.Code != tc.wantStatus {
				t.Fatalf("want status %d, got %d, body: %s", tc.wantStatus, rr.Code, rr.Body.String())
			}
		})
	}
}
