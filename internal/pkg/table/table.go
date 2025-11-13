package table

import (
	"fmt"
	"strings"
)

// Table represents a simple ASCII table
type Table struct {
	headers []string
	rows    [][]string
	widths  []int
}

// NewTable creates a new table with the given headers
func NewTable(headers ...string) *Table {
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	return &Table{
		headers: headers,
		rows:    [][]string{},
		widths:  widths,
	}
}

// AddRow adds a row to the table
func (t *Table) AddRow(values ...string) {
	if len(values) != len(t.headers) {
		panic(fmt.Sprintf("expected %d columns, got %d", len(t.headers), len(values)))
	}

	// Update column widths
	for i, v := range values {
		if len(v) > t.widths[i] {
			t.widths[i] = len(v)
		}
	}

	t.rows = append(t.rows, values)
}

// Render renders the table as a string
func (t *Table) Render() string {
	if len(t.rows) == 0 {
		return ""
	}

	var sb strings.Builder

	// Print header
	sb.WriteString("\n")
	for i, h := range t.headers {
		sb.WriteString(t.pad(h, t.widths[i]))
		if i < len(t.headers)-1 {
			sb.WriteString("  ")
		}
	}
	sb.WriteString("\n")

	// Print separator
	for i := range t.headers {
		sb.WriteString(strings.Repeat("-", t.widths[i]))
		if i < len(t.headers)-1 {
			sb.WriteString("  ")
		}
	}
	sb.WriteString("\n")

	// Print rows
	for _, row := range t.rows {
		for i, cell := range row {
			sb.WriteString(t.pad(cell, t.widths[i]))
			if i < len(row)-1 {
				sb.WriteString("  ")
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// Print prints the table to stdout
func (t *Table) Print() {
	fmt.Print(t.Render())
}

// pad pads a string to the given width
func (t *Table) pad(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}

// RowCount returns the number of rows in the table
func (t *Table) RowCount() int {
	return len(t.rows)
}
