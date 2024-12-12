package application

import (
	"encoding/csv"
	"fmt"
	"io"

	xlsx "github.com/thedatashed/xlsxreader"
)

func ParseDocument(document string) string {
	switch len(document) {
	case 11:
		return parseCPF(document)
	case 14:
		return parseCNPJ(document)
	default:
		return document
	}
}

func parseCPF(cpf string) string {
	return fmt.Sprintf("%s.%s.%s-%s", cpf[:3], cpf[3:6], cpf[6:9], cpf[9:11])
}

func parseCNPJ(cnpj string) string {
	return fmt.Sprintf("%s.%s.%s/%s-%s", cnpj[:2], cnpj[2:5], cnpj[5:8], cnpj[8:12], cnpj[12:14])
}

func readCSV(file io.Reader) ([][]string, error) {
	fileCSV := csv.NewReader(file)

	csvRows := make([][]string, 0)

	for {
		row, err := fileCSV.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		csvRows = append(csvRows, row)
	}

	return csvRows, nil
}

func readXLS(file io.Reader) ([][]string, error) {
	xlsFile, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	fileXLS, err := xlsx.NewReader(xlsFile)
	if err != nil {
		return nil, err
	}

	xlsRows := make([][]string, 0)

	for row := range fileXLS.ReadRows(fileXLS.Sheets[0]) {
		xlsxRowCells := make([]string, 0)
		for _, cell := range row.Cells {
			xlsxRowCells = append(xlsxRowCells, cell.Value)
		}

		xlsRows = append(xlsRows, xlsxRowCells)
	}

	return xlsRows, nil
}

func getColumnHeadersIndex(header []string) map[string]int {
	columnsIndex := make(map[string]int)
	for i, column := range header {
		columnsIndex[column] = i
	}
	return columnsIndex
}
