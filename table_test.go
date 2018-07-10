package pretty

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/rubrikinc/testwell/assert"
)

func TestBasicTable(t *testing.T) {
	table := createBasicTable(t)
	assertExpectedTable(t, table, "basic_table.txt")
}

func TestBasicTableWithHeader(t *testing.T) {
	table := createBasicTable(t)
	table.SetHeader("Employees")

	assertExpectedTable(t, table, "basic_table_with_header.txt")
}

func TestBasicTableWithRowCount(t *testing.T) {
	table := createBasicTable(t)
	table.ShowRowCount(true)

	assertExpectedTable(t, table, "basic_table_with_row_count.txt")
}

func TestTableWithSingleRow(t *testing.T) {
	table, err := NewPrettyTable(
		NewColumnDef("Only Column"))
	assert.Nil(t, err)

	err = table.AddRow("Some stuff")
	assert.Nil(t, err)
	err = table.AddRow("More stuff")
	assert.Nil(t, err)
	err = table.AddRow("hello")
	assert.Nil(t, err)
	err = table.AddRow("bye")
	assert.Nil(t, err)
	err = table.AddRow("hi")
	assert.Nil(t, err)

	assertExpectedTable(t, table, "table_with_single_row.txt")
}

func TestTableCreationWithNoRows(t *testing.T) {
	table, err := NewPrettyTable()
	assert.NotNil(t, err)
	assert.Nil(t, table)
}

func TestTableWithLongHeader(t *testing.T) {
	// Test the formatting of a table where the header is longer than the
	// width of the entire table.
	table := createBasicTable(t)
	table.SetHeader(
		"This is a really really really really really pretty long header")

	assertExpectedTable(t, table, "table_with_long_header.txt")
}

func TestTableWithSlightlyLongHeader(t *testing.T) {
	// Test the formatting of a table where the header is only 1 character
	// longer than the width of the table.
	table := createBasicTable(t)
	table.SetHeader(
		"This is a really really really really slightly big header")

	assertExpectedTable(t, table, "table_with_slightly_long_header.txt")
}

func TestTableWithSpecialCharacters(t *testing.T) {
	table, err := NewPrettyTable(
		NewColumnDef("Name"),
		NewColumnDef("Chars"))
	assert.Nil(t, err)

	// There are some special characters that do not visually result in added
	// horizontal space. If unaccounted for, these characters would mess up
	// the horizontal alignment of rows within columns.
	err = table.AddRow("no special chars", "just some written stuff")
	assert.Nil(t, err)
	err = table.AddRow(
		"couple special",
		"Funda̤̣o Municipal de Tecnologia da Informa̤̣o e Comunica̤̣o")
	assert.Nil(t, err)
	err = table.AddRow("just 1", "this one a̤̣")
	assert.Nil(t, err)

	assertExpectedTable(t, table, "table_with_special_chars.txt")
}

func TestTableWithColumnLimitAndSpecialCharactersAtEnd(t *testing.T) {
	table, err := NewPrettyTable(
		NewColumnDef("Name"),
		NewColumnDefWithWidth("Chars", 10))
	assert.Nil(t, err)

	// We must be careful when counting special characters that do not take up
	// horizontal space when truncating a string to make it fit in the cell.
	err = table.AddRow("no special chars", "no special chars")
	assert.Nil(t, err)
	// If we naively indexed into this string on truncation, result would be:
	// 01234a�...
	err = table.AddRow("special char right on border", "01234ạ56789")
	assert.Nil(t, err)
	// Naively, this would be: 0123o̤...
	err = table.AddRow("3-len char on border", "0123o̤̣456789")
	assert.Nil(t, err)
	err = table.AddRow(
		"all 3-lens",
		"o̤̣o̤̣o̤̣o̤̣o̤̣o̤̣o̤̣o̤̣o̤̣o̤̣o̤̣o̤̣o̤̣o̤̣o̤̣")
	assert.Nil(t, err)

	assertExpectedTable(
		t,
		table,
		"table_with_column_limit_and_special_chars.txt")
}

func TestBasicColumnLengthLimit(t *testing.T) {
	table, err := NewPrettyTable(
		NewColumnDef("Name"),
		NewColumnDefWithWidth("Words", 10))
	assert.Nil(t, err)

	err = table.AddRow("A", "Short")
	assert.Nil(t, err)
	err = table.AddRow("B", "short one")
	assert.Nil(t, err)
	err = table.AddRow("C", "exactly 10")
	assert.Nil(t, err)
	err = table.AddRow("D", "one too big")
	assert.Nil(t, err)
	err = table.AddRow("E", "this one is way too long")
	assert.Nil(t, err)

	assertExpectedTable(t, table, "table_with_column_limit.txt")
}

func TestColumnWidthLimitErrorWithLongColumnName(t *testing.T) {
	table, err := NewPrettyTable(
		NewColumnDef("Name"),
		NewColumnDefWithWidth("A Very Long Description", 10))

	assert.NotNil(t, err)
	assert.Nil(t, table)
}

func createBasicTable(t *testing.T) *Table {
	table, err := NewPrettyTable(
		NewColumnDef("Employee Number"),
		NewColumnDef("Name"),
		NewColumnDef("Type"),
		NewColumnDef("Phone Number"))
	assert.Nil(t, err)

	err = table.AddRow("23", "Noel", "Human", "(123) 456-7899")
	assert.Nil(t, err)
	err = table.AddRow("83", "David", "Cyborg", "987-654-3211")
	assert.Nil(t, err)
	err = table.AddRow("52", "Pranava", "Crusher", "1-800-123-4567")
	assert.Nil(t, err)
	err = table.AddRow("1182", "Postnava", "Kitten", "1 (800) 987-6543")
	assert.Nil(t, err)

	return table
}

func assertExpectedTable(
	t *testing.T,
	table *Table,
	filename string,
) {
	strOut, err := table.PrettyString()
	assert.Nil(t, err)

	filepath := path.Join("test", filename)
	expectedStr := readFileAsString(t, filepath)
	assert.EqualString(t, expectedStr, strOut)
}

func readFileAsString(t *testing.T, filename string) string {
	b, err := ioutil.ReadFile(filename)
	assert.Nil(t, err)
	return string(b)
}
