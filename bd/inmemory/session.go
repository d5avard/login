package inmemory

import "log"

var SessionDB map[string]string

func init() {
	SessionDB = map[string]string{}
}

func AddSession(id string, email string) {
	SessionDB[id] = email
	log.Println("create session:", id)
}

func DeleteSession(id string) {
	if len(SessionDB) < 1 {
		return
	}

	if id == "" {
		return
	}

	if _, ok := SessionDB[id]; !ok {
		return
	}
	delete(SessionDB, id)
	log.Println("delete session:", id)
}

func FindSession(email string) string {
	for _, e := range SessionDB {
		if e == email {
			return e
		}
	}
	return ""
}
