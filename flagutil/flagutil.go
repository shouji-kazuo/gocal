package flagutil

import (
	"fmt"

	cli "gopkg.in/urfave/cli.v2"
)

func IsAllSpecified(ctx *cli.Context, flagNames ...string) error {
	for _, flagname := range flagNames {
		if !ctx.IsSet(flagname) {
			return fmt.Errorf("%s flag is not set.", flagname)
		}
	}
	return nil
}
