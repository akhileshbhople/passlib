package raw
import "golang.org/x/crypto/scrypt"
import "encoding/base64"
import "strings"
import "strconv"
import "fmt"

const ScryptSHA256RecommendedN = 16384
const ScryptSHA256Recommendedr = 8
const ScryptSHA256Recommendedp = 1

// Wrapper for golang.org/x/crypto/scrypt implementing a sensible
// modular crypt interface.
//
// password should be a UTF-8 plaintext password.
// salt should be a random salt value in binary form.
//
// N, r and p are parameters to scrypt.
//
// Returns a modular crypt hash.
func ScryptSHA256(password string, salt []byte, N, r, p int) string {
  passwordb := []byte(password)

  hash, err := scrypt.Key(passwordb, salt, N, r, p, 32)
  if err != nil {
    panic(err)
  }

  hstr := base64.StdEncoding.EncodeToString(hash)
  sstr := base64.StdEncoding.EncodeToString(salt)

  return fmt.Sprintf("$s2$%d$%d$%d$%s$%s", N, r, p, sstr, hstr)
}

func ParseScrypt(stub string) (salt, hash []byte, N,r,p int, err error) {
  if len(stub) < 10 || !strings.HasPrefix(stub, "$s2$") {
    err = ErrInvalidStub
    return
  }

  // $s2$  N$r$p$salt-base64$hash-base64
  parts := strings.Split(stub[4:], "$")

  if len(parts) < 4 {
    err = ErrInvalidStub
    return
  }

  var Ni, ri, pi uint64

  Ni, err = strconv.ParseUint(parts[0], 10, 31)
  if err != nil {
    return
  }

  ri, err = strconv.ParseUint(parts[1], 10, 31)
  if err != nil {
    return
  }

  pi, err = strconv.ParseUint(parts[2], 10, 31)
  if err != nil {
    return
  }

  N, r, p = int(Ni), int(ri), int(pi)

  salt, err = base64.StdEncoding.DecodeString(parts[3])
  if err != nil {
    return
  }

  if len(parts) >= 5 {
    hash, err = base64.StdEncoding.DecodeString(parts[4])
  }

  return
}
