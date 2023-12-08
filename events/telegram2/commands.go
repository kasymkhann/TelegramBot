package telegram2

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"

	"sqlit-lessonTEST/storage"
)

const (
	RndCmd   = "/rnd"
	Help     = "/help"
	StartCmd = "/Start"
)

func (p *Processor) doCmd(text string, chatId int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new commad '%s' from '%s'", text, username)

	if isAddCmd(text) {
		return p.savePage(chatId, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatId, username)
	case Help:
		return p.sendHelp(chatId)
	case StartCmd:
		return p.sendHello(chatId)
	default:
		return p.tgClient.SendMessage(chatId, msgUnknowcommand)
	}
}

func (p *Processor) savePage(chatId int, UrlPage string, username string) (err error) {

	page := &storage.Page{
		URL:      UrlPage,
		UserName: username,
	}

	isExist, err := p.Storage.IsExists(page)
	if err != nil {
		return fmt.Errorf("error in a savePage %w", err)
	}
	if isExist {
		return p.tgClient.SendMessage(chatId, msgAlreadyExists)
	}

	if err := p.Storage.Save(page); err != nil {
		return fmt.Errorf("error no save in savePage: %w", err)
	}
	if err := p.tgClient.SendMessage(chatId, msgSaved); err != nil {
		return err
	}
	return nil

}

func (p *Processor) sendRandom(ChatId int, username string) (err error) {

	page, err := p.Storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.E) {
		return err
	}
	if errors.Is(err, storage.E) {
		return p.tgClient.SendMessage(ChatId, msgNoSaveMessage)
	}
	if err := p.tgClient.SendMessage(ChatId, page.URL); err != nil {
		return err
	}
	return p.Storage.Remove(page)
}

func (p *Processor) sendHelp(chatId int) error {
	return p.tgClient.SendMessage(chatId, msgHelp)
}

func (p *Processor) sendHello(chatId int) error {
	return p.tgClient.SendMessage(chatId, msgHello)
}
func isAddCmd(text string) bool {
	return isUrl(text)
}

func isUrl(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
