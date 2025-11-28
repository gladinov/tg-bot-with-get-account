package telegram

import (
	"context"
	"errors"
	"fmt"
	"time"

	"log"

	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
	"main.go/internal/models"
	"main.go/lib/cryptoToken"
	"main.go/lib/e"
	"main.go/lib/valuefromcontext"
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

type TokenStatus int

const (
	TokenError TokenStatus = iota
	TokenFound
	TokenInserted
)

var constList = []string{HelpCmd,
	StartCmd,
	AccountsCmd,
	GetBondReport,
	GetGeneralBondReport,
	GetUSD,
	GetPortfolioStructure,
	GetUnionPortfolioStructure,
	GetUnionWithSber}

func ContainsInConstantCommands(text string) bool {
	for _, v := range constList {
		if text == v {
			return true
		}
	}
	return false
}

var ErrIncorrectToken = errors.New("incorrect token")

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)
	if ContainsInConstantCommands(text) {
		log.Printf("got new command '%s' from '%s' in chat: %v", text, username, chatID)
	} else {
		log.Printf("got new other command from '%s' in chat: %v", username, chatID)
	}

	chatIDStr := strconv.Itoa(chatID)
	ctx := context.Background()
	ctx = context.WithValue(ctx, models.ChatIdKey, chatIDStr)

	// TODO: Написать разные приветствия в зависимости от наличия токена
	if text == StartCmd {
		return p.sendHello(chatID)
	}

	tokenStatus, err := p.auth(ctx, text, username)

	switch {
	case err != nil && !errors.Is(err, ErrIncorrectToken):
		return err
	case errors.Is(err, ErrIncorrectToken):
		return p.tg.SendMessage(chatID, msgNoToken)
	case err == nil:
		switch tokenStatus {
		case TokenFound:

		case TokenInserted:
			p.tg.SendMessage(chatID, msgTrueToken)
		}
	}

	switch text {
	case HelpCmd:
		return p.sendHelp(chatID)
	case AccountsCmd:
		return p.sendAccounts(ctx, chatID)
	case GetBondReport:
		return p.getBondReports(ctx, chatID)
	case GetGeneralBondReport:
		return p.getBondRepotsWithPng(ctx, chatID)
	case GetUSD:
		return p.getUSD(ctx, chatID)
	case GetPortfolioStructure:
		return p.GetPortfolioStructure(ctx, chatID)
	case GetUnionPortfolioStructure:
		return p.GetUnionPortfolioStructure(ctx, chatID)
	case GetUnionWithSber:
		return p.GetUnionPortfolioStructureWithSber(ctx, chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) getUSD(ctx context.Context, chatId int) error {
	usdResponce, err := p.bondReportService.GetUsd(ctx)
	if err != nil {
		return e.WrapIfErr("can't get usd", err)
	}
	usd := strconv.FormatFloat(usdResponce.Usd, 'f', 5, 64)
	p.tg.SendMessage(chatId, usd)
	return nil

}

func (p *Processor) sendAccounts(ctx context.Context, chatID int) error {
	accountsResponce, err := p.bondReportService.GetAccountsList(ctx)
	if err != nil {
		return e.WrapIfErr("can't get account", err)
	}
	accounts := accountsResponce.Accounts
	p.tg.SendMessage(chatID, accounts)
	return nil
}

func (p *Processor) getBondReports(ctx context.Context, chatID int) (err error) {
	if err = p.bondReportService.GetBondReportsByFifo(ctx); err != nil {
		return e.WrapIfErr("getBondReport: can't get Bond reports", err)
	}

	p.tg.SendMessage(chatID, "Отчет по облигациям по методу FIFO успешно сохранен в базу данных")
	return nil
}

func (p *Processor) getBondRepotsWithPng(ctx context.Context, chatID int) (err error) {
	bondReportsResponce, err := p.bondReportService.GetBondReports(ctx)
	if err != nil {
		return e.WrapIfErr("can't get bond report with png", err)
	}
	reportsInByteByAccount := bondReportsResponce.Media
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

func (p *Processor) GetPortfolioStructure(ctx context.Context, chatID int) (err error) {
	portfolioStructures, err := p.bondReportService.GetPortfolioStructure(ctx)
	if err != nil {
		return e.WrapIfErr("can't get portfolio structure", err)
	}
	for _, report := range portfolioStructures.PortfolioStructures {
		p.tg.SendMessage(chatID, report)
	}
	return nil
}

func (p *Processor) GetUnionPortfolioStructure(ctx context.Context, chatID int) (err error) {

	unionPortfolioStructureResponce, err := p.bondReportService.GetUnionPortfolioStructure(ctx)
	if err != nil {
		return e.WrapIfErr("processor: can't get union portfolio structure", err)
	}
	unionPortfolioStructure := unionPortfolioStructureResponce.Report
	p.tg.SendMessage(chatID, unionPortfolioStructure)
	return nil
}

func (p *Processor) GetUnionPortfolioStructureWithSber(ctx context.Context, chatID int) (err error) {

	unionPortfolioStructureResponce, err := p.bondReportService.GetUnionPortfolioStructureWithSber(ctx)
	if err != nil {
		return e.WrapIfErr("processor: can't get union portfolio structure", err)
	}
	unionPortfolioStructure := unionPortfolioStructureResponce.Report
	p.tg.SendMessage(chatID, unionPortfolioStructure)
	return nil

}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func (p *Processor) checkUserToken(ctx context.Context) (res bool, err error) {
	const op = "processor:checkUserToken"

	isExists, err := p.storage.IsExistsToken(ctx)
	if err != nil {
		return false, err
	}

	return isExists, nil
}

func (p *Processor) tokenToBase64(token string) (string, error) {
	const op = "telegram.tokenToBase64"
	encryptedToken, err := cryptoToken.EncryptToken(token, p.tokenCrypter.Key)
	if err != nil {
		return "", e.WrapIfErr("could not encrypt token", err)
	}
	tokenInBase64, err := encryptedToken.ToBase64()
	if err != nil {
		return "", fmt.Errorf("%s:%w", op, err)
	}
	return tokenInBase64, nil
}

func (p *Processor) auth(ctx context.Context, text string, username string) (TokenStatus, error) {
	const op = "telegram.auth"
	chatIDStr, err := valuefromcontext.GetChatIDFromCtxStr(ctx)
	if err != nil {
		return TokenError, fmt.Errorf("%s:%w", op, err)
	}

	haveTokenInRedisErr := p.redis.Get(ctx, chatIDStr).Err()
	if haveTokenInRedisErr != nil && haveTokenInRedisErr != redis.Nil {
		return TokenError, fmt.Errorf("%s:%w", op, err)
	}
	if haveTokenInRedisErr == nil {
		return TokenFound, nil
	}

	haveToken, err := p.checkUserToken(ctx)
	if err != nil {
		return TokenError, fmt.Errorf("%s:%w", op, err)
	}

	switch haveToken {
	case true:
		tokenInBase64, err := p.storage.PickToken(ctx)
		if err != nil {
			return TokenError, fmt.Errorf("%s:%w", op, err)
		}

		err = p.redis.Set(ctx, chatIDStr, tokenInBase64, time.Until(time.Date(2030, time.December, 31, 0, 0, 0, 0, time.UTC))).Err()
		if err != nil {
			return TokenError, fmt.Errorf("%s:%w", op, err)
		}
		return TokenInserted, nil

	case false:
		err := p.isToken(ctx, text)
		if err != nil {
			return TokenError, ErrIncorrectToken
		}
		tokenInBase64, err := p.tokenToBase64(text)
		if err != nil {
			return TokenError, fmt.Errorf("%s:%w", op, err)
		}
		err = p.storage.Save(ctx, username, tokenInBase64)
		if err != nil {
			return TokenError, fmt.Errorf("%s:%w", op, err)
		}
		err = p.redis.Set(ctx, chatIDStr, tokenInBase64, time.Until(time.Date(2030, time.December, 31, 0, 0, 0, 0, time.UTC))).Err()
		if err != nil {
			return TokenError, fmt.Errorf("%s:%w", op, err)
		}
		return TokenInserted, nil
	}
	return TokenFound, nil
}

func (p *Processor) isToken(ctx context.Context, text string) error {
	const op = "telegram.isToken"
	chatID, err := valuefromcontext.GetChatIDFromCtxStr(ctx)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	if len(text) == 88 { // TODO:модифицировать проверку
		tokenInBase64, err := p.tokenToBase64(text)
		if err != nil {
			return fmt.Errorf("%s:%w", op, err)
		}
		err = p.redis.Set(ctx, chatID, tokenInBase64, time.Until(time.Date(2030, time.December, 31, 0, 0, 0, 0, time.UTC))).Err()
		if err != nil {
			return fmt.Errorf("%s:%w", op, err)
		}
		err = p.tinkoffApi.CheckToken(ctx)
		if err != nil {
			err := p.redis.Del(ctx, chatID).Err()
			if err != nil {
				return fmt.Errorf("%s:%w", op, err)
			}
			return fmt.Errorf("%s:%w", op, err)
		}

		return nil
	}

	return errors.New("is not token")
}
