package report

import (
	"bonds-report-service/internal/models/domain"
	"bonds-report-service/internal/models/domain/report"
	"bonds-report-service/internal/utils/logging"
	"context"
	"log/slog"

	"github.com/gladinov/e"
)

const (
	WithholdingOfPersonalIncomeTaxOnCoupons        = 2  // 2	Удержание НДФЛ по купонам.
	WithholdingOfPersonalIncomeTaxOnDividends      = 8  // 8    Удержание налога по дивидендам.
	PartialRedemptionOfBonds                       = 10 // 10	Частичное погашение облигаций.
	PurchaseOfSecurities                           = 15 // 15	Покупка ЦБ.
	PurchaseOfSecuritiesWithACard                  = 16 // 16	Покупка ЦБ с карты.
	TransferOfSecuritiesFromAnotherDepository      = 17 // 17	Перевод ценных бумаг из другого депозитария.
	WithhouldingACommissionForTheTransaction       = 19 // 19	Удержание комиссии за операцию.
	PaymentOfDividends                             = 21 // 21	Выплата дивидендов.
	SaleOfSecurities                               = 22 // 22	Продажа ЦБ.
	PaymentOfCoupons                               = 23 // 23 Выплата купонов.
	StampDuty                                      = 47 // 47	Гербовый сбор.
	TransferOfSecuritiesFromIISToABrokerageAccount = 57 // 57   Перевод ценных бумаг с ИИС на Брокерский счет
)

func (p *ReportProcessor) ProcessOperations(ctx context.Context, reportLine *domain.ReportLine) (_ *report.ReportPositions, err error) {
	const op = "report.ProcessOperations"

	defer logging.LogOperation_Debug(ctx, p.logger, op, &err)()

	processPosition := report.NewReportPositons()
	for _, operation := range reportLine.Operation {
		switch operation.Type {
		// 2	Удержание НДФЛ по купонам.
		// 8    Удержание налога по дивидендам.
		case WithholdingOfPersonalIncomeTaxOnCoupons, WithholdingOfPersonalIncomeTaxOnDividends:
			if err := processPosition.ProcessWithholdingOfPersonalIncomeTaxOnCouponsOrDividends(
				ctx,
				p.logger,
				operation); err != nil {
				return nil, e.WrapIfErr("failed to Process Withholding Of Personal Income Tax On Coupons Or Dividends", err)
			}

			// 10	Частичное погашение облигаций.
		case PartialRedemptionOfBonds:
			if err := processPosition.ProcessPartialRedemptionOfBonds(
				ctx,
				p.logger,
				operation); err != nil {
				return nil, err
			}

			// 15	Покупка ЦБ.
			// 16	Покупка ЦБ с карты.
			// 57   Перевод ценных бумаг с ИИС на Брокерский счет
		case PurchaseOfSecurities, PurchaseOfSecuritiesWithACard, TransferOfSecuritiesFromIISToABrokerageAccount:
			// Проверяем операцию на выполнение.
			// Т.е. операция может быть неисполнена, когда по заявленой цене не было встречного предложения
			if operation.QuantityDone == 0 {
				continue
			} else {
				processPosition.ProcessPurchaseOfSecurities(
					ctx,
					p.logger,
					operation,
					reportLine.Bond,
					reportLine.LastPrice,
					reportLine.Vunit_rate)
			}
			// 17	Перевод ценных бумаг из другого депозитария.
		case TransferOfSecuritiesFromAnotherDepository:
			// Проверяем операцию на выполнение.
			// Т.е. операция может быть неисполнена, когда по заявленой цене не было встречного предложения
			if operation.QuantityDone == 0 {
				continue
			} else {
				processPosition.ProcessTransferOfSecuritiesFromAnotherDepository(
					ctx,
					p.logger,
					operation,
					reportLine.Bond,
					reportLine.LastPrice,
					reportLine.Vunit_rate)
			}
			// 19	Удержание комиссии за операцию.
		case WithhouldingACommissionForTheTransaction:
			// Посчитали комисссию в операции покупки(15,16.17,57) и операции продажи(22)

			// 21	Выплата дивидендов.
		case PaymentOfDividends:
			err := processPosition.ProcessPaymentOfDividends(
				ctx,
				p.logger,
				operation)
			if err != nil {
				return nil, e.WrapIfErr("failed to Process Payment Of Dividends", err)
			}
			// 22	Продажа ЦБ.
		case SaleOfSecurities:
			// Проверяем операцию на выполнение.
			// Т.е. операция может быть неисполнена, когда по заявленой цене не было встречного предложения
			if operation.QuantityDone == 0 {
				continue
			} else {
				err := processPosition.ProcessSellOfSecurities(
					ctx,
					p.logger,
					&operation) // TODO: изменяем операцию . Надо здесь обдумать верна ли логика
				if err != nil {
					return nil, e.WrapIfErr("failed to Process Sell Of Securities", err)
				}
			}

			// 23 Выплата купонов.
		case PaymentOfCoupons:
			err := processPosition.ProcessPaymentOfCoupons(
				ctx,
				p.logger,
				operation,
			)
			if err != nil {
				return nil, e.WrapIfErr("failed to Process Payment Of Coupons", err)
			}

			// 47	Гербовый сбор.
		case StampDuty:
			err := processPosition.ProcessStampDuty(
				ctx,
				p.logger,
				operation,
			)
			if err != nil {
				return nil, e.WrapIfErr("failed to Process Stamp Duty", err)
			}
		default:
			p.logger.WarnContext(ctx, "unkown opperation type", slog.String("op", op))
			continue

		}
	}
	return processPosition, nil
}
