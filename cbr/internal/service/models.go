package service

const (
	layout = "02/01/2006"
)

type Currency struct {
	NumCode   string `xml:"NumCode" json:"numCode,omitempty"`
	CharCode  string `xml:"CharCode" json:"charCode,omitempty"`
	Nominal   string `xml:"Nominal" json:"nominal,omitempty"`
	Name      string `xml:"Name" json:"name,omitempty"`
	Value     string `xml:"Value" json:"value,omitempty"`
	VunitRate string `xml:"VunitRate" json:"vunitRate,omitempty"`
}

type CurrenciesResponce struct {
	Date   string     `xml:"Date,attr" json:"date,omitempty"`
	Valute []Currency `xml:"Valute" json:"valute,omitempty"`
}

var HappyPathCurrenciesInBytes = `{
    "date": "06.11.2025",
    "valute": [
        {
            "numCode": "036",
            "charCode": "AUD",
            "nominal": "1",
            "name": "Австралийский доллар",
            "value": "52,7076",
            "vunitRate": "52,7076"
        },
        {
            "numCode": "944",
            "charCode": "AZN",
            "nominal": "1",
            "name": "Азербайджанский манат",
            "value": "47,7579",
            "vunitRate": "47,7579"
        },
        {
            "numCode": "012",
            "charCode": "DZD",
            "nominal": "100",
            "name": "Алжирских динаров",
            "value": "62,1249",
            "vunitRate": "0,621249"
        },
        {
            "numCode": "826",
            "charCode": "GBP",
            "nominal": "1",
            "name": "Фунт стерлингов",
            "value": "105,9185",
            "vunitRate": "105,9185"
        },
        {
            "numCode": "051",
            "charCode": "AMD",
            "nominal": "100",
            "name": "Армянских драмов",
            "value": "21,2219",
            "vunitRate": "0,212219"
        },
        {
            "numCode": "048",
            "charCode": "BHD",
            "nominal": "1",
            "name": "Бахрейнский динар",
            "value": "215,8802",
            "vunitRate": "215,8802"
        },
        {
            "numCode": "933",
            "charCode": "BYN",
            "nominal": "1",
            "name": "Белорусский рубль",
            "value": "27,2673",
            "vunitRate": "27,2673"
        },
        {
            "numCode": "975",
            "charCode": "BGN",
            "nominal": "1",
            "name": "Болгарский лев",
            "value": "47,7004",
            "vunitRate": "47,7004"
        },
        {
            "numCode": "068",
            "charCode": "BOB",
            "nominal": "1",
            "name": "Боливиано",
            "value": "11,7494",
            "vunitRate": "11,7494"
        },
        {
            "numCode": "986",
            "charCode": "BRL",
            "nominal": "1",
            "name": "Бразильский реал",
            "value": "15,0787",
            "vunitRate": "15,0787"
        },
        {
            "numCode": "348",
            "charCode": "HUF",
            "nominal": "100",
            "name": "Форинтов",
            "value": "24,0551",
            "vunitRate": "0,240551"
        },
        {
            "numCode": "704",
            "charCode": "VND",
            "nominal": "10000",
            "name": "Донгов",
            "value": "32,3499",
            "vunitRate": "0,00323499"
        },
        {
            "numCode": "344",
            "charCode": "HKD",
            "nominal": "1",
            "name": "Гонконгский доллар",
            "value": "10,4597",
            "vunitRate": "10,4597"
        },
        {
            "numCode": "981",
            "charCode": "GEL",
            "nominal": "1",
            "name": "Лари",
            "value": "29,9489",
            "vunitRate": "29,9489"
        },
        {
            "numCode": "208",
            "charCode": "DKK",
            "nominal": "1",
            "name": "Датская крона",
            "value": "12,4961",
            "vunitRate": "12,4961"
        },
        {
            "numCode": "784",
            "charCode": "AED",
            "nominal": "1",
            "name": "Дирхам ОАЭ",
            "value": "22,1071",
            "vunitRate": "22,1071"
        },
        {
            "numCode": "840",
            "charCode": "USD",
            "nominal": "1",
            "name": "Доллар США",
            "value": "81,1885",
            "vunitRate": "81,1885"
        },
        {
            "numCode": "978",
            "charCode": "EUR",
            "nominal": "1",
            "name": "Евро",
            "value": "93,5131",
            "vunitRate": "93,5131"
        },
        {
            "numCode": "818",
            "charCode": "EGP",
            "nominal": "10",
            "name": "Египетских фунтов",
            "value": "17,1376",
            "vunitRate": "1,71376"
        },
        {
            "numCode": "356",
            "charCode": "INR",
            "nominal": "100",
            "name": "Индийских рупий",
            "value": "91,5964",
            "vunitRate": "0,915964"
        },
        {
            "numCode": "360",
            "charCode": "IDR",
            "nominal": "10000",
            "name": "Рупий",
            "value": "48,5461",
            "vunitRate": "0,00485461"
        },
        {
            "numCode": "364",
            "charCode": "IRR",
            "nominal": "100000",
            "name": "Иранских риалов",
            "value": "14,1128",
            "vunitRate": "0,000141128"
        },
        {
            "numCode": "398",
            "charCode": "KZT",
            "nominal": "100",
            "name": "Тенге",
            "value": "15,5388",
            "vunitRate": "0,155388"
        },
        {
            "numCode": "124",
            "charCode": "CAD",
            "nominal": "1",
            "name": "Канадский доллар",
            "value": "57,6214",
            "vunitRate": "57,6214"
        },
        {
            "numCode": "634",
            "charCode": "QAR",
            "nominal": "1",
            "name": "Катарский риал",
            "value": "22,3045",
            "vunitRate": "22,3045"
        },
        {
            "numCode": "417",
            "charCode": "KGS",
            "nominal": "100",
            "name": "Сомов",
            "value": "92,8399",
            "vunitRate": "0,928399"
        },
        {
            "numCode": "156",
            "charCode": "CNY",
            "nominal": "1",
            "name": "Юань",
            "value": "11,3362",
            "vunitRate": "11,3362"
        },
        {
            "numCode": "192",
            "charCode": "CUP",
            "nominal": "10",
            "name": "Кубинских песо",
            "value": "33,8285",
            "vunitRate": "3,38285"
        },
        {
            "numCode": "498",
            "charCode": "MDL",
            "nominal": "10",
            "name": "Молдавских леев",
            "value": "47,4021",
            "vunitRate": "4,74021"
        },
        {
            "numCode": "496",
            "charCode": "MNT",
            "nominal": "1000",
            "name": "Тугриков",
            "value": "22,6548",
            "vunitRate": "0,0226548"
        },
        {
            "numCode": "566",
            "charCode": "NGN",
            "nominal": "1000",
            "name": "Найр",
            "value": "56,6303",
            "vunitRate": "0,0566303"
        },
        {
            "numCode": "554",
            "charCode": "NZD",
            "nominal": "1",
            "name": "Новозеландский доллар",
            "value": "45,7944",
            "vunitRate": "45,7944"
        },
        {
            "numCode": "578",
            "charCode": "NOK",
            "nominal": "10",
            "name": "Норвежских крон",
            "value": "79,5583",
            "vunitRate": "7,95583"
        },
        {
            "numCode": "512",
            "charCode": "OMR",
            "nominal": "1",
            "name": "Оманский риал",
            "value": "211,1534",
            "vunitRate": "211,1534"
        },
        {
            "numCode": "985",
            "charCode": "PLN",
            "nominal": "1",
            "name": "Злотый",
            "value": "21,8890",
            "vunitRate": "21,889"
        },
        {
            "numCode": "682",
            "charCode": "SAR",
            "nominal": "1",
            "name": "Саудовский риял",
            "value": "21,6503",
            "vunitRate": "21,6503"
        },
        {
            "numCode": "946",
            "charCode": "RON",
            "nominal": "1",
            "name": "Румынский лей",
            "value": "18,3311",
            "vunitRate": "18,3311"
        },
        {
            "numCode": "960",
            "charCode": "XDR",
            "nominal": "1",
            "name": "СДР (специальные права заимствования)",
            "value": "110,0129",
            "vunitRate": "110,0129"
        },
        {
            "numCode": "702",
            "charCode": "SGD",
            "nominal": "1",
            "name": "Сингапурский доллар",
            "value": "62,1135",
            "vunitRate": "62,1135"
        },
        {
            "numCode": "972",
            "charCode": "TJS",
            "nominal": "10",
            "name": "Сомони",
            "value": "87,5461",
            "vunitRate": "8,75461"
        },
        {
            "numCode": "764",
            "charCode": "THB",
            "nominal": "10",
            "name": "Батов",
            "value": "24,9350",
            "vunitRate": "2,4935"
        },
        {
            "numCode": "050",
            "charCode": "BDT",
            "nominal": "100",
            "name": "Так",
            "value": "66,5471",
            "vunitRate": "0,665471"
        },
        {
            "numCode": "949",
            "charCode": "TRY",
            "nominal": "10",
            "name": "Турецких лир",
            "value": "19,3113",
            "vunitRate": "1,93113"
        },
        {
            "numCode": "934",
            "charCode": "TMT",
            "nominal": "1",
            "name": "Новый туркменский манат",
            "value": "23,1967",
            "vunitRate": "23,1967"
        },
        {
            "numCode": "860",
            "charCode": "UZS",
            "nominal": "10000",
            "name": "Узбекских сумов",
            "value": "67,9366",
            "vunitRate": "0,00679366"
        },
        {
            "numCode": "980",
            "charCode": "UAH",
            "nominal": "10",
            "name": "Гривен",
            "value": "19,2973",
            "vunitRate": "1,92973"
        },
        {
            "numCode": "203",
            "charCode": "CZK",
            "nominal": "10",
            "name": "Чешских крон",
            "value": "38,2766",
            "vunitRate": "3,82766"
        },
        {
            "numCode": "752",
            "charCode": "SEK",
            "nominal": "10",
            "name": "Шведских крон",
            "value": "84,9167",
            "vunitRate": "8,49167"
        },
        {
            "numCode": "756",
            "charCode": "CHF",
            "nominal": "1",
            "name": "Швейцарский франк",
            "value": "100,1462",
            "vunitRate": "100,1462"
        },
        {
            "numCode": "230",
            "charCode": "ETB",
            "nominal": "100",
            "name": "Эфиопских быров",
            "value": "53,0511",
            "vunitRate": "0,530511"
        },
        {
            "numCode": "941",
            "charCode": "RSD",
            "nominal": "100",
            "name": "Сербских динаров",
            "value": "79,5863",
            "vunitRate": "0,795863"
        },
        {
            "numCode": "710",
            "charCode": "ZAR",
            "nominal": "10",
            "name": "Рэндов",
            "value": "46,4832",
            "vunitRate": "4,64832"
        },
        {
            "numCode": "410",
            "charCode": "KRW",
            "nominal": "1000",
            "name": "Вон",
            "value": "56,4711",
            "vunitRate": "0,0564711"
        },
        {
            "numCode": "392",
            "charCode": "JPY",
            "nominal": "100",
            "name": "Иен",
            "value": "52,9122",
            "vunitRate": "0,529122"
        },
        {
            "numCode": "104",
            "charCode": "MMK",
            "nominal": "1000",
            "name": "Кьятов",
            "value": "38,6612",
            "vunitRate": "0,0386612"
        }
    ]
}`

var HappyPathCurrencies = CurrenciesResponce{Date: "06.11.2025",
	Valute: []Currency{
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

var xmlData = CurrenciesResponce{Date: "03.03.1995",
	Valute: []Currency{Currency{
		NumCode:   "036",
		CharCode:  "AUD",
		Nominal:   "1",
		Name:      "Австралийский доллар",
		Value:     "3334,8200",
		VunitRate: "3334,82",
	}}}

var xmlDataInBytesErr = []byte(``)
