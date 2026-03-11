package facades

import (
	"log"

	"github.com/desygner/desygner-toolchain/packages/aws"
	"github.com/desygner/desygner-toolchain/packages/aws/contracts"
)

func Aws() contracts.Aws {
	instance, err := aws.App.Make(aws.Binding)
	if err != nil {
		log.Println(err)
		return nil
	}

	return instance.(contracts.Aws)
}
