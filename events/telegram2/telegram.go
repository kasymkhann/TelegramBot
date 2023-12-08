package telegram2

import (
	"errors"

	"sqlit-lessonTEST/clients/telegram"
	"sqlit-lessonTEST/events"
	"sqlit-lessonTEST/lib/mistake"
	"sqlit-lessonTEST/storage"
)

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

type Processor struct {
	tgClient *telegram.Client
	offset   int
	Storage  storage.Storage
}

type Meta struct {
	ChatId   int
	Username string
}

func New(tgClient *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tgClient: tgClient,
		Storage:  storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tgClient.Update(p.offset, limit)
	if err != nil {
		return nil, mistake.WrapErr("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}
	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}
	p.offset = updates[len(updates)-1].Id + 1

	return res, nil
}

func (p Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return mistake.WrapErr("can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(e events.Event) error {
	meta, err := meta(e)
	if err != nil {
		return mistake.WrapErr("can't proccess message", err)
	}
	if err := p.doCmd(e.Text, meta.ChatId, meta.Username); err != nil {
		return mistake.WrapErr("can't proccess message", err)
	}
	return nil
}

func meta(e events.Event) (Meta, error) {
	res, ok := e.Meta.(Meta)
	if !ok {
		return Meta{}, mistake.WrapErr("can't get meta", ErrUnknownMetaType)
	}
	return res, nil
}

func event(u telegram.Update) events.Event {
	upType := fetchType(u)
	var res = events.Event{
		Type: upType,
		Text: fetchText(u),
	}
	if upType == events.Message {
		res.Meta = Meta{
			ChatId:   u.Message.Chat.Id,
			Username: u.Message.From.Username,
		}
	}
	return res
}

func fetchText(u telegram.Update) string {
	if u.Message == nil {
		return ""

	}
	return u.Message.Text
}

func fetchType(u telegram.Update) events.Type {
	if u.Message == nil {
		return events.Unknown
	}
	return events.Message
}
