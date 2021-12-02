package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/d5avard/login/bd/inmemory"
	"github.com/d5avard/login/models"
	"github.com/d5avard/login/utils"
	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestIndexOK(t *testing.T) {
	router := httprouter.New()
	Routes(router)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "http://localhost:8080/index", nil)
	assert := assert.New(t)

	u := createUser()
	inmemory.AddUser(u)
	// request.AddCookie(createCookie(u.UUID, "true"))

	sessionId := uuid.NewV4().String()
	inmemory.AddSession(sessionId, u.UUID)
	request.AddCookie(createCookie(utils.Session, u.UUID))

	router.ServeHTTP(recorder, request)
	resp := recorder.Result()

	assert.EqualValues(resp.StatusCode, http.StatusOK)
	assert.EqualValues(utils.ContentTypeHTML, resp.Header.Get(utils.HeaderContentType))
}

func TestIndexUnauthorized(t *testing.T) {
	assert := assert.New(t)
	router := httprouter.New()
	Routes(router)
	recorder := httptest.NewRecorder()

	request := httptest.NewRequest("GET", "http://localhost:8080/index", nil)

	router.ServeHTTP(recorder, request)
	resp := recorder.Result()

	var body []byte
	var err error
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		t.Error(err.Error())
	}

	var restErr = utils.RestError{}
	if err = json.Unmarshal(body, &restErr); err != nil {
		t.Error(err.Error())
	}

	assert.EqualValues(resp.StatusCode, http.StatusUnauthorized)
	assert.EqualValues(restErr.Status, http.StatusUnauthorized)
	assert.EqualValues(utils.ContentTypeJSON, resp.Header.Get(utils.HeaderContentType))
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
