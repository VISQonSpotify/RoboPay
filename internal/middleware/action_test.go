package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	x402 "github.com/x402-foundation/x402/go"
)

func TestCaptureActionPayloadMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(CaptureActionPayloadMiddleware())

	var capturedFromBody []byte
	var capturedFromContext []byte

	router.POST("/action", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		capturedFromBody = body
		capturedFromContext = ActionPayloadFromContext(c)
		c.Status(http.StatusOK)
	})

	reqBody := []byte(`{"command":"forward","speed":2}`)
	req := httptest.NewRequest(http.MethodPost, "/action", bytes.NewReader(reqBody))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.Code)
	}
	if !bytes.Equal(capturedFromBody, reqBody) {
		t.Fatalf("expected handler body %q, got %q", string(reqBody), string(capturedFromBody))
	}
	if !bytes.Equal(capturedFromContext, reqBody) {
		t.Fatalf("expected context body %q, got %q", string(reqBody), string(capturedFromContext))
	}
}

func TestBuildActionZenohEvent(t *testing.T) {
	payload := []byte(`{"command":"dock"}`)
	settle := &x402.SettleResponse{
		Transaction: "0xabc",
		Network:     "eip155:84532",
		Payer:       "0x123",
	}

	eventBytes, err := BuildActionZenohEvent(payload, settle)
	if err != nil {
		t.Fatalf("expected no error building event, got %v", err)
	}

	var event map[string]interface{}
	if err := json.Unmarshal(eventBytes, &event); err != nil {
		t.Fatalf("expected valid JSON event, got %v", err)
	}

	payloadField, ok := event["payload"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected payload object, got %T", event["payload"])
	}
	if payloadField["command"] != "dock" {
		t.Fatalf("expected command 'dock', got %v", payloadField["command"])
	}

	txField, ok := event["transaction_details"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected transaction details object, got %T", event["transaction_details"])
	}
	if txField["transaction"] != "0xabc" {
		t.Fatalf("expected transaction 0xabc, got %v", txField["transaction"])
	}
	if txField["network"] != "eip155:84532" {
		t.Fatalf("expected network eip155:84532, got %v", txField["network"])
	}
	if txField["payer"] != "0x123" {
		t.Fatalf("expected payer 0x123, got %v", txField["payer"])
	}
}
