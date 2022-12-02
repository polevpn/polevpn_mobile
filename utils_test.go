package polevpnmobile

import (
	"fmt"
	"testing"
)

func TestSubNetMask(t *testing.T) {
	fmt.Println(GetSubNetMask("192.168.0.1/24"))
}
