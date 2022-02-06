package cli

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"unicode/utf8"
)

func formatRow(row []string) string {
	return strings.Join(row, "\t")
}

func generateSeparator(columns []string) []string {
	var separators []string

	for _, column := range columns {
		separators = append(separators, strings.Repeat("-", utf8.RuneCountInString(column)))
	}

	return separators
}

func generateHeader(headerList []string) string {
	var separatorList []string

	for _, column := range headerList {
		separatorList = append(separatorList, strings.Repeat("-", utf8.RuneCountInString(column)))
	}
	header := []string{formatRow(headerList), formatRow(separatorList)}

	return strings.Join(header, "\n")
}

func formatTable(table [][][]string) string {
	var formattedRows []string

	for _, sections := range table {
		for _, row := range sections {
			formattedRows = append(formattedRows, formatRow(row))
		}
		formattedRows = append(formattedRows, formatRow(generateSeparator(sections[len(sections)-1])))
	}

	return strings.Join(formattedRows, "\n")
}

func RenderTable(table ...[][]string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprint(w, formatTable(table))
}
