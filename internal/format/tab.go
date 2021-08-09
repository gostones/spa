package format

import (
	"os"
	"text/tabwriter"
	"text/template"

	"github.com/gostones/spa/internal/log"
)

// PrintTab prints data in tabular form.
func PrintTab(data interface{}) {
	var tpl string
	t := template.New("tab")
	// switch data.(type) {
	// default:
	// 	log.Errorln("formatting template not found")
	// 	return
	// }
	t, _ = t.Parse(tpl)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	if err := t.Execute(w, data); err != nil {
		log.Errorln(err)
	}
	w.Flush()
}
