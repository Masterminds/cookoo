// Basic commands for Cookoo.
package cookoo

import (
	"log"
)

// Print a message to the log.
//
// Params:
// - msg: The message to print
func LogMessage(cxt Context, params *Params) (interface{}, Interrupt) {
	msg := params.Get("msg", "tick")
	log.Print(msg)
	return nil, nil
}
