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

// func (s *TestSuite) TearDownSuite() {
// 	if s.secret != nil {
// 		secrets.DeleteSecret(s.secret.Name, s.secret.Namespace)
// 	}
// }

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

// func (s *TestSuite) Test2GetSecretByName() {
// 	secret, err := secrets.GetSecret(s.secret.Name, s.secret.Namespace)
// 	s.NoError(err)
// 	s.Equal(string(secret.Data["foo"]), "bar")
// }

// func (s *TestSuite) Test3GetSecretByLabel() {
// 	secret, err := secrets.GetSecretsByLabels(s.secret.Namespace, map[string]string{
// 		"foo": "bar",
// 	})
// 	s.NoError(err)
// 	s.Equal(len(secret), 1)
// }

// func (s *TestSuite) Test4DeleteSecret() {
// 	err := secrets.DeleteSecret(s.secret.Name, s.secret.Namespace)
// 	s.NoError(err)
// }
