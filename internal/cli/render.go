package cli

import (
	"bufio"
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

	return strings.Join(formattedRows, "\n") + "\n"
}

func RenderTable(table ...[][]string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprint(w, formatTable(table))
}

// Prompts the user for a confirmation [yes/no].
// Returns true/false if the user answers yes/no.
// Default choice defined by 'def'
func Confirmation(str string, def bool) (bool, error) {
	var defAns string
	reader := bufio.NewReader(os.Stdin)

	switch def {
	case true:
		defAns = "[Y/n]"
	case false:
		defAns = "[y/N]"
	}
	fmt.Printf("%s %s: ", str, defAns)
	answer, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	answer = strings.ToLower(strings.TrimSpace(answer))

	if answer == "y" || answer == "yes" {
		return true, nil
	} else if answer == "n" || answer == "no" {
		return false, nil
	}

	return def, nil
}
