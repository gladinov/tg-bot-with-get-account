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
		return p.getBondReportEachAccount(chatID, token)
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

func (p *Processor) getBondReportEachAccount(chatID int, token string) (err error) {
	defer func() { err = e.WrapIfErr("can't get bond reports", err) }()
	client := p.service.Tinkoffapi

	err = client.FillClient(token)
	if err != nil {
		return err
	}

	// Загружаем список валют
	//  Если они обновлены менее 12 часов назад, то достаем из БД
	// Если Более, то обновляем БД
	// Можно сохранять данные по другим датам, что бы не запрашивать постоянно одно и то же в цб
	// Сохранять дату курса валюты

	assetUidInstrumentUidMap, err := p.service.Tinkoffapi.GetAllAssetUids() // TODO: Переписать так чтобы запрос проходил один раз в день без вызова пользователя
	if err != nil {
		return err
	}
	accounts, err := p.service.Tinkoffapi.GetAcc()
	if err != nil {
		return err
	}

	for _, account := range accounts {
		err := p.service.Tinkoffapi.GetOpp(&account)
		if err != nil {
			return err
		}
		operations := p.service.TransOperations(account.Operations)

		err = p.storage.SaveOperations(context.Background(), chatID, account.Id, operations)
		if err != nil {
			return err
		}

		err = p.service.Tinkoffapi.GetPortf(&account)
		if err != nil {
			return err
		}

		portfolio, err := p.service.TransPositions(&account, assetUidInstrumentUidMap)
		if err != nil {
			return err
		}

		for _, v := range portfolio.BondPositions {
			operationsDb, err := p.storage.GetOperations(context.Background(), chatID, v.Identifiers.AssetUid, account.Id)
			if err != nil {
				return err
			}
			resultBondPosition, err := p.service.ProcessOperations(operationsDb)
			if err != nil {
				return err
			}
			bondReport, err := p.service.CreateBondReport(*resultBondPosition)
			if err != nil {
				return err
			}
			err = p.storage.SaveBondReport(context.Background(), chatID, account.Id, bondReport.BondsInRUB)
			if err != nil {
				return err
			}
		}
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
