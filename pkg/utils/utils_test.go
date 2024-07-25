package utils

import (
	"main/pkg/models"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenCreation(t *testing.T) {
	os.Setenv("SECRET_KEY", "123")
	user := models.User{
		Username: "tester",
		AppRole:  "admin",
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

func TestRoleDeterminationShouldReturnCorrect(t *testing.T) {
	g, err := DetermineRole("guest")
	if err != nil {
		t.Fail()
	}
	a, err := DetermineRole("admin")
	if err != nil {
		t.Fail()
	}
	m, err := DetermineRole("maintainer")
	if err != nil {
		t.Fail()
	}
	c, err := DetermineRole("creator")
	if err != nil {
		t.Fail()
	}

	if g != models.Guest || a != models.Admin || c != models.Creator || m != models.Maintainer {
		t.Fail()
	}
}

func TestRoleDeterminationShouldRetunUndefined(t *testing.T) {
	g, err := DetermineRole("other")
	if err == nil || g != models.Undefined {
		t.Fail()
	}
}
