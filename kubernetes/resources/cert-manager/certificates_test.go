package certmanager

import (
	"fmt"
	"testing"
)

func TestNewCertificate(t *testing.T) {
	cert, err := NewCertificate("test", "default", []string{"test.prod.riser.sh"}, IssuerEnvironmentProd)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(cert)
}

func TestListCertificates(t *testing.T) {
	certificates, err := ListCertificates()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(certificates)
}
