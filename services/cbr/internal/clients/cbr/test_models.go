package cbr

import "cbr/internal/models"

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

var xmlData = models.CurrenciesResponce{
	Date: "03.03.1995",
	Currencies: []models.Currency{{
		NumCode:   "036",
		CharCode:  "AUD",
		Nominal:   "1",
		Name:      "Австралийский доллар",
		Value:     "3334,8200",
		VunitRate: "3334,82",
	}},
}

var xmlDataInBytesErr = []byte(``)
