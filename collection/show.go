package collection

import "fmt"
import "log"
import "github.com/gosuri/uitable"

import "github.com/tu-pm/contrail-introspect-cli/utils"

// Several representations of resources
type Shower interface {
	Long(maxColWidth uint)
	Short()
	Xml()
}

func (e Element) Xml() {
	fmt.Printf("%s", e.Node)
}
func (elts Elements) Xml() {
	for _, e := range elts {
		e.Xml()
	}
}
func (c Collection) Xml() {
	fmt.Printf("%s", c.rootNode)
}

func (e Element) Short() {
	if e.desc.ShortDetailXpath == "" {
		log.Fatal("ShortDetailXpath is not defined. Use XML output instead (-x option)")
		return
	}
	s, _ := e.Node.Search(e.desc.ShortDetailXpath)
	if len(s) != 1 {
		return
		// TODO Debug mode
		//log.Fatal("Xpath '" + e.desc.ShortDetailXpath + "' is not valid (check your ShortDetailXpath definition)")
	}
	fmt.Printf("%s\n", s[0])
}
func (col Collection) Short() {
	Elements(col.elements).Short()
}
func (elts Elements) Short() {
	for _, e := range elts {
		e.Short()
	}
}
func (e Element) Long(maxColWidth uint) {
	table := uitable.New()
	table.MaxColWidth = maxColWidth
	e.desc.LongDetail.LongFormat(table, FORMAT_TEXT, e)
	fmt.Println(table)
}
func (col Collection) Long(maxColWidth uint) {
	Elements(col.elements).Long(maxColWidth)
}
func (elts Elements) Long(maxColWidth uint) {
	table := uitable.New()
	table.MaxColWidth = maxColWidth
	for i, e := range elts {
		format := FORMAT_TABLE
		if i == 0 {
			format = FORMAT_TABLE_HEADER
		}
		e.desc.LongDetail.LongFormat(table, format, e)
	}
	fmt.Println(table)
}

// This is used to show the long version of an Element.
type LongFormatter interface {
	LongFormat(t *uitable.Table, f Format, e Element)
}

type LongFormatFn (func(*uitable.Table, Element))
type LongFormatXpaths []string

type Format uint8

const (
	FORMAT_TEXT         Format = 1
	FORMAT_TABLE_HEADER Format = 2
	FORMAT_TABLE        Format = 3
)

func (fn LongFormatFn) LongFormat(table *uitable.Table, format Format, e Element) {
	fn(table, e)
}

func (xpaths LongFormatXpaths) LongFormat(table *uitable.Table, format Format, e Element) {
	if format == FORMAT_TABLE_HEADER || format == FORMAT_TABLE {
		longFormatTable(table, format, e, xpaths)
	} else {
		for _, xpath := range xpaths {
			s, _ := e.Node.Search(xpath + "/text()")
			if len(s) == 1 {
				table.AddRow(utils.Pretty(s))
			}
		}
	}
}

func longFormatTable(table *uitable.Table, format Format, e Element, xpaths LongFormatXpaths) {
	if format == FORMAT_TABLE_HEADER {
		tmp := make([]interface{}, len(xpaths))
		for i, v := range xpaths {
			tmp[i] = v
		}
		table.AddRow(tmp...)
	}

	tmp := make([]interface{}, len(xpaths))
	for i, xpath := range xpaths {
		s, _ := e.Node.Search(xpath + "/text()")
		if len(s) == 1 {
			tmp[i] = utils.Pretty(s)
		}
	}
	table.AddRow(tmp...)
}
