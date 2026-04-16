//go:build !release

package mode

import (
	"fmt"
)

const CountPlayers int = 1

func init() {
	fmt.Printf("Debug mode is active. Use 'go build -tags release' to change\n\n")
}
