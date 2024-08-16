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
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		Method: http.MethodGet,
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	p := models.Permission{}
	p.SetDefaultPermissions()

	c.Request.Header.Add("Authorization", createMockToken(p))
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
	p := models.Permission{}
	p.SetDefaultPermissions()
	expirationTime := time.Now().Add(5 * time.Minute)
	mockClaims := models.Claims{
		Username:    "test",
		Permissions: p,
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
	p := models.Permission{}
	p.Admin = true

	c.Request.Header.Add("Authorization", createMockToken(p))
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
	p := models.Permission{}
	p.SetDefaultPermissions()
	expirationTime := time.Now().Add(5 * time.Minute)
	mockClaims := models.Claims{
		Username:    "test",
		Permissions: p,
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

func TestAuthenticateAdmin_ShouldReturn403IfNotAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("SECRET_KEY", "blaah")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	p := models.Permission{
		Admin: false,
		Write: true,
		Read:  true,
	}

	c.Request.Header.Add("Authorization", createMockToken(p))
	AuthenticateAdmin(c)

	if w.Code != 403 {
		t.Fail()
	}
	os.Unsetenv("SECRET_KEY")
}

func TestAuthenticate_ShouldAddUserIdAndNotAdminToStoreIfSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("SECRET_KEY", "blaah")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Method: http.MethodGet,
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	p := models.Permission{}
	p.SetDefaultPermissions()

	c.Request.Header.Add("Authorization", createMockToken(p))
	Authenticate(c)

	if w.Code != 200 {
		t.Fail()
	}
	if c.GetString("userId") != "123456789012345678901234" {
		t.Fail()
	}
	if c.GetBool("isAdmin") != false {
		t.Fail()
	}
	os.Unsetenv("SECRET_KEY")
}

func TestAuthenticate_ShouldAddUserIdAndIsAdminToStoreIfSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("SECRET_KEY", "blaah")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	p := models.Permission{}
	p.Admin = true

	c.Request.Header.Add("Authorization", createMockToken(p))
	AuthenticateAdmin(c)

	if w.Code != 200 {
		t.Fail()
	}
	if c.GetString("userId") != "123456789012345678901234" {
		t.Fail()
	}
	if c.GetBool("isAdmin") != true {
		t.Fail()
	}
	os.Unsetenv("SECRET_KEY")
}

func TestAuthenticate_ShouldCheckPermissionsByMethod(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("SECRET_KEY", "blaah")

	tests := []struct {
		method       string
		permission   models.Permission
		expectedCode int
	}{
		{http.MethodGet, models.Permission{Admin: false, Write: false, Read: true}, 200},
		{http.MethodGet, models.Permission{Admin: false, Write: false, Read: false}, 401},

		{http.MethodPost, models.Permission{Admin: false, Write: false}, 401},
		{http.MethodPost, models.Permission{Admin: false, Write: true}, 200},
		{http.MethodPut, models.Permission{Admin: false, Write: false}, 401},
		{http.MethodPut, models.Permission{Admin: false, Write: true}, 200},
		{http.MethodDelete, models.Permission{Admin: false, Write: false}, 401},
		{http.MethodDelete, models.Permission{Admin: false, Write: true}, 200},
	}

	for _, v := range tests {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = &http.Request{
			/* method setup */
			Method: v.method,
			Header: make(http.Header),
			URL:    &url.URL{},
		}

		/* permission setup */
		c.Request.Header.Add("Authorization", createMockToken(v.permission))
		Authenticate(c)

		/* expected setup */
		if w.Code != v.expectedCode {
			t.Fail()
		}

	}
	os.Unsetenv("SECRET_KEY")
}

func TestAuthenticateAdmin_ShouldAddUserIdAndIsAdminToStoreIfSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("SECRET_KEY", "blaah")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	p := models.Permission{}
	p.Admin = true

	c.Request.Header.Add("Authorization", createMockToken(p))
	AuthenticateAdmin(c)

	if w.Code != 200 {
		t.Fail()
	}
	if c.GetString("userId") != "123456789012345678901234" {
		t.Fail()
	}
	if c.GetBool("isAdmin") != true {
		t.Fail()
	}
	os.Unsetenv("SECRET_KEY")
}

func createMockToken(p models.Permission) string {

	expirationTime := time.Now().Add(5 * time.Minute)
	id, _ := primitive.ObjectIDFromHex("123456789012345678901234")
	mockClaims := models.Claims{
		UserID:      id,
		Username:    "test",
		Permissions: p,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	mockToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		mockClaims)
	mockSecret := []byte("blaah")
	tokenstring, _ := mockToken.SignedString(mockSecret)
	return tokenstring
}
