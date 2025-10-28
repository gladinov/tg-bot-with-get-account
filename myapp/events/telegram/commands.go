package telegram

import (
	"context"
	"log"
	"strconv"
	"strings"

	"main.go/lib/e"
)

const (
	HelpCmd                    = "/help"
	StartCmd                   = "/start"
	AccountsCmd                = "/accounts"
	GetBondReport              = "/bondfifo"
	GetGeneralBondReport       = "/bondreport"
	GetUSD                     = "/usd"
	GetPortfolioStructure      = "/portfoliostructure"
	GetUnionPortfolioStructure = "/unionportfoliostructure"
	GetUnionWithSber           = "/unionpswithsber"
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
		p.service.Tinkoffapi.Token = token
	case false:
		err := p.service.Tinkoffapi.IsToken(text)
		switch err {
		case nil:
			token = text
			p.storage.Save(context.Background(), username, chatID, text)
			return p.tg.SendMessage(chatID, msgTrueToken)
		default:
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
	case GetGeneralBondReport:
		return p.getBondRepotsWithPng(chatID, token)
	case GetUSD:
		return p.getUSD(chatID)
	case GetPortfolioStructure:
		return p.GetPortfolioStructure(chatID, token)
	case GetUnionPortfolioStructure:
		return p.GetUnionPortfolioStructure(chatID, token)
	case GetUnionWithSber:
		return p.GetUnionPortfolioStructureWithSber(chatID, token)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) getUSD(chatId int) error {
	usd, err := p.service.GetUsd()
	if err != nil {
		return e.WrapIfErr("can't get usd", err)
	}
	usdRes := strconv.FormatFloat(usd, 'f', 5, 64)
	p.tg.SendMessage(chatId, usdRes)
	return nil

}

func (p *Processor) sendAccounts(chatID int, token string) error {
	accounts, err := p.service.GetAccountsList(token)
	if err != nil {
		return e.WrapIfErr("can't get account", err)
	}
	p.tg.SendMessage(chatID, accounts)
	return nil
}

func (p *Processor) getBondReports(chatID int, token string) (err error) {
	if err = p.service.GetBondReportsByFifo(chatID, token); err != nil {
		return e.WrapIfErr("getBondReport: can't get Bond reports", err)
	}

	p.tg.SendMessage(chatID, "Отчет по облигациям по методу FIFO успешно сохранен в базу данных")
	return nil
}

func (p *Processor) getGeneralBondReport(chatID int, token string) (err error) {
	if err = p.service.GetBondReportsWithEachGeneralPosition(chatID, token); err != nil {
		return e.WrapIfErr("getBondReport: can't get general bond reports", err)
	}
	p.tg.SendMessage(chatID, "Общий отчет по облигациям успешно сохранен в базу данных")
	return nil
}

func (p *Processor) getBondRepotsWithPng(chatID int, token string) (err error) {
	reportsInByteByAccount, err := p.service.GetBondReports(chatID, token)
	if err != nil {
		return e.WrapIfErr("can't get bond report with png", err)
	}
	for _, reportsInByte := range reportsInByteByAccount {
		for _, v := range reportsInByte {

			switch len(v.Reports) {
			case 0:
				continue
			case 1:
				err = p.tg.SendImageFromBuffer(chatID, v.Reports[0].Data, v.Reports[0].Caption)
				if err != nil {
					return e.WrapIfErr("can't get bond report with png", err)
				}
			default:
				err = p.tg.SendMediaGroupFromBuffer(chatID, v.Reports)
				if err != nil {
					return e.WrapIfErr("can't get bond report with png", err)
				}
			}
		}
	}
	return nil
}

func (p *Processor) GetPortfolioStructure(chatID int, token string) (err error) {
	accounts, err := p.service.GetAccounts()
	if err != nil {
		return e.WrapIfErr("can't get portfolio structure", err)
	}
	for _, account := range accounts {
		if account.Status == 3 {
			continue
		}
		report, err := p.service.GetPortfolioStructure(account)
		if err != nil {
			return e.WrapIfErr("can't get portfolio structure", err)
		}
		p.tg.SendMessage(chatID, report)
	}
	return nil
}

func (p *Processor) GetUnionPortfolioStructure(chatID int, token string) (err error) {
	accounts, err := p.service.GetAccounts()
	if err != nil {
		return e.WrapIfErr("processor: can't get union portfolio structure", err)
	}
	unionPortfolioStructure, err := p.service.GetUnionPortfolioStructure(token, accounts)
	if err != nil {
		return e.WrapIfErr("processor: can't get union portfolio structure", err)
	}
	p.tg.SendMessage(chatID, unionPortfolioStructure)
	return nil
}

func (p *Processor) GetUnionPortfolioStructureWithSber(chatID int, token string) (err error) {
	accounts, err := p.service.GetAccounts()
	if err != nil {
		return e.WrapIfErr("processor: can't get union portfolio structure", err)
	}
	unionPortfolioStructure, err := p.service.GetUnionPortfolioStructureWithSber(accounts)
	if err != nil {
		return e.WrapIfErr("processor: can't get union portfolio structure", err)
	}
	p.tg.SendMessage(chatID, unionPortfolioStructure)
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
