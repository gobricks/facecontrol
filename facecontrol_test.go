package facecontrol

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func getFacecontrol() *Facecontrol {
	fc, _ := New(Config{
		RunAt:     ":50058",
		JwtSecret: "ShwiftyJWTSecret",
		Validator: validatorTestFunction,
	})

	return fc
}

func validatorTestFunction(r *http.Request) (Payload, error) {
	return map[string]interface{}{
		"is_admin": true,
		"can_edit": []string{"posts", "comments"},
	}, nil
}

func TestTokenIssue(t *testing.T) {
	req, err := http.NewRequest("POST", "/issue?username=admin&password=12345", nil)
	if err != nil {
		t.Fatal(err)
	}

	fc := getFacecontrol()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(fc.issueToken)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if rr.Body.Len() == 0 {
		t.Errorf("Expected JWT token in response body")
	}
}

func TestTokenValidate(t *testing.T) {
	testTokenString := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjp7ImNhbl9lZGl0IjpbInBvc3RzIiwiY29tbWVudHMiXSwiaXNfYWRtaW4iOnRydWV9LCJpYXQiOjE1MDI3OTg0MzYsImlzcyI6ImZhY2Vjb250cm9sIn0.VoafopbLAZUzmf2FfkafGzqDtIqC4XpqmHbFBGvHihkgxaiGHSTZlWH83vPRLQW0yxkyqwJJU0rBBmI-pkFMkg"

	req, err := http.NewRequest("GET", "/validate", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+testTokenString)

	fc := getFacecontrol()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(fc.validateToken)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if rr.Body.Len() == 0 {
		t.Errorf("Expected decoded token payload in response body")
	}
}
