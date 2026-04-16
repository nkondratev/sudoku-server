//go:build release

package mode

import "fmt"

const CountPlayers int = 2

func init() {
	fmt.Printf("Release mode is active.\n\n")
}
