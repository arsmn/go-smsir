package smsir

import (
	"sync"
	"time"
)

const expiryDelta = 10 * time.Second

type Token struct {
	AccessToken string
	Expiry      time.Time
}

func (t *Token) expired() bool {
	if t.Expiry.IsZero() {
		return false
	}
	return t.Expiry.Round(0).Add(-expiryDelta).Before(time.Now())
}

func (t *Token) Valid() bool {
	return t != nil && t.AccessToken != "" && !t.expired()
}

type TokenSource interface {
	Token() (*Token, error)
}

///// StaticTokenSource /////

func StaticTokenSource(t *Token) TokenSource {
	return staticTokenSource{t}
}

type staticTokenSource struct {
	t *Token
}

func (s staticTokenSource) Token() (*Token, error) {
	return s.t, nil
}

///// ReuseTokenSource /////

func ReuseTokenSource(t *Token, src TokenSource) TokenSource {
	if rt, ok := src.(*reuseTokenSource); ok {
		if t == nil {
			return rt
		}
		src = rt.new
	}
	return &reuseTokenSource{
		t:   t,
		new: src,
	}
}

type reuseTokenSource struct {
	new TokenSource

	mu sync.Mutex
	t  *Token
}

func (s *reuseTokenSource) Token() (*Token, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.t.Valid() {
		return s.t, nil
	}
	t, err := s.new.Token()
	if err != nil {
		return nil, err
	}
	s.t = t
	return t, nil
}
