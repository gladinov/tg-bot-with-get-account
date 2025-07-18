package servicet_sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"main.go/lib/e"
	"main.go/service/service_models"
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

func (s *Storage) Init(ctx context.Context) error {
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

	_, err := s.db.ExecContext(ctx, q_bondReport)
	if err != nil {
		return fmt.Errorf("can't create bond reports table: %w", err)
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

	q_currencies := `CREATE TABLE IF NOT EXISTS currencies (
		on_date DATETIME,
		num_code TEXT,
		char_code TEXT,
		nominal INTEGER,
		name TEXT,
		value REAL,
		vunit_rate REAL
		)`
	_, err = s.db.ExecContext(ctx, q_currencies)
	if err != nil {
		return fmt.Errorf("can't create currencies table: %w", err)
	}

	return nil
}

func (s *Storage) SaveOperations(ctx context.Context, chatID int, accountId string, operations []service_models.Operation) error {
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

func (s *Storage) SaveBondReport(ctx context.Context, chatID int, accountId string, bondReport []service_models.BondReport) error {
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

func (s *Storage) GetOperations(ctx context.Context, chatId int, assetUid string, accountId string) ([]service_models.Operation, error) {
	q := "select name,date,type, figi, operation_id,quantity_done,instrument_type,instrument_uid,price,currency,accrued_int,commission, payment from operations where chatId == ? and broker_account_id == ? and asset_uid == ? order by date"

	operationRes := make([]service_models.Operation, 0)

	rows, err := s.db.QueryContext(ctx, q, chatId, accountId, assetUid)
	if err != nil {
		return nil, e.WrapIfErr("query error", err)
	}
	for rows.Next() {
		var operation service_models.Operation
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

// func (s *Storage) GetCurrencies(ctx context.Context) {

// }
