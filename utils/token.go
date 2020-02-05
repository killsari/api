package utils

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"io"
	"math/rand"
	"strconv"
	"time"
)

func SetAuthorization(c *gin.Context, token string) {
	//send token in header
	c.Header("Authorization", fmt.Sprintf("Bearer %s", token))
}

func CodeGen() string {
	return fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
}

func TokenGen(code string) string {
	time := time.Now().Unix()
	h := md5.New()
	io.WriteString(h, fmt.Sprintf("%s-%s", code, strconv.FormatInt(time, 10)))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Md5Gen(s string) string {
	hash := md5.New()
	hash.Write([]byte(s))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func PasswordHashGen(password string) (passwordHash string) {
	passwordHashByte, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(passwordHashByte)
}

func ComparePasswordHash(passwordHash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
}
