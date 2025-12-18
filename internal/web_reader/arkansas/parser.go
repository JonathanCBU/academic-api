package arkansas

import (
	"strconv"
)

type IParser interface {
	buildUrl(academicYear int, schoolId string) (string, error)
}

type Parser struct {
	IParser
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) buildUrl(academicYear int, schoolId string) (string, error) {
	fy := strconv.Itoa(academicYear - 2000 + 10)

	return base_url + "lea=" + schoolId + "&fy=" + fy + "&format=Excel", nil
}
