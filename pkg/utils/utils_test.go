package utils

import (
	"main/pkg/models"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenCreation(t *testing.T) {
	os.Setenv("SECRET_KEY", "123")
	permission := models.Permission{}
	permission.SetDefaultPermissions()
	user := models.User{
		Username:   "tester",
		Permission: permission,
	}
	token, err := CreateToken(user)
	if err != nil {
		t.Fail()
	}
	assert.NotEmpty(t, token)

	os.Unsetenv("SECRET_KEY")
}

func TestCheckPassword_ShouldBeSuccess(t *testing.T) {
	passwd := "test-password"
	hash, err := HashPwd(passwd)
	if err != nil {
		t.Fail()
	}
	err = CheckPassword(hash, passwd)
	assert.Nil(t, err)
}

func TestCheckPassword_ShouldBeError(t *testing.T) {
	passwd := "test-password"
	hash, err := HashPwd(passwd)
	if err != nil {
		t.Fail()
	}
	err = CheckPassword(hash, "wrong-password")
	assert.NotNil(t, err)
}

func TestCheckPassword_ShouldBeErrorOnNullValue(t *testing.T) {
	passwd := ""
	hash, err := HashPwd(passwd)
	if err != nil {
		t.Fail()
	}
	err = CheckPassword(hash, passwd)
	assert.NotNil(t, err)
}
