package facades

import (
	"log"

	"github.com/desygner/desygner-toolchain/packages/mcp"
	"github.com/desygner/desygner-toolchain/packages/mcp/contracts"
)

func Mcp() contracts.Mcp {
	instance, err := mcp.App.Make(mcp.Binding)
	if err != nil {
		log.Println(err)
		return nil
	}

	return instance.(contracts.Mcp)
}
