package storage

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/zhayt/read-adviser-bot/lib/e"
	"io"
)

type Storage interface {
	Save(ctx context.Context, p *Page) error
	PickRandom(ctx context.Context, userName string) (*Page, error)
	Remove(ctx context.Context, p *Page) error
	IsExists(ctx context.Context, p *Page) (bool, error)
}

var ErrNoSavedPages = errors.New("no saves page")

// Page basic type data that Storage will work with
type Page struct {
	URL      string // link that gives client
	UserName string // client name
}

// Hash generate hash using Page.URL
func (p Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
