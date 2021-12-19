package session

import (
	"encoding/gob"
	"net/http"
)

func init() {
	gob.Register(Flash{})
}

type Flash struct {
	Message string
}

func (s *Store) GetFlashes(r *http.Request, w http.ResponseWriter) ([]Flash, error) {
	session, err := s.store.Get(r, sessionKey)
	if err != nil {
		return nil, err
	}

	sfs := session.Flashes()
	flashes := make([]Flash, len(sfs))
	for i, sf := range sfs {
		flash, ok := sf.(Flash)
		if !ok {
			return nil, ErrFlashInvalid
		}
		flashes[i] = flash
	}

	return flashes, session.Save(r, w)
}

func (s *Store) SetFlash(r *http.Request, w http.ResponseWriter, msg string) error {
	session, err := s.store.Get(r, sessionKey)
	if err != nil {
		return err
	}

	session.AddFlash(Flash{Message: msg})

	return session.Save(r, w)
}
