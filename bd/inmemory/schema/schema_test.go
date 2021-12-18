package inmemory

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)

	r := m.Run()
	os.Exit(r)
}
func TestAdd(t *testing.T) {
	assert := assert.New(t)

	newSchema := NewSchema()
	id := uuid.NewV4()
	value := uuid.NewV4()
	newSchema.Add(id, value)
	actual := newSchema.Get(id)

	assert.EqualValues(value, actual)
}
func TestDeleteNoValueExists(t *testing.T) {
	assert := assert.New(t)

	newSchema := NewSchema()
	id := uuid.NewV4()
	newSchema.DeleteAll()
	err := newSchema.Delete(id)
	expected := "error: no value exists"

	assert.EqualValues(expected, err.Error())
}
func TestDeleteIdIsNotValid(t *testing.T) {
	assert := assert.New(t)

	newSchema := NewSchema()
	id := uuid.NewV4()
	value := uuid.NewV4()

	newSchema.Add(id, value)
	err := newSchema.Delete(uuid.Nil)
	expected := "error: id not valid"

	assert.EqualValues(expected, err.Error())
}
func TestDeleteIdNotExits(t *testing.T) {
	assert := assert.New(t)

	newSchema := NewSchema()
	id := uuid.NewV4()
	value := uuid.NewV4()
	otherId := uuid.NewV4()

	newSchema.Add(id, value)
	err := newSchema.Delete(otherId)

	expected := "error: value id not exists"
	assert.EqualValues(expected, err.Error())
}
func TestSessionExists(t *testing.T) {
	assert := assert.New(t)

	newSchema := NewSchema()
	id := uuid.NewV4()
	value := uuid.NewV4()

	newSchema.Add(id, value)
	actual := newSchema.Exists(id)

	assert.EqualValues(true, actual)
}
