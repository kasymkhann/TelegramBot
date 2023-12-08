package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"sqlit-lessonTEST/lib/mistake"
	"sqlit-lessonTEST/storage"
)

const (
	defaultPerm = 0774
)

type Storage struct {
	basePath string
}

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return mistake.WrapErr("can't remove file", err)
	}
	path := filepath.Join(s.basePath, p.UserName, fileName)
	msg := fmt.Sprintf("can't remove file %s", err)

	if err := os.Remove(path); err != nil {
		return mistake.WrapErr(msg, err)
	}
	return nil
}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	filName, err := fileName(p)
	if err != nil {
		return false, mistake.WrapErr("can't exists file", err)
	}

	path := filepath.Join(s.basePath, p.UserName, filName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, err
	case err != nil:
		msg := fmt.Sprintf("can't exists file %s", path)

		return false, mistake.WrapErr(msg, err)

	}
	return true, err
}

func (s Storage) PickRandom(UserName string) (*storage.Page, error) {
	path := filepath.Join(s.basePath, UserName)

	file, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("error ReadDir path", err)
	}

	if len(file) == 0 {
		return nil, storage.E

	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(file))

	files := file[n]

	return s.decodePage(filepath.Join(path, files.Name()))

}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("can't decode page %w", err)
	}

	defer func() { _ = f.Close() }()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, fmt.Errorf("can't decode page %w", err)
	}

	return &p, nil
}

func (s Storage) Save(page *storage.Page) error {

	filePAth := filepath.Join(s.basePath, page.UserName)

	if err := os.Mkdir(filePAth, defaultPerm); err != nil {
		return fmt.Errorf("error Save in a files: %s", err)
	}

	fName, err := fileName(page)
	if err != nil {
		return err
	}

	filePAth = filepath.Join(filePAth, fName)

	file, err := os.Create(filePAth)
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}
	return nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
