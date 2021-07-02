// Author: Steve Zhang
// Date: 2020/10/16 3:02 下午

package password

import (
	"crypto/md5"
	"encoding/hex"

	"go-sever/util/uuid"
)

func MD5(data string) string {
	m := md5.Sum([]byte(data))
	return hex.EncodeToString(m[:])
}

func Hash(password string) (passwordHash, passwordSalt string) {
	passwordSalt = uuid.NewUUID().String()
	passwordHash = MD5(MD5(password) + passwordSalt)
	return
}

func Verify(password, passwordHash, passwordSalt string) bool {
	return passwordHash == MD5(MD5(password)+passwordSalt)
}
