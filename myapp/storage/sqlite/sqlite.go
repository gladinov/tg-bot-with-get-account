package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"main.go/lib/e"
	"main.go/service"
	"main.go/storage"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) Save(ctx context.Context, user_name string, chatID int, token string) error {
	q := `INSERT INTO users(user_name, chatID, token) VALUES (?,?,?)`

	if _, err := s.db.ExecContext(ctx, q, user_name, chatID, token); err != nil {
		return fmt.Errorf("can't save page: %w", err)
	}
	return nil
}

func (s *Storage) SavePositions(ctx context.Context, chatID int, accountId string, positions []service.PortfolioPosition) error {
	q := `INSERT INTO portfolios (
		chatId, 
		broker_account_id,
        figi,
        instrumentType,
        currency,
        quantity,
        averagePositionPrice,
        expectedYield,
        currentNkd,
        currentPrice,
        averagePositionPriceFifo,
        blocked,
        blockedLots,
        positionUid,
        instrumentUid,
        asset_uid,
        varMargin,
        expectedYieldFifo,
        dailyYield
    	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?,?)
		`
	for _, vals := range positions {

		if _, err := s.db.ExecContext(
			ctx,
			q,
			chatID,
			accountId,
			vals.Figi,
			vals.InstrumentType,
			vals.Currency,
			vals.Quantity,
			vals.AveragePositionPrice,
			vals.ExpectedYield,
			vals.CurrentNkd,
			vals.CurrentPrice,
			vals.AveragePositionPriceFifo,
			vals.Blocked,
			vals.BlockedLots,
			vals.PositionUid,
			vals.InstrumentUid,
			vals.AssetUid,
			vals.VarMargin,
			vals.ExpectedYieldFifo,
			vals.DailyYield); err != nil {
			return e.WrapIfErr("can't save portfolio positions", err)
		}
	}
	return nil
}

func (s *Storage) SaveOperations(ctx context.Context, chatID int, accountId string, operations []service.Operation) error {
	q := `
    INSERT INTO operations (
        chatId,
        broker_account_id,
		currency,
        operation_id,
        parent_operation_id,
        name,
        date,
        type,
        description,
        instrument_uid,
        figi,
        instrument_type,
        instrument_kind,
        position_uid,
        payment,
        price,
        commission,
        yield,
        yield_relative,
        accrued_int,
        quantity_done,
        asset_uid
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?)
    `

	for _, val := range operations {
		if _, err := s.db.ExecContext(
			ctx,
			q,
			chatID,
			val.BrokerAccountId,
			val.Currency,
			val.Operation_Id,
			val.ParentOperationId,
			val.Name,
			val.Date,
			val.Type,
			val.Description,
			val.InstrumentUid,
			val.Figi,
			val.InstrumentType,
			val.InstrumentKind,
			val.PositionUid,
			val.Payment,
			val.Price,
			val.Commission,
			val.Yield,
			val.YieldRelative,
			val.AccruedInt,
			val.QuantityDone,
			val.AssetUid); err != nil {
			return e.WrapIfErr("can't save operations", err)
		}
	}
	return nil
}

func (s *Storage) SaveBondReport(ctx context.Context, chatID int, accountId string, bondReport []service.BondReport) error {
	q := `
    INSERT INTO bond_reports (
		chatId,
		broker_account_id,
        name,
        ticker,
        maturity_date,
        offer_date,
        duration,
        buy_date,
        buy_price,
        yield_to_maturity_on_purchase,
        yield_to_offer_on_purchase,
        yield_to_maturity,
        yield_to_offer,
        current_price,
        nominal,
        profit,
        annualized_return
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?,?)
`
	for _, report := range bondReport {

		if _, err := s.db.ExecContext(
			ctx,
			q,
			chatID,
			accountId,
			report.Name,
			report.Ticker,
			report.MaturityDate,
			report.OfferDate,
			report.Duration,
			report.BuyDate,
			report.BuyPrice,
			report.YieldToMaturityOnPurchase,
			report.YieldToOfferOnPurchase,
			report.YieldToMaturity,
			report.YieldToOffer,
			report.CurrentPrice,
			report.Nominal,
			report.Profit,
			report.AnnualizedReturn); err != nil {
			return e.WrapIfErr("can't save bond report", err)
		}
	}
	return nil
}

func (s *Storage) PickToken(ctx context.Context, chatId int) (string, error) {
	q := `SELECT token FROM users WHERE chatID = ? LIMIT 1`

	var token string

	err := s.db.QueryRowContext(ctx, q, chatId).Scan(&token)
	if err == sql.ErrNoRows {
		return "", storage.ErrNoSavePages
	}
	if err != nil {
		return "", fmt.Errorf("can't pick random page: %w", err)

	}

	return token, nil
}

func (s *Storage) Remove(ctx context.Context, page *storage.Page) error {
	q := `DELETE FROM pages WHERE url = ? AND chatID = ?`
	if _, err := s.db.ExecContext(ctx, q, page.URL, page.ChatId); err != nil {
		return fmt.Errorf("can't remove page: %w", err)

	}
	return nil
}

func (s *Storage) IsExists(ctx context.Context, chatId int) (bool, error) {
	q := `SELECT COUNT(*) FROM users WHERE chatID = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, q, chatId).Scan(&count); err != nil {
		return false, fmt.Errorf("can't check user in storage: %w", err)
	}

	return count > 0, nil
}

func (s *Storage) Init(ctx context.Context) error {
	q_users := `CREATE TABLE IF NOT EXISTS users (user_name TEXT, chatID INTEGER, token TEXT)`

	_, err := s.db.ExecContext(ctx, q_users)
	if err != nil {
		return fmt.Errorf("can't create users table: %w", err)
	}

	q_portfolio := `CREATE TABLE IF NOT EXISTS portfolios (
            id integer primary key,
			chatId 				  REAL,
			broker_account_id         TEXT,
            figi                      TEXT,
            instrumentType            TEXT,
            currency                  TEXT,
            quantity                  REAL,
            averagePositionPrice      REAL,
            expectedYield             REAL,
            currentNkd                REAL,
            currentPrice              REAL,
            averagePositionPriceFifo  REAL,
            blocked                   BOOLEAN,
            blockedLots               REAL,
            positionUid               TEXT,
            instrumentUid             TEXT,
            asset_uid                 TEXT,
            varMargin                 REAL,
            expectedYieldFifo         REAL,
            dailyYield                REAL
        );`

	_, err = s.db.ExecContext(ctx, q_portfolio)
	if err != nil {
		return fmt.Errorf("can't create portfolio table: %w", err)
	}

	q_operations := `CREATE TABLE IF NOT EXISTS operations (
        id integer primary key,
		chatId 				REAL,
		broker_account_id       TEXT,
        currency                TEXT,
        operation_id            TEXT,
        parent_operation_id     TEXT,
        name                    TEXT,
        date                    DATETIME, 
        type                    INTEGER,
        description             TEXT,
        instrument_uid          TEXT,
        figi                    TEXT,
        instrument_type         TEXT,
        instrument_kind         TEXT,
        position_uid            TEXT,
        payment                 REAL,
        price                   REAL,
        commission              REAL,
        yield                   REAL,
        yield_relative          REAL,
        accrued_int             REAL,
        quantity_done           INTEGER,
        asset_uid               TEXT  
    )`

	_, err = s.db.ExecContext(ctx, q_operations)
	if err != nil {
		return fmt.Errorf("can't create operations table: %w", err)
	}

	q_bondReport := `CREATE TABLE IF NOT EXISTS bond_reports (
        id integer primary key,
		chatId 				REAL,
		broker_account_id       TEXT,
        name                            TEXT,
        ticker                          TEXT,
        maturity_date                   DATETIME,
        offer_date                      DATETIME,
        duration                        INTEGER,
        buy_date                        DATETIME,
        buy_price                       REAL,
        yield_to_maturity_on_purchase   REAL,
        yield_to_offer_on_purchase      REAL,
        yield_to_maturity               REAL,
        yield_to_offer                  REAL,
        current_price                   REAL,
        nominal                         REAL,
        profit                          REAL,
        annualized_return               REAL
    )`

	_, err = s.db.ExecContext(ctx, q_bondReport)
	if err != nil {
		return fmt.Errorf("can't create bond reports table: %w", err)
	}

	return nil
}

func (s *Storage) GetOperations(ctx context.Context, chatId int, assetUid string, accountId string) ([]service.Operation, error) {
	q := "select name,date,type, figi, operation_id,quantity_done,instrument_type,instrument_uid,price,currency,accrued_int,commission, payment from operations where chatId == ? and broker_account_id == ? and asset_uid == ? order by date"

	operationRes := make([]service.Operation, 0)

	rows, err := s.db.QueryContext(ctx, q, chatId, accountId, assetUid)
	if err != nil {
		return nil, e.WrapIfErr("query error", err)
	}
	for rows.Next() {
		var operation service.Operation
		err := rows.Scan(&operation.Name,
			&operation.Date,
			&operation.Type,
			&operation.Figi,
			&operation.Operation_Id,
			&operation.QuantityDone,
			&operation.InstrumentType,
			&operation.InstrumentUid,
			&operation.Price,
			&operation.Currency,
			&operation.AccruedInt,
			&operation.Commission,
			&operation.Payment)
		if err != nil {
			return nil, e.WrapIfErr("can't get operations", err)
		}
		operationRes = append(operationRes, operation)
	}
	if err := rows.Err(); err != nil {
		return nil, e.WrapIfErr("GetOperationsFromDB: rows.Err()", err)
	}

	return operationRes, nil

}
