package snowflake

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node
var workerID int64

// init default snowflake with random workerID strategy,
func init() {
	var err error
	workerID, err = randomNumber(0, 1023)
	if err != nil {
		panic(fmt.Sprintf("Snowflake init default failed: %s", err))
	}
	node, err = snowflake.NewNode(workerID)
	if err != nil {
		panic(fmt.Sprintf("Snowflake create default node failed: %s", err))
	}
}

func GetID() string {
	return node.Generate().String()
}

func randomNumber(min, max int64) (int64, error) {
	// calculate the max we will be using
	bg := big.NewInt(max - min)

	// get big.Int between zero and bg
	// in this case 0 to 20
	n, err := rand.Int(rand.Reader, bg)
	if err != nil {
		return 0, err
	}

	// add n and min to support the passed in range
	return n.Int64() + min, nil
}
