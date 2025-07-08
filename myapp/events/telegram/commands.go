package telegram

import (
	"context"
	"errors"
	"log"
	"net/url"
	"strings"

	tinkoffapi "main.go/clients/tinkoffApi"
	"main.go/lib/e"
)

const (
	RndCmd      = "/rnd"
	HelpCmd     = "/help"
	StartCmd    = "/start"
	AccountsCmd = "/accounts"
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
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) isToken(token string) (bool, error) {
	if len(token) == 88 { // TODO:модифицировать проверку
		logger := p.tinkoffapi.Logg
		client := p.tinkoffapi
		err := client.FillClient(token)
		if err != nil {
			return false, e.WrapIfErr("isTokent: can't fillClient with token", err)
		}
		defer func() {
			logger.Infof("closing client connection")
			err := client.Client.Stop()
			if err != nil {
				logger.Errorf("client shutdown error %v", err.Error())
			}
		}()
		_, err = tinkoffapi.GetAcc(client.Client)
		if err != nil {
			return false, e.WrapIfErr("isTokent: can't get user account with token", err)
		}
		return true, nil
	}
	return false, errors.New("incorrect length of token")
}

func (p *Processor) sendAccounts(chatID int, token string) error {
	logger := p.tinkoffapi.Logg
	client := p.tinkoffapi
	err := client.FillClient(token)
	if err != nil {
		return e.WrapIfErr("sendAccounts: can't get tinkoffAPI client ", err)
	}
	defer func() {
		logger.Infof("closing client connection")
		err := client.Client.Stop()
		if err != nil {
			logger.Errorf("client shutdown error %v", err.Error())
		}
	}()
	accounts, err := tinkoffapi.GetAcc(client.Client)
	if err != nil {
		return e.WrapIfErr("sendAccounts: can't get accounts from tinkoffAPI client", err)
	}
	p.tg.SendMessage(chatID, accounts)
	return nil
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}
func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}

func (p *Processor) checkUser(chatId int) (res bool, err error) {
	defer func() { err = e.WrapIfErr("can't do command: checkUser", err) }()

	isExists, err := p.storage.IsExists(context.Background(), chatId)
	if err != nil {
		return false, err
	}

	return isExists, nil
}
