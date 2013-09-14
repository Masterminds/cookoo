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

func AddToContext(cxt Context, params *Params) (interface{}, Interrupt) {
	p := params.AsMap()
	for k, v := range p {
		cxt.Add(k, v)
	}
	return true, nil
}
