package inmemory

import (
	"errors"
	"log"
)

var sessionDB map[string]string

func init() {
	sessionDB = map[string]string{}
}

func AddSession(sessionId string, userId string) {
	sessionDB[sessionId] = userId
	log.Println("create session:", sessionId)
}

func GetSessionById(sessionId string) string {
	return sessionDB[sessionId]
}

func SessionExists(sessionId string) bool {
	_, ok := sessionDB[sessionId]
	return ok
}

func DeleteSession(id string) error {
	if sessionDB == nil || len(sessionDB) < 1 {
		return errors.New("error: no session exists")
	}

	if id == "" {
		return errors.New("error: session id is not valid")
	}

	if _, ok := sessionDB[id]; !ok {
		return errors.New("error: session id not exists")
	}

	delete(sessionDB, id)
	log.Println("delete session in memory:", id)
	return nil
}

func DeleteAllSession() {
	for k := range sessionDB {
		delete(sessionDB, k)
	}
}

func FindSession(userId string) string {
	for k, v := range sessionDB {
		if v == userId {
			return k
		}
	}
	return ""
}
