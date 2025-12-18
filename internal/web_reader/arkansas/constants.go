package arkansas

/*
2023-24
curl https://myschoolinfo.arkansas.gov/StandardReports/SRC\?lea\=6040704\&fy\=34\&format\=Excel --output test.xlsx

2022-23
https://myschoolinfo.arkansas.gov/StandardReports/SRC?lea=6040704&fy=33&format=Excel

lea: Local Education Agency (school id number)
fy: Fiscal Year (seems like xy -> 20xy-11 -- 20xy-10)
*/

const (
	base_url string = "https://myschoolinfo.arkansas.gov/StandardReports/SRC?"
)

var (
	demographic_groups = []string{"all"}
)
