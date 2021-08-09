package format

import (
	"encoding/json"
	"fmt"

	"github.com/gostones/spa/internal/log"
)

// PrintJSON formats data in json.
func PrintJSON(data interface{}) {
	d, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Errorln(err)
		return
	}
	fmt.Printf("%s\n", string(d))
}
