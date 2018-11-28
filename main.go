package main

import (
	"fmt"

	"github.com/satori/go.uuid"
)

func main() {
	// uu1 := uuid.Must(uuid.NewV1())
	// uu2 := uuid.Must(uuid.NewV2())
	// uu3 := uuid.Must(uuid.NewV3())
	uu4 := uuid.Must(uuid.NewV4())
	//u5 := uuid.Must(uuid.NewV5())
	fmt.Print( /*uu1, uu2, uu3, */ uu4 /*, uu5*/)
}
