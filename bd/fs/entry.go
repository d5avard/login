package fs

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/d5avard/diary/models"
	"google.golang.org/api/iterator"
)

const path string = "entry"

type entry struct {
	ctx    context.Context
	fscl   *firestore.Client
	fscoll *firestore.CollectionRef
}

var e *entry

func GetEntry() *entry {
	if e != nil {
		return e
	}
	log.Fatalln("entry is not init")
	return nil
}

func InitEntry(c context.Context, fc *firestore.Client) *entry {
	if e != nil {
		return e
	}
	e = new(entry)
	e.ctx = c
	e.fscl = fc
	e.fscoll = e.fscl.Collection(path)
	return e
}

func (e *entry) Add(data interface{}) error {
	_, _, err := e.fscoll.Add(e.ctx, data)
	if err != nil {
		log.Printf("failed adding new %s : %s", path, err.Error())
		return err
	}
	return nil
}

func (e *entry) GetAll() []*models.Entry {
	iter := e.fscoll.Documents(e.ctx)
	var ex []*models.Entry
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("failed to iterate: %v", err)
		}
		var e models.Entry
		m := doc.Data()
		e.Body = m["Body"].(string)
		e.Date = m["Date"].(string)
		e.Username = m["Username"].(string)
		ex = append(ex, &e)
	}
	return ex
}
