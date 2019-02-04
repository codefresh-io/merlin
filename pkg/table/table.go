package table

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

type (
	Writer interface {
		Table() *tablewriter.Table
	}

	writer struct {
		table *tablewriter.Table
	}

	Options struct {
		Headers []string
	}
)

func New(opt *Options) Writer {
	t := tablewriter.NewWriter(os.Stdout)
	if len(opt.Headers) > 0 {
		t.SetHeader(opt.Headers)
	}
	t.SetAlignment(tablewriter.ALIGN_LEFT)
	t.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	t.SetColWidth(100)
	t.SetRowLine(true)
	t.SetAutoMergeCells(true)
	return &writer{
		table: t,
	}
}

func (w *writer) Table() *tablewriter.Table {
	return w.table
}
