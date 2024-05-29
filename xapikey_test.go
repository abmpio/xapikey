package xapikey

import "testing"

func Test_xapikey(t *testing.T) {
	ak, sk := GenerateAKSK()

	if len(ak) <= 0 && len(sk) <= 0 {
		t.Fatalf("Expect: ak and sk is not null,but actual is empty string")
	}
}
