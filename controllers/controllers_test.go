package controllers

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/d5avard/login/bd/inmemory"
	"github.com/d5avard/login/models"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	result := m.Run()
	os.Exit(result)
}

func TestAlreadySignedInCookieDoesNotExist(t *testing.T) {
	assert := assert.New(t)
	r := httptest.NewRequest("GET", "http://localhost:8080/index", nil)

	var expected = "session cookie does not exist"
	err := AlreadySignedIn(r)

	assert.EqualValues(expected, err.Error())
}

func TestAlreadySignedInSessionDoesNotExist(t *testing.T) {
	assert := assert.New(t)
	r := httptest.NewRequest("GET", "http://localhost:8080/index", nil)

	r.AddCookie(CreateCookie(session, "true"))

	var expected = "session in memory does not exist"
	err := AlreadySignedIn(r)

	assert.EqualValues(expected, err.Error())
}

func TestAlreadySignedInIUserCookieDoesNotExist(t *testing.T) {
	assert := assert.New(t)
	r := httptest.NewRequest("GET", "http://localhost:8080/index", nil)

	u := CreateUser()
	inmemory.AddUser(u)
	inmemory.AddSession(u.UUID, u.Email)
	r.AddCookie(CreateCookie(session, u.UUID))

	expected := "user cookie does not exist"
	err := AlreadySignedIn(r)

	assert.EqualValues(expected, err.Error())
}

func TestAlreadySignedInIUserInMemoryDoesNotExist(t *testing.T) {
	assert := assert.New(t)
	r := httptest.NewRequest("GET", "http://localhost:8080/index", nil)

	u := CreateUser()
	r.AddCookie(CreateCookie(u.UUID, "true"))
	inmemory.AddSession(u.UUID, u.Email)
	r.AddCookie(CreateCookie(session, u.UUID))

	expected := "user in memory does not exist"
	err := AlreadySignedIn(r)

	assert.EqualValues(expected, err.Error())
}

func TestAlreadySignedIn(t *testing.T) {
	assert := assert.New(t)
	r := httptest.NewRequest("GET", "http://localhost:8080/index", nil)

	u := CreateUser()
	inmemory.AddUser(u)
	r.AddCookie(CreateCookie(u.UUID, "true"))
	inmemory.AddSession(u.UUID, u.Email)
	r.AddCookie(CreateCookie(session, u.UUID))

	err := AlreadySignedIn(r)

	assert.EqualValues(nil, err)
}

func CreateUser() *models.User {
	u := &models.User{}
	u.UUID = uuid.NewV4().String()
	u.Email = "dsavard@example.com"
	var hash []byte
	var err error
	if hash, err = bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost); err != nil {
		log.Fatal("error:", err.Error())
		return nil
	}
	u.HashPassword = hash
	return u
}

func CreateCookie(n string, v string) *http.Cookie {
	c := &http.Cookie{
		Name:  n,
		Value: v,
	}
	return c
}
