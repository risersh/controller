package certmanager

import (
	"testing"

	"github.com/risersh/controller/conf"
	"github.com/stretchr/testify/suite"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type TestSuite struct {
	suite.Suite
	issuer *unstructured.Unstructured
}

func (s *TestSuite) SetupTest() {
	conf.Init()
}

func TestIssuerSuiteRun(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TearDownSuite() {
	if s.issuer != nil {
		DeleteIssuer(s.issuer.GetName(), s.issuer.GetNamespace())
	}
}

func (s *TestSuite) Test1NewIssuerWithHTTPSolver() {
	res, err := NewIssuer(NewIssuerArgs{
		Name:      "test",
		Namespace: "default",
		Labels: map[string]string{
			"foo": "bar",
		},
		Solver: IssuerSolverTypeHTTP,
	})
	s.NoError(err)
	s.issuer = res
}
func (s *TestSuite) Test2GetIssuerByName() {
	issuer, err := GetIssuerByName(s.issuer.GetName(), s.issuer.GetNamespace())
	s.NoError(err)
	s.Equal(issuer.GetName(), s.issuer.GetName())
}

func (s *TestSuite) Test3DeleteIssuer() {
	err := DeleteIssuer(s.issuer.GetName(), s.issuer.GetNamespace())
	s.NoError(err)
}
