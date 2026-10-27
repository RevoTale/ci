// Command buildx runs the Woodpecker Buildx wrapper with required RevoTale metadata.
package main

import "github.com/RevoTale/ci/internal/buildx"

func main() {
	buildx.Run(true)
}
