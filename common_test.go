package fibercasbinrest_test

import (
	fibercasbinrest "github.com/prongbang/fiber-casbinrest"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseRolesSuccess(t *testing.T) {
	// Given
	expect := []string{"ADMIN", "USER"}
	roles := []interface{}{"ADMIN", "USER"}

	// When
	actual := fibercasbinrest.ParseRoles(roles)

	// Then
	assert.Equal(t, actual, expect)
}

func TestGetValueSuccess(t *testing.T) {
	// Given
	expect := "John Doe"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiSm9obiBEb2UifQ.3fazvmF342WiHp5uhY-wkWArn-YJxq1IO7Msrtfk-OQ"
	key := "name"
	secret := []byte("test")

	// When
	actual, _ := fibercasbinrest.GetValue(token, key, secret)

	// Then
	assert.Equal(t, actual, expect)
}

func TestGetValueError(t *testing.T) {
	// Given
	expect := ""
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiSm9obiBEb2UifQ.3fazvmF342WiHp5uhY-wkWArn-YJxq1IO7Msrtfk-OQ"
	key := "name"
	secret := []byte("invalid-secret")

	// When
	actual, _ := fibercasbinrest.GetValue(token, key, secret)

	// Then
	assert.Equal(t, actual, expect)
}

func TestVerifyTrue(t *testing.T) {
	// Given
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiSm9obiBEb2UifQ.3fazvmF342WiHp5uhY-wkWArn-YJxq1IO7Msrtfk-OQ"
	secret := []byte("test")

	// When
	actual := fibercasbinrest.Verify(token, secret)

	// Then
	assert.Equal(t, actual, true)
}

func TestVerifyFalse(t *testing.T) {
	// Given
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiSm9obiBEb2UifQ.3fazvmF342WiHp5uhY-wkWArn-YJxq1IO7Msrtfk-OQ"
	secret := []byte("invalid-secret")

	// When
	actual := fibercasbinrest.Verify(token, secret)

	// Then
	assert.Equal(t, actual, false)
}
