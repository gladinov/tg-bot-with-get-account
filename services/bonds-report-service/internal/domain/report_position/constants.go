package report

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

const (
	EuroTransBuyCost       = 240.0 // Стоимость Евротранса при переводе из другого депозитария
	EuroTransInstrumentUID = "02b2ea14-3c4b-47e8-9548-45a8dbcc8f8a"
	baseTaxRate            = 0.13 // Налог с продажи ЦБ
)

const (
	layout     = "2006-01-02"
	hoursInDay = 24
	daysInYear = 365.25
)
