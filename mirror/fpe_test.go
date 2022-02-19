package mirror

import (
	"crypto/sha256"
	"fmt"
	"github.com/capitalone/fpe/ff3"
	"github.com/pkg/profile"
	"math/rand"
	"testing"
)

func BenchmarkFPE(b *testing.B) {
	defer profile.Start(profile.ProfilePath("."), profile.MemProfile).Stop()

	pass_bob := "hello"
	ch := 52
	for i := 0; i < b.N; i++ {
		key := sha256.Sum256([]byte(pass_bob))

		nonce := make([]byte, 8)
		rand.Read(nonce)

		FF3, _ := ff3.NewCipher(ch, key[:32], nonce)
		msg := "sdfgh"
		ciphertext, _ := FF3.Encrypt(msg)
		_, _ = FF3.Decrypt(ciphertext)
	}
}

func TestFPE(t *testing.T) {

	pass_bob := "hello"

	key := sha256.Sum256([]byte(pass_bob))

	nonce := make([]byte, 8)
	rand.Read(nonce)

	msg := "sdfgh"

	ch := 52

	//argCount := len(os.Args[1:])

	//if (argCount>0) {msg= string(os.Args[1])}
	//if (argCount>1) {pass_bob= string(os.Args[2])}
	//if (argCount>2) {ch,_= strconv.Atoi(os.Args[3])}

	FF3, _ := ff3.NewCipher(ch, key[:32], nonce)

	ciphertext, _ := FF3.Encrypt(msg)

	plaintext, _ := FF3.Decrypt(ciphertext)

	fmt.Printf("Message:\t\t%s\n", msg)
	fmt.Printf("Number of chars:\t%d\n", ch)
	fmt.Printf("Passphrase:\t\t%s\nSalt:\t\t\t%x\n\n", pass_bob, nonce)
	fmt.Printf("Cipher:\t\t\t%s\n", ciphertext)
	fmt.Printf("Decipher:\t\t%s", plaintext)
}
