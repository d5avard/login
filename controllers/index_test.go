package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/d5avard/login/bd/inmemory"
	"github.com/d5avard/login/utils"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestIndexOK(t *testing.T) {
	router := httprouter.New()
	Routes(router)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "http://localhost:8080/index", nil)
	assert := assert.New(t)

	u := CreateUser()
	inmemory.AddUser(u)
	request.AddCookie(CreateCookie(u.UUID, "true"))
	inmemory.AddSession(u.UUID, u.Email)
	request.AddCookie(CreateCookie(session, u.UUID))

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
