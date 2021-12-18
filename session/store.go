package session

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

const (
	sessionKey       = "thousand"
	currentUserIDKey = "currentUserID"
)

type Store struct {
	store sessions.Store
}

type StoreOptions struct {
	Key string
}

func NewStore(opts StoreOptions) *Store {
	return &Store{
		store: sessions.NewCookieStore([]byte(opts.Key)),
	}
}

func (s *Store) SetCurrentUserID(r *http.Request, w http.ResponseWriter, id uuid.UUID) error {
	session, _ := s.store.Get(r, sessionKey)

	session.Values[currentUserIDKey] = id.String()
	return session.Save(r, w)
}

func (s *Store) GetCurrentUserID(r *http.Request) (uuid.UUID, bool) {
	session, _ := s.store.Get(r, sessionKey)

	id, ok := session.Values[currentUserIDKey]
	if !ok {
		return uuid.UUID{}, false
	}

	idAsString, ok := id.(string)
	if !ok {
		return uuid.UUID{}, false
	}

	idAsUUID, err := uuid.Parse(idAsString)
	if err != nil {
		return uuid.UUID{}, false
	}

	return idAsUUID, true
}

func (s *Store) ClearCurrentUserID(r *http.Request, w http.ResponseWriter) error {
	session, _ := s.store.Get(r, sessionKey)

	delete(session.Values, currentUserIDKey)
	return session.Save(r, w)
}
