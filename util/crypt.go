package util

import (
	"crypto/rand"
	"crypto/sha256"
)

/*func GenerateToken() string {
	h := sha256.New()
	fmt.Println(h.Sum(nil))

	b := make([]byte, 50)
	for i := range b {
		b[i] = i
	}

	return string(b)
}
*/

func HashPassword(password string) []byte {
	hash := sha256.New()
	hash.Write([]byte(password))

	return hash.Sum(nil)
}

func GenerateSession() []byte {
	session := make([]byte, 32)

	_, err := rand.Read(session)
	if err != nil {
		panic(err)
	}

	return session
}
