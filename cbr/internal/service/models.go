package service

const (
	layout = "02/01/2006"
)

type Valute struct {
	NumCode   string `xml:"NumCode" json:"numCode,omitempty"`
	CharCode  string `xml:"CharCode" json:"charCode,omitempty"`
	Nominal   string `xml:"Nominal" json:"nominal,omitempty"`
	Name      string `xml:"Name" json:"name,omitempty"`
	Value     string `xml:"Value" json:"value,omitempty"`
	VunitRate string `xml:"VunitRate" json:"vunitRate,omitempty"`
}

type ValCurs struct {
	Date   string   `xml:"Date,attr" json:"date,omitempty"`
	Valute []Valute `xml:"Valute" json:"valute,omitempty"`
}

var happyPathCurrencies = ValCurs{Date: "06.11.2025",
	Valute: []Valute{
		{NumCode: "036", CharCode: "AUD", Nominal: "1", Name: "Австралийский доллар", Value: "52,7076", VunitRate: "52,7076"},
		{NumCode: "944", CharCode: "AZN", Nominal: "1", Name: "Азербайджанский манат", Value: "47,7579", VunitRate: "47,7579"},
		{NumCode: "012", CharCode: "DZD", Nominal: "100", Name: "Алжирских динаров", Value: "62,1249", VunitRate: "0,621249"},
		{NumCode: "826", CharCode: "GBP", Nominal: "1", Name: "Фунт стерлингов", Value: "105,9185", VunitRate: "105,9185"},
		{NumCode: "051", CharCode: "AMD", Nominal: "100", Name: "Армянских драмов", Value: "21,2219", VunitRate: "0,212219"},
		{NumCode: "048", CharCode: "BHD", Nominal: "1", Name: "Бахрейнский динар", Value: "215,8802", VunitRate: "215,8802"},
		{NumCode: "933", CharCode: "BYN", Nominal: "1", Name: "Белорусский рубль", Value: "27,2673", VunitRate: "27,2673"},
		{NumCode: "975", CharCode: "BGN", Nominal: "1", Name: "Болгарский лев", Value: "47,7004", VunitRate: "47,7004"},
		{NumCode: "068", CharCode: "BOB", Nominal: "1", Name: "Боливиано", Value: "11,7494", VunitRate: "11,7494"},
		{NumCode: "986", CharCode: "BRL", Nominal: "1", Name: "Бразильский реал", Value: "15,0787", VunitRate: "15,0787"},
		{NumCode: "348", CharCode: "HUF", Nominal: "100", Name: "Форинтов", Value: "24,0551", VunitRate: "0,240551"},
		{NumCode: "704", CharCode: "VND", Nominal: "10000", Name: "Донгов", Value: "32,3499", VunitRate: "0,00323499"},
		{NumCode: "344", CharCode: "HKD", Nominal: "1", Name: "Гонконгский доллар", Value: "10,4597", VunitRate: "10,4597"},
		{NumCode: "981", CharCode: "GEL", Nominal: "1", Name: "Лари", Value: "29,9489", VunitRate: "29,9489"},
		{NumCode: "208", CharCode: "DKK", Nominal: "1", Name: "Датская крона", Value: "12,4961", VunitRate: "12,4961"},
		{NumCode: "784", CharCode: "AED", Nominal: "1", Name: "Дирхам ОАЭ", Value: "22,1071", VunitRate: "22,1071"},
		{NumCode: "840", CharCode: "USD", Nominal: "1", Name: "Доллар США", Value: "81,1885", VunitRate: "81,1885"},
		{NumCode: "978", CharCode: "EUR", Nominal: "1", Name: "Евро", Value: "93,5131", VunitRate: "93,5131"},
		{NumCode: "818", CharCode: "EGP", Nominal: "10", Name: "Египетских фунтов", Value: "17,1376", VunitRate: "1,71376"},
		{NumCode: "356", CharCode: "INR", Nominal: "100", Name: "Индийских рупий", Value: "91,5964", VunitRate: "0,915964"},
		{NumCode: "360", CharCode: "IDR", Nominal: "10000", Name: "Рупий", Value: "48,5461", VunitRate: "0,00485461"},
		{NumCode: "364", CharCode: "IRR", Nominal: "100000", Name: "Иранских риалов", Value: "14,1128", VunitRate: "0,000141128"},
		{NumCode: "398", CharCode: "KZT", Nominal: "100", Name: "Тенге", Value: "15,5388", VunitRate: "0,155388"},
		{NumCode: "124", CharCode: "CAD", Nominal: "1", Name: "Канадский доллар", Value: "57,6214", VunitRate: "57,6214"},
		{NumCode: "634", CharCode: "QAR", Nominal: "1", Name: "Катарский риал", Value: "22,3045", VunitRate: "22,3045"},
		{NumCode: "417", CharCode: "KGS", Nominal: "100", Name: "Сомов", Value: "92,8399", VunitRate: "0,928399"},
		{NumCode: "156", CharCode: "CNY", Nominal: "1", Name: "Юань", Value: "11,3362", VunitRate: "11,3362"},
		{NumCode: "192", CharCode: "CUP", Nominal: "10", Name: "Кубинских песо", Value: "33,8285", VunitRate: "3,38285"},
		{NumCode: "498", CharCode: "MDL", Nominal: "10", Name: "Молдавских леев", Value: "47,4021", VunitRate: "4,74021"},
		{NumCode: "496", CharCode: "MNT", Nominal: "1000", Name: "Тугриков", Value: "22,6548", VunitRate: "0,0226548"},
		{NumCode: "566", CharCode: "NGN", Nominal: "1000", Name: "Найр", Value: "56,6303", VunitRate: "0,0566303"},
		{NumCode: "554", CharCode: "NZD", Nominal: "1", Name: "Новозеландский доллар", Value: "45,7944", VunitRate: "45,7944"},
		{NumCode: "578", CharCode: "NOK", Nominal: "10", Name: "Норвежских крон", Value: "79,5583", VunitRate: "7,95583"},
		{NumCode: "512", CharCode: "OMR", Nominal: "1", Name: "Оманский риал", Value: "211,1534", VunitRate: "211,1534"},
		{NumCode: "985", CharCode: "PLN", Nominal: "1", Name: "Злотый", Value: "21,8890", VunitRate: "21,889"},
		{NumCode: "682", CharCode: "SAR", Nominal: "1", Name: "Саудовский риял", Value: "21,6503", VunitRate: "21,6503"},
		{NumCode: "946", CharCode: "RON", Nominal: "1", Name: "Румынский лей", Value: "18,3311", VunitRate: "18,3311"},
		{NumCode: "960", CharCode: "XDR", Nominal: "1", Name: "СДР (специальные права заимствования)", Value: "110,0129", VunitRate: "110,0129"},
		{NumCode: "702", CharCode: "SGD", Nominal: "1", Name: "Сингапурский доллар", Value: "62,1135", VunitRate: "62,1135"},
		{NumCode: "972", CharCode: "TJS", Nominal: "10", Name: "Сомони", Value: "87,5461", VunitRate: "8,75461"},
		{NumCode: "764", CharCode: "THB", Nominal: "10", Name: "Батов", Value: "24,9350", VunitRate: "2,4935"},
		{NumCode: "050", CharCode: "BDT", Nominal: "100", Name: "Так", Value: "66,5471", VunitRate: "0,665471"},
		{NumCode: "949", CharCode: "TRY", Nominal: "10", Name: "Турецких лир", Value: "19,3113", VunitRate: "1,93113"},
		{NumCode: "934", CharCode: "TMT", Nominal: "1", Name: "Новый туркменский манат", Value: "23,1967", VunitRate: "23,1967"},
		{NumCode: "860", CharCode: "UZS", Nominal: "10000", Name: "Узбекских сумов", Value: "67,9366", VunitRate: "0,00679366"},
		{NumCode: "980", CharCode: "UAH", Nominal: "10", Name: "Гривен", Value: "19,2973", VunitRate: "1,92973"},
		{NumCode: "203", CharCode: "CZK", Nominal: "10", Name: "Чешских крон", Value: "38,2766", VunitRate: "3,82766"},
		{NumCode: "752", CharCode: "SEK", Nominal: "10", Name: "Шведских крон", Value: "84,9167", VunitRate: "8,49167"},
		{NumCode: "756", CharCode: "CHF", Nominal: "1", Name: "Швейцарский франк", Value: "100,1462", VunitRate: "100,1462"},
		{NumCode: "230", CharCode: "ETB", Nominal: "100", Name: "Эфиопских быров", Value: "53,0511", VunitRate: "0,530511"},
		{NumCode: "941", CharCode: "RSD", Nominal: "100", Name: "Сербских динаров", Value: "79,5863", VunitRate: "0,795863"},
		{NumCode: "710", CharCode: "ZAR", Nominal: "10", Name: "Рэндов", Value: "46,4832", VunitRate: "4,64832"},
		{NumCode: "410", CharCode: "KRW", Nominal: "1000", Name: "Вон", Value: "56,4711", VunitRate: "0,0564711"},
		{NumCode: "392", CharCode: "JPY", Nominal: "100", Name: "Иен", Value: "52,9122", VunitRate: "0,529122"},
		{NumCode: "104", CharCode: "MMK", Nominal: "1000", Name: "Кьятов", Value: "38,6612", VunitRate: "0,0386612"},
	},
}

var xmlDataInBytes = []byte(`
<ValCurs Date="03.03.1995" name="Foreign Currency Market">
    <Valute ID="R01010">
        <NumCode>036</NumCode>
        <CharCode>AUD</CharCode>
        <Nominal>1</Nominal>
        <Name>Австралийский доллар</Name>
        <Value>3334,8200</Value>
        <VunitRate>3334,82</VunitRate>
    </Valute>
</ValCurs>
`)

var xmlData = ValCurs{Date: "03.03.1995",
	Valute: []Valute{Valute{
		NumCode:   "036",
		CharCode:  "AUD",
		Nominal:   "1",
		Name:      "Австралийский доллар",
		Value:     "3334,8200",
		VunitRate: "3334,82",
	}}}

var xmlDataInBytesErr = []byte(``)
