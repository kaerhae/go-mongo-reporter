package middleware

import (
	"io"
	"main/pkg/models"
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

	if w.Code != 400 {
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
	expirationTime := time.Now().Add(5 * time.Minute)
	mockClaims := models.Claims{
		Username: "test",
		AppRole:  models.Admin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	mockToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		mockClaims)
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

	expirationTime := time.Now().Add(5 * time.Minute)
	mockClaims := models.Claims{
		Username: "test",
		AppRole:  models.Admin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	mockToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		mockClaims)
	mockSecret := []byte("DEFINITELYWRONG")
	tokenstring, _ := mockToken.SignedString(mockSecret)
	c.Request.Header.Add("Authorization", tokenstring)
	Authenticate(c)

	if w.Code != 400 {
		t.Fail()
	}
	os.Unsetenv("SECRET_KEY")
}

func TestAuthenticateAdmin_ShouldReturnErrIfNoAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	AuthenticateAdmin(c)
	data, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fail()
	}
	if w.Code != 401 || string(data) != `{"message":"401 Unauthorized"}` {
		t.Fail()
	}
}

func TestAuthenticateAdmin_ShouldReturnErrIfNoSecret(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	c.Request.Header.Add("Authorization", "blaah.blaah.blaah")
	AuthenticateAdmin(c)

	if w.Code != 400 {
		t.Fail()
	}
}

func TestAuthenticateAdmin_ShouldReturn200IfCorrectTokenKey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("SECRET_KEY", "blaah")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	expirationTime := time.Now().Add(5 * time.Minute)
	mockClaims := models.Claims{
		Username: "test",
		AppRole:  models.Admin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	mockToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		mockClaims)
	mockSecret := []byte("blaah")
	tokenstring, _ := mockToken.SignedString(mockSecret)
	c.Request.Header.Add("Authorization", tokenstring)
	AuthenticateAdmin(c)

	if w.Code != 200 {
		t.Fail()
	}
	os.Unsetenv("SECRET_KEY")
}

func TestAuthenticateAdmin_ShouldReturnErrorIfIncorrectToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("SECRET_KEY", "blaah")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	expirationTime := time.Now().Add(5 * time.Minute)
	mockClaims := models.Claims{
		Username: "test",
		AppRole:  models.Admin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	mockToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		mockClaims)
	mockSecret := []byte("DEFINITELYWRONG")
	tokenstring, _ := mockToken.SignedString(mockSecret)
	c.Request.Header.Add("Authorization", tokenstring)
	AuthenticateAdmin(c)

	if w.Code != 400 {
		t.Fail()
	}
	os.Unsetenv("SECRET_KEY")
}

func TestAuthenticateAdmin_ShouldReturn403IfIncorrectRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("SECRET_KEY", "blaah")
	tests := []struct {
		input  models.Role
		output int
	}{
		{models.Admin, 200},
		{models.Maintainer, 403},
		{models.Creator, 403},
		{models.Guest, 403},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = &http.Request{
			Header: make(http.Header),
			URL:    &url.URL{},
		}
		expirationTime := time.Now().Add(5 * time.Minute)
		mockClaims := models.Claims{
			Username: "test",
			AppRole:  test.input,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
				IssuedAt:  time.Now().Unix(),
			},
		}

		mockToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
			mockClaims)
		mockSecret := []byte("blaah")
		tokenstring, _ := mockToken.SignedString(mockSecret)
		c.Request.Header.Add("Authorization", tokenstring)
		AuthenticateAdmin(c)

		if w.Code != test.output {
			t.Fail()
		}
	}
	os.Unsetenv("SECRET_KEY")
}
