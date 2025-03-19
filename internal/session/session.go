package session

import (
	"encoding/gob"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

var (
	store *sessions.CookieStore
)

const (
	sessionCookie = "microblog_session"
)

type Session struct {
	Timestamp int64
}

func init() {
	gob.Register(Session{})
}

func InitStore() {
	store = sessions.NewCookieStore([]byte(uuid.NewString()))
}

func Start(r *http.Request, w http.ResponseWriter) error {
	if store == nil {
		return errors.New("session store is nil")
	}

	session, err := store.New(r, sessionCookie)
	if err != nil {
		End(r, w)
	}

	session.Values["account"] = Session{Timestamp: time.Now().UTC().Unix()}
	return session.Save(r, w)
}

func End(r *http.Request, w http.ResponseWriter) error {
	if store == nil {
		return errors.New("session store is nil")
	}

	session, err := store.Get(r, sessionCookie)
	if err != nil {
		return err
	}
	session.Options.MaxAge = -1
	return session.Save(r, w)
}

func Get(r *http.Request, w http.ResponseWriter) (*Session, error) {
	if store == nil {
		return nil, errors.New("session store is nil")
	}

	session, err := store.Get(r, sessionCookie)
	if err != nil {
		return nil, err
	}

	if value, ok := session.Values["account"].(Session); ok {
		return &value, nil
	}
	return nil, errors.New("session is empty")

}
