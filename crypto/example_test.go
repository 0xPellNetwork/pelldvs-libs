package crypto_test

import (
	"fmt"

	"github.com/0xPellNetwork/golibs/crypto"
)

func ExampleSha256() {
	sum := crypto.Sha256([]byte("This is PellDVS Interactor"))
	fmt.Printf("%x\n", sum)
	// Output:
	// 9bbb94f02b807c94fee1f898615ef3fedd54345ec326edbfbf311c141f818b0e
}
