package inmemory

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

var sessionId string
var userId string

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)

	sessionId = uuid.NewV4().String()
	userId = uuid.NewV4().String()

	r := m.Run()
	os.Exit(r)
}

func TestAddSession(t *testing.T) {
	assert := assert.New(t)

	AddSession(sessionId, userId)
	actual, ok := sessionDB[sessionId]

	assert.NotEqualValues(ok, false)
	assert.EqualValues(userId, actual)
}

func TestGetSession(t *testing.T) {
	assert := assert.New(t)

	AddSession(sessionId, userId)
	actual := GetSessionById(sessionId)

	assert.EqualValues(userId, actual)
}

func TestDeleteSessionNoSessionExists(t *testing.T) {
	assert := assert.New(t)

	DeleteAllSession()
	err := DeleteSession(sessionId)
	expected := "error: no session exists"

	assert.EqualValues(expected, err.Error())
}
func TestDeleteSessionSessionIdIsNotValid(t *testing.T) {
	assert := assert.New(t)

	AddSession(sessionId, userId)
	err := DeleteSession("")
	expected := "error: session id is not valid"

	assert.EqualValues(expected, err.Error())
}

func TestDeleteSessionIdNotExits(t *testing.T) {
	assert := assert.New(t)

	AddSession(sessionId, userId)
	err := DeleteSession("idnotexits")

	expected := "error: session id not exists"
	assert.EqualValues(expected, err.Error())
}

func TestFindSession(t *testing.T) {
	assert := assert.New(t)

	AddSession(sessionId, userId)
	err := DeleteSession("idnotexits")

	expected := "error: session id not exists"
	assert.EqualValues(expected, err.Error())
}

func TestSessionExists(t *testing.T) {
	assert := assert.New(t)

	AddSession(sessionId, userId)
	actual := SessionExists(sessionId)

	assert.EqualValues(true, actual)
}
