package certmanager

import (
	"fmt"
	"testing"
)

func TestListCertificates(t *testing.T) {
	certificates, err := ListCertificates()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(certificates)
}
