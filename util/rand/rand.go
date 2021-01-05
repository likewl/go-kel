package rand

import (
	"crypto/md5"
	"fmt"
	"io"
	"time"
)

func GenerateString(s string) string{
	h := md5.New()
	io.WriteString(h, time.Now().String())
	return fmt.Sprintf("%x",h.Sum([]byte(s)))
}
func SplitStringToName(s string) string {
	name :=s[:8]+"-"+s[8:13]+"-"+s[13:18]+"-"+s[18:25]
	return name
}