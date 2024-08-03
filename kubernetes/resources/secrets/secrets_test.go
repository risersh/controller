package secrets_test

import (
	"testing"

	"github.com/risersh/controller/kubernetes/resources/secrets"
	"github.com/stretchr/testify/suite"
	corev1 "k8s.io/api/core/v1"
)

type TestSuite struct {
	suite.Suite
	secret *corev1.Secret
}

func (s *TestSuite) SetupTest() {

}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TearDownSuite() {
	if s.secret != nil {
		secrets.DeleteSecret(s.secret.Name, s.secret.Namespace)
	}
}

func (s *TestSuite) Test1NewSecret() {
	secret, err := secrets.NewSecret(secrets.NewSecretArgs{
		Name:      "test-secret",
		Namespace: "default",
		Labels: map[string]string{
			"foo": "bar",
		},
		Data: map[string][]byte{
			"foo": []byte("bar"),
		},
	})
	s.NoError(err)
	s.secret = secret
}

func (s *TestSuite) Test2GetSecretByName() {
	secret, err := secrets.GetSecret(s.secret.Name, s.secret.Namespace)
	s.NoError(err)
	s.Equal(string(secret.Data["foo"]), "bar")
}

func (s *TestSuite) Test3GetSecretByLabel() {
	secret, err := secrets.GetSecretsByLabels(s.secret.Namespace, map[string]string{
		"foo": "bar",
	})
	s.NoError(err)
	s.Equal(len(secret), 1)
}

func (s *TestSuite) Test4DeleteSecret() {
	err := secrets.DeleteSecret(s.secret.Name, s.secret.Namespace)
	s.NoError(err)
}
