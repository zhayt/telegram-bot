package files

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/zhayt/read-adviser-bot/lib/e"
	"github.com/zhayt/read-adviser-bot/storage"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

// Storage implements interface
type Storage struct {
	basePath string // stored basePath where we locate all links
}

const defaultPerm = 0775

// New create Storage
func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(ctx context.Context, page *storage.Page) (err error) {
	defer func() { err = e.WrapIfErr("can't save page", err) }()

	// form the path where our file will be saved
	fPath := filepath.Join(s.basePath, page.UserName)

	// create all dir in fPath
	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}

	// form file name
	fName, err := fileName(page)
	if err != nil {
		return err
	}

	// join filename in path
	fPath = filepath.Join(fPath, fName)

	// create file
	file, err := os.Create(fPath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	// write page to the file in desired format
	if err = gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil
}

func (s Storage) PickRandom(ctx context.Context, userName string) (page *storage.Page, err error) {
	defer func() { err = e.WrapIfErr("can't pick random", err) }()

	// form the path where our files saved
	path := filepath.Join(s.basePath, userName)

	// get list of files
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	// get "random" number
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))

	// get random file
	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))
}

func (s Storage) Remove(ctx context.Context, p *storage.Page) error {
	// form file name
	fileName, err := fileName(p)
	if err != nil {
		return e.Wrap("can't remove file", err)
	}

	// form the path to file that will be deleted
	path := filepath.Join(s.basePath, p.UserName, fileName)

	// delete file
	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("can't remove file %s", path)

		return e.Wrap(msg, err)
	}

	return nil
}

func (s Storage) IsExists(ctx context.Context, p *storage.Page) (bool, error) {
	// form file name
	fileName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("can't check if file exists", err)
	}

	// form path to the file
	path := filepath.Join(s.basePath, p.UserName, fileName)

	// check file is exist or not
	switch _, err := os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file %s exists", path)

		return false, e.Wrap(msg, err)
	}

	return true, nil
}

// decodePage read the file content and return page or error
func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	// open file
	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("can't decode page", err)
	}
	defer func() { _ = f.Close() }()

	// create var where file will be decoding
	var p storage.Page

	// decoding file
	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, err
	}

	return &p, nil
}

// fileName generates a name for the file
func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
