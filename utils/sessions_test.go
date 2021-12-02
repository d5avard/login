package utils

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

	r.AddCookie(createCookie(Session, "true"))

	var expected = "session in memory does not exist"
	err := AlreadySignedIn(r)

	assert.EqualValues(expected, err.Error())
}

// func TestAlreadySignedInIUserCookieDoesNotExist(t *testing.T) {
// 	assert := assert.New(t)
// 	r := httptest.NewRequest("GET", "http://localhost:8080/index", nil)

// 	u := createUser()
// 	inmemory.AddUser(u)

// 	id := uuid.NewV4().String()
// 	inmemory.AddSession(id, u.UUID)
// 	r.AddCookie(createCookie(Session, u.UUID))

// 	expected := "user cookie does not exist"
// 	err := AlreadySignedIn(r)

// 	assert.EqualValues(expected, err.Error())
// }

func TestAlreadySignedInIUserInMemoryDoesNotExist(t *testing.T) {
	assert := assert.New(t)
	r := httptest.NewRequest("GET", "http://localhost:8080/index", nil)

	u := createUser()
	// r.AddCookie(createCookie(u.UUID, "true"))
	sessionId := uuid.NewV4().String()
	inmemory.AddSession(sessionId, u.UUID)
	r.AddCookie(createCookie(Session, u.UUID))

	expected := "user in memory does not exist"
	err := AlreadySignedIn(r)

	assert.EqualValues(expected, err.Error())
}

func TestAlreadySignedIn(t *testing.T) {
	assert := assert.New(t)
	r := httptest.NewRequest("GET", "http://localhost:8080/index", nil)

	u := createUser()
	inmemory.AddUser(u)
	// r.AddCookie(createCookie(u.UUID, "true"))
	sessionId := uuid.NewV4().String()
	inmemory.AddSession(sessionId, u.UUID)
	r.AddCookie(createCookie(Session, u.UUID))

	err := AlreadySignedIn(r)

	assert.EqualValues(nil, err)
}

func createUser() *models.User {
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

func createCookie(n string, v string) *http.Cookie {
	c := &http.Cookie{
		Name:  n,
		Value: v,
	}
	return c
}
