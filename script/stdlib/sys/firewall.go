package sys

import (
	"github.com/google/nftables"
	"github.com/kcarretto/paragon/script"
)

// Tables uses runtime.GOOS to detect what OS the agent is running on.
//
//
// @return (osStr, nil) iff success; (nil, err) o/w
func Tables(parser script.ArgParser) (script.Retval, error) {
	c := &nftables.Conn{}
	tables, _ := c.ListTables()
	var tablesList []string
	for _, table := range tables {
		tablesList = append(tablesList, table.Name)
	}
	return tablesList, nil
}
