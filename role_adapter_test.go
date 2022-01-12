package fibercasbinrest_test

import (
	fibercasbinrest "github.com/prongbang/fiber-casbinrest"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRoleByTokenSuccess(t *testing.T) {
	// Given
	secret := "test"
	expect := []string{"ADMIN", "USER"}
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlcyI6WyJBRE1JTiIsIlVTRVIiXX0.7fOkJTmCRWJTVR5lLh0He_wtUTEUfWEFkrcfoArgPIw"
	adpt := fibercasbinrest.NewRoleAdapter(secret)

	// When
	actual, _ := adpt.GetRoleByToken(token)

	// Then
	assert.Equal(t, actual, expect)
}
