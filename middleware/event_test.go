package middleware

import "testing"

func TestAuthMiddleware(t *testing.T) {
	if err := requestWebSocketApi("http://localhost:3456/message", map[string]interface{}{
		"type": 1,
		"user_id": "5db8066d0d833f36af881c5a",
		"num": 23,
	}); err != nil {
		t.Fatal(err)
	}
}
