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

func generateHeader(headerList []string) string {
	var separatorList []string

	for _, column := range headerList {
		separatorList = append(separatorList, strings.Repeat("-", utf8.RuneCountInString(column)))
	}
	header := []string{formatRow(headerList), formatRow(separatorList)}

	return strings.Join(header, "\n")
}

func formatTable(headerList []string, table [][]string) string {
	var formattedRows []string

	for _, row := range table {
		formattedRows = append(formattedRows, formatRow(row))
	}

	joinedRows := strings.Join(formattedRows, "\n")

	return strings.Join([]string{generateHeader(headerList), joinedRows}, "\n") + "\n"
}

func RenderTable(headerList []string, table [][]string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprint(w, formatTable(headerList, table)+"\n")
}
