// Package pretty provides a utility to print out organized data in a pretty
// manner.
//
// Table can be used as thus:
// prettyTable, _ := NewPrettyTable(
//   NewColumnDef("Name"),
//   NewColumnDef("Type"))
// prettyTable.AddRow("Noel", "Human")
// prettyTable.AddRow("David", "Cyborg")
// prettyTable.AddRow("Pranava", "Crusher")
// prettyTable.Print()
//
// Output looks like:
// +---------+---------+
// | Name    | Type    |
// +---------+---------+
// |    Noel |   Human |
// |   David |  Cyborg |
// | Pranava | Crusher |
// +---------+---------+
//
package pretty

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/fatih/color"
)

// Table creates formatted tables for human readability.
type Table struct {
	header              *string
	columnDefs          []ColumnDef
	rows                [][]string
	shouldPrintRowCount bool
}

// ColumnDef is a representation of a column definition with a name and a
// maximum width. The max width must be > 3, and the name must be shorter than
// the max width. Errors will happen on instantiation of the table.
type ColumnDef struct {
	name     string
	maxWidth *int
}

// NewColumnDef creates a ColumnDef with a name and no maximum width.
func NewColumnDef(name string) ColumnDef {
	return ColumnDef{name: name}
}

// NewColumnDefWithWidth creates a ColumnDef with a name and maximum width.
func NewColumnDefWithWidth(name string, maxWidth int) ColumnDef {
	return ColumnDef{
		name:     name,
		maxWidth: &maxWidth,
	}
}

type alignment uint

const (
	leftJustify  alignment = iota
	rightJustify alignment = iota
)

var (
	columnColors = []color.Attribute{
		color.FgRed,
		color.FgMagenta,
		color.FgBlue,
		color.FgWhite,
	}

	rowColors = []color.Attribute{
		color.FgYellow,
		color.FgGreen,
	}
)

// NewPrettyTable creates a new Table.
func NewPrettyTable(columnDefs ...ColumnDef) (*Table, error) {
	if len(columnDefs) < 1 {
		return nil, fmt.Errorf("must have at least 1 column")
	}

	for _, columnDef := range columnDefs {
		if columnDef.maxWidth == nil {
			continue
		}

		if *columnDef.maxWidth <= 3 {
			return nil, fmt.Errorf(
				"column %s max width %d must be greater than 3",
				columnDef.name,
				columnDef.maxWidth)
		}
		if strLengthWithEncoding(columnDef.name) > *columnDef.maxWidth {
			return nil, fmt.Errorf(
				"column name %s cannot be longer than max width %d",
				columnDef.name,
				columnDef.maxWidth)
		}
	}

	return &Table{
		columnDefs: columnDefs,
		rows:       make([][]string, 0),
	}, nil
}

// SetHeader creates a header for the table.
func (table *Table) SetHeader(header string) {
	table.header = &header
}

// ShowRowCount is a configuration, defaulted to false, that can be toggled
// on to print row count when Print() is called.
func (table *Table) ShowRowCount(showRowCount bool) {
	table.shouldPrintRowCount = showRowCount
}

// SetRows sets the rows of the table, overriding any that might
// currently be there.
func (table *Table) SetRows(rows [][]string) error {
	for _, row := range rows {
		if len(row) != len(table.columnDefs) {
			return fmt.Errorf(
				"row length %d must match columns %d",
				len(row),
				len(table.columnDefs))
		}
	}

	table.rows = rows
	return nil
}

// AddRow adds a row to the table.
func (table *Table) AddRow(row ...string) error {
	if err := table.validateRowSize(row); err != nil {
		return err
	}
	table.rows = append(table.rows, row)
	return nil
}

// PrettyString creates the pretty string representing this table.
func (table *Table) PrettyString() (string, error) {
	for _, row := range table.rows {
		err := table.validateRowSize(row)
		if err != nil {
			return "", err
		}
	}

	columnSizes := make([]int, len(table.columnDefs))
	for i, columnDef := range table.columnDefs {
		columnSize := strLengthWithEncoding(columnDef.name)
		for _, row := range table.rows {
			if strLengthWithEncoding(row[i]) > columnSize {
				columnSize = strLengthWithEncoding(row[i])
			}
		}

		if columnDef.maxWidth != nil && columnSize > *columnDef.maxWidth {
			columnSizes[i] = *columnDef.maxWidth
		} else {
			columnSizes[i] = columnSize
		}
	}

	var buffer bytes.Buffer

	var columnNames []string
	for _, columnDef := range table.columnDefs {
		columnNames = append(columnNames, columnDef.name)
	}

	// Write the header. Keep track of the length of the materialized header,
	// so that we can extend the header line in the case that the header is
	// longer than the width of the table.
	headerLength := 0
	if table.header != nil {
		var headerStr string
		headerStr, headerLength = renderHeader(*table.header)
		buffer.WriteString(headerStr)
	}

	// Write and create table borders
	headerLineStrings := make([]string, len(columnSizes))
	for i := range columnSizes {
		// Add 2 for the single space at beginning and end of cell
		headerLineStrings[i] = strings.Repeat("-", columnSizes[i]+2)
	}
	border := "+" + strings.Join(headerLineStrings, "+") + "+"

	// Extend upper border if the header is longer than the width of table.
	upperBorder := border
	if headerLength > len(upperBorder) {
		upperBorder = upperBorder +
			strings.Repeat("-", headerLength-len(upperBorder))
	}
	buffer.WriteString(upperBorder + "\n")
	border += "\n"

	// Write the column headers
	err := renderRow(&buffer, columnSizes, columnNames, columnColors, leftJustify)
	if err != nil {
		return "", err
	}
	buffer.WriteString("\n")

	// Write another border between columns and data rows.
	buffer.WriteString(border)

	// Write the content rows
	for _, row := range table.rows {
		err = renderRow(&buffer, columnSizes, row, rowColors, rightJustify)
		if err != nil {
			return "", err
		}
		buffer.WriteString("\n")
	}

	// Write the last border.
	buffer.WriteString(border)

	// Write row count, if needed.
	if table.shouldPrintRowCount {
		buffer.WriteString(
			fmt.Sprintf("Count: %d\n", len(table.rows)))
	}

	// Pretty print!
	return buffer.String(), nil
}

// Print prints the table to stdout.
func (table *Table) Print() error {
	strOutput, err := table.PrettyString()
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(os.Stdout, strOutput)
	return err
}

func (table *Table) validateRowSize(row []string) error {
	if len(row) != len(table.columnDefs) {
		return fmt.Errorf(
			"row length %d must match columns %d",
			len(row),
			len(table.columnDefs))
	}
	return nil
}

func renderRow(
	buffer *bytes.Buffer,
	columnSizes []int,
	contents []string,
	colors []color.Attribute,
	justification alignment,
) error {
	contentStrings := make([]string, len(contents))
	for i := range contents {
		cell, err := renderCell(
			contents[i],
			columnSizes[i],
			justification,
			colors[i%len(colors)])
		if err != nil {
			return err
		}
		contentStrings[i] = cell
	}
	_, err := buffer.WriteString(
		"|" + strings.Join(contentStrings, "|") + "|")
	return err
}

func renderCell(
	content string,
	cellLength int,
	justification alignment,
	textAttribute color.Attribute,
) (string, error) {
	truncatedContent := content
	if strLengthWithEncoding(content) > cellLength {
		truncatedContent = fmt.Sprintf(
			"%s...",
			truncateStringWithEncoding(content, cellLength-3))
	}

	paddingLength := cellLength - strLengthWithEncoding(truncatedContent)
	padding := strings.Repeat(" ", paddingLength)

	textColor := color.New(textAttribute, color.Bold)
	switch justification {
	case leftJustify:
		return textColor.Sprintf(" %s%s ", truncatedContent, padding), nil
	case rightJustify:
		return textColor.Sprintf(" %s%s ", padding, truncatedContent),
			nil
	default:
		return "", fmt.Errorf("did not match alignment")
	}
}

// renderHeader renders the header, as well as returns its horizontal length.
func renderHeader(header string) (string, int) {
	horizontalBorder := strings.Repeat("-", strLengthWithEncoding(header)+2)
	rendered := fmt.Sprintf(
		"%s\n %s |\n",
		horizontalBorder,
		header)

	return rendered, strLengthWithEncoding(horizontalBorder)
}

func strLengthWithEncoding(str string) int {
	length := 0
	for _, strRune := range str {
		if shouldCountEncodedRune(strRune) {
			length++
		}
	}
	return length
}

func truncateStringWithEncoding(str string, truncateLength int) string {
	if truncateLength == 0 {
		return ""
	}

	// Find the index at which we must truncate the string. Only truncate when
	// we absolutely must, i.e. when a counted rune puts us over the
	// truncateLength.
	strTruncateIndex := 0
	runeCount := 0
	for _, strRune := range str {
		if shouldCountEncodedRune(strRune) {
			if runeCount == truncateLength {
				break
			}
			runeCount++
		}
		strTruncateIndex++
	}

	return string([]rune(str)[:strTruncateIndex])
}

func shouldCountEncodedRune(r rune) bool {
	// DO NOT count non-spacing marks in the output!
	return !unicode.IsMark(r)
}
