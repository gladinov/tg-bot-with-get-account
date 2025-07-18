package telegram

import (
	"context"
	"log"
	"strings"

	"main.go/lib/e"
)

const (
	RndCmd        = "/rnd"
	HelpCmd       = "/help"
	StartCmd      = "/start"
	AccountsCmd   = "/accounts"
	GetBondReport = "/bondReport"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)
	var token string
	log.Printf("got new command '%s' from '%s' in chat: %v", text, username, chatID)

	haveToken, err := p.checkUser(chatID)
	if err != nil {
		return e.WrapIfErr("doCmd: can`t check availability of token", err)
	}

	if text == StartCmd {
		return p.sendHello(chatID)
	}

	switch haveToken {
	case true:
		token, err = p.storage.PickToken(context.Background(), chatID)
		if err != nil {
			return e.WrapIfErr("doCmd: can't pick token from storage", err)
		}
	case false:
		istoken, _ := p.isToken(text)
		switch istoken {
		case true:
			token = text
			p.storage.Save(context.Background(), username, chatID, text)
			return p.tg.SendMessage(chatID, msgTrueToken)
		case false:
			return p.tg.SendMessage(chatID, msgNoToken)
		}
	}

	switch text {
	case HelpCmd:
		return p.sendHelp(chatID)
	case AccountsCmd:
		return p.sendAccounts(chatID, token)
	case GetBondReport:
		return p.getBondReports(chatID, token)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) isToken(token string) (res bool, err error) {
	defer func() { err = e.WrapIfErr("isTokent error", err) }()
	if len(token) == 88 { // TODO:модифицировать проверку
		client := p.service.Tinkoffapi
		err := client.FillClient(token)
		if err != nil {
			return false, err
		}
		_, err = p.service.Tinkoffapi.GetAccToTgBot()
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, err
}

func (p *Processor) sendAccounts(chatID int, token string) error {
	client := p.service.Tinkoffapi
	err := client.FillClient(token)
	if err != nil {
		return e.WrapIfErr("sendAccounts: can't get tinkoffAPI client ", err)
	}

	accounts, err := p.service.Tinkoffapi.GetAccToTgBot()
	if err != nil {
		return e.WrapIfErr("sendAccounts: can't get accounts from tinkoffAPI client", err)
	}
	p.tg.SendMessage(chatID, accounts)
	return nil
}

func (p *Processor) getBondReports(chatID int, token string) (err error) {
	if err = p.service.GetBondReports(chatID, token); err != nil {
		return e.WrapIfErr("getBondReport: can't get Bond reports", err)
	}

	p.tg.SendMessage(chatID, "Отчеты по облигациям успешно сохранены в базу данных")
	return nil
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func (p *Processor) checkUser(chatId int) (res bool, err error) {
	defer func() { err = e.WrapIfErr("can't do command: checkUser", err) }()

	isExists, err := p.storage.IsExists(context.Background(), chatId)
	if err != nil {
		return false, err
	}

	return isExists, nil
}
