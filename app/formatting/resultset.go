package formatting

import (
	"github.com/tvanriel/db-gui/app/domain"
)

type ResultSetJsonColumn struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type ResultSetJsonCell struct {
	Null  bool   `json:"null"`
	Value string `json:"value"`
}

type ResultSetJson struct {
	Columns []ResultSetJsonColumn `json:"columns"`
	Error   string                `json:"error"`
	Result  [][]ResultSetJsonCell `json:"result"`
}

func (res *ResultSetJson) AddColumn(Name string, Type string) {
	res.Columns = append(res.Columns, ResultSetJsonColumn{
		Name: Name,
		Type: Type,
	})
}

func (res *ResultSetJson) AddRow(Value []*string) {
	row := make([]ResultSetJsonCell, len(Value))

	for i := range Value {
		if Value[i] == nil {
			row[i] = ResultSetJsonCell{
				Null:  true,
				Value: "",
			}
		} else {
			row[i] = ResultSetJsonCell{Null: false, Value: *Value[i]}
		}
	}

	res.Result = append(res.Result, row)
}

func ToJSON(resultSet domain.Resultset) []ResultSetJson {

	result := resultSet.Get()
	response := make([]ResultSetJson, len(result))

	for i := range result {
		for j := range result[i].ColNames {
			response[i].AddColumn(result[i].ColNames[j], result[i].ColTypes[j])
		}
		for j := range result[i].Result {
			response[i].AddRow(result[i].Result[j])
		}
	}
	return response
}
