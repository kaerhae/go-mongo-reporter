package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func TestAuthenticate_ShouldReturnErrIfNoAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	Authenticate(c)
	data, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fail()
	}
	if w.Code != 401 || string(data) != `{"message":"401 Unauthorized"}` {
		t.Fail()
	}
}

func TestAuthenticate_ShouldReturnErrIfNoSecret(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	c.Request.Header.Add("Authorization", "blaah.blaah.blaah")
	Authenticate(c)

	if w.Code != 500 {
		t.Fail()
	}
}

func TestAuthenticate_ShouldReturn200IfCorrectTokenKey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("SECRET_KEY", "blaah")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}

	mockToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp": time.Now().Add(time.Minute * 30).Unix(),
		})
	mockSecret := []byte("blaah")
	tokenstring, _ := mockToken.SignedString(mockSecret)
	c.Request.Header.Add("Authorization", tokenstring)
	Authenticate(c)

	if w.Code != 200 {
		t.Fail()
	}
	os.Unsetenv("SECRET_KEY")
}

func TestAuthenticate_ShouldReturnErrorIfIncorrectToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("SECRET_KEY", "blaah")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}

	mockToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp": time.Now().Add(time.Minute * 30).Unix(),
		})
	mockSecret := []byte("DEFINITELYWRONG")
	tokenstring, _ := mockToken.SignedString(mockSecret)
	c.Request.Header.Add("Authorization", tokenstring)
	Authenticate(c)

	if w.Code != 400 {
		t.Fail()
	}
	os.Unsetenv("SECRET_KEY")
}
