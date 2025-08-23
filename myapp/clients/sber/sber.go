package sber

import (
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
	"main.go/lib/e"
)

type ConfigSber struct {
	Bonds string `yaml:"Bonds"`
}

func LoadConfigSber(filename string) (_ ConfigSber, err error) {
	defer func() { err = e.WrapIfErr("load sber config error", err) }()
	var c ConfigSber
	input, err := os.ReadFile(filename)
	if err != nil {
		return ConfigSber{}, err
	}
	err = yaml.Unmarshal(input, &c)
	if err != nil {
		return ConfigSber{}, err
	}
	return c, nil
}

func ProcessConfigSber(config ConfigSber) (map[string]float64, error) {
	retBonds := make(map[string]float64)
	bonds := strings.Split(config.Bonds, ",")
	for _, v := range bonds {
		bond := strings.Split(v, ":")
		ticker := bond[0]
		quantity, err := strconv.Atoi(bond[1])
		if err != nil {
			return nil, e.WrapIfErr("can't process sber config", err)
		}
		retBonds[ticker] = float64(quantity)
	}
	return retBonds, nil
}
