package servicet_sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"main.go/lib/e"
	"main.go/service/service_models"
)

const (
	hoursToUpdate = 12.0
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

	q_general_bond_report := `CREATE TABLE IF NOT EXISTS general_bond_report (
        id INTEGER PRIMARY KEY,
        chatId REAL,
        broker_account_id TEXT,
        name TEXT,
        ticker TEXT,
        currencies TEXT,
        quantity INTEGER,
        percent_of_portfolio TEXT,
        maturity_date DATETIME,
        duration INTEGER,
        buy_date DATETIME,
        position_price REAL,
        yield_to_maturity_on_purchase REAL,
        yield_to_maturity REAL,
        current_price REAL,
        nominal REAL,
        profit REAL,
        profit_in_percentage REAL
    )`

	_, err = s.db.ExecContext(ctx, q_general_bond_report)
	if err != nil {
		return fmt.Errorf("can't create general bond report positions table: %w", err)
	}

	q_uids := `CREATE TABLE IF NOT EXISTS uids (
		update_time DATETIME,
		instrument_uid TEXT,
		asset_uid TEXT
	)`
	_, err = s.db.ExecContext(ctx, q_uids)
	if err != nil {
		return fmt.Errorf("can't create uids table: %w", err)
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
		date DATETIME,
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

func (s *Storage) LastOperationTime(ctx context.Context, chatID int, accountId string) (time.Time, error) {
	q := "select date from operations where chatId == ? and broker_account_id == ? order by date desc LIMIT 1"

	var date time.Time

	err := s.db.QueryRowContext(ctx, q, chatID, accountId).Scan(&date)

	if err != nil {
		if err == sql.ErrNoRows {
			return date, service_models.ErrNoOpperations
		}
		return date, e.WrapIfErr("query error", err)
	}
	return date, nil

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

func (s *Storage) DeleteBondReport(ctx context.Context, chatID int, accountId string) (err error) {
	defer func() { err = e.WrapIfErr("can't delete bond report by chatId and accountId", err) }()
	q := "DELETE FROM bond_reports WHERE chatId = ? AND broker_account_id = ?"

	if _, err := s.db.ExecContext(ctx,
		q,
		chatID, accountId); err != nil {
		return err
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

func (s *Storage) DeleteGeneralBondReport(ctx context.Context, chatID int, accountId string) (err error) {
	defer func() { err = e.WrapIfErr("can't delete general bond report by chatId and accountId", err) }()
	q := "DELETE FROM general_bond_report WHERE chatId = ? AND broker_account_id = ?"

	if _, err := s.db.ExecContext(ctx,
		q,
		chatID, accountId); err != nil {
		return err
	}
	return nil
}
func (s *Storage) SaveGeneralBondReport(ctx context.Context, chatID int, accountId string, positions []service_models.GeneralBondReporPosition) error {
	q := `
    INSERT INTO general_bond_report (
        chatId,
        broker_account_id,
        name,
        ticker,
        currencies,
        quantity,
        percent_of_portfolio,
        maturity_date,
        duration,
        buy_date,
        position_price,
        yield_to_maturity_on_purchase,
        yield_to_maturity,
        current_price,
        nominal,
        profit,
        profit_in_percentage
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `

	for _, pos := range positions {
		if _, err := s.db.ExecContext(
			ctx,
			q,
			chatID,
			accountId,
			pos.Name,
			pos.Ticker,
			pos.Currencies,
			pos.Quantity,
			pos.PercentOfPortfolio,
			pos.MaturityDate,
			pos.Duration,
			pos.BuyDate,
			pos.PositionPrice,
			pos.YieldToMaturityOnPurchase,
			pos.YieldToMaturity,
			pos.CurrentPrice,
			pos.Nominal,
			pos.Profit,
			pos.ProfitInPercentage,
		); err != nil {
			return e.WrapIfErr("can't save general bond position", err)
		}
	}
	return nil
}

func (s *Storage) SaveUids(ctx context.Context, uids map[string]string) (err error) {
	defer func() { err = e.WrapIfErr("can't save uids", err) }()

	if len(uids) == 0 {
		return errors.New("empty uids map")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q_delete := "DELETE FROM uids"
	if _, err := tx.ExecContext(ctx, q_delete); err != nil {
		return err
	}

	q_insert := `
    INSERT INTO uids (
		update_time, 
		instrument_uid,
		asset_uid
		) VALUES (?, ?, ?)
	`
	now := time.Now()

	stmt, err := tx.PrepareContext(ctx, q_insert)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for instrument_uid, asset_uid := range uids {
		if _, err := stmt.ExecContext(ctx, now, instrument_uid, asset_uid); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Storage) IsUpdatedUids(ctx context.Context) (bool, error) {
	q := `SELECT update_time FROM uids LIMIT 1`

	var date time.Time

	err := s.db.QueryRowContext(ctx, q).Scan(&date)
	if err == sql.ErrNoRows {
		return false, service_models.ErrEmptyUids
	}
	if err != nil {
		return false, e.WrapIfErr("can't check update uids:", err)

	}

	if time.Since(date).Hours() > hoursToUpdate {
		return false, nil
	}
	return true, nil
}

func (s *Storage) GetUid(ctx context.Context, instrumentUid string) (string, error) {
	q := "SELECT asset_uid FROM uids WHERE instrument_uid = ?"

	var asset_uid string

	err := s.db.QueryRowContext(ctx, q, instrumentUid).Scan(&asset_uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", service_models.ErrEmptyUids
		} else {
			return "", e.WrapIfErr("can't get uid from DB", err)
		}
	}
	return asset_uid, nil
}

func (s *Storage) SaveCurrency(ctx context.Context, currencies service_models.Currencies, date time.Time) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO currencies 
        (date, num_code, char_code, nominal, name, value, vunit_rate)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, c := range currencies.CurrenciesMap {
		_, err = stmt.ExecContext(ctx,
			date,
			c.NumCode,
			c.CharCode,
			c.Nominal,
			c.Name,
			c.Value,
			c.VunitRate,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *Storage) GetCurrency(ctx context.Context, charCode string, date time.Time) (float64, error) {
	q := `
    SELECT vunit_rate
    FROM currencies
    WHERE char_code = ? 
    AND date(date) = date(?)
    LIMIT 1
    `

	var vunit_rate float64
	err := s.db.QueryRowContext(ctx, q, charCode, date.Format("2006-01-02")).Scan(&vunit_rate)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return vunit_rate, service_models.ErrNoCurrency
		}
		return vunit_rate, e.WrapIfErr("can't get currency from DB", err)
	}

	return vunit_rate, nil
}
