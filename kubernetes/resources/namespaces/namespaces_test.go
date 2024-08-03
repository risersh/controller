package namespaces

import (
	"testing"

	"github.com/risersh/controller/conf"
	"github.com/stretchr/testify/suite"
	apiv1 "k8s.io/api/core/v1"
)

type TestSuite struct {
	suite.Suite
	namespace *apiv1.Namespace
}

func (s *TestSuite) SetupTest() {
	conf.Init()
}

func TestNamespaceSuiteRun(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TearDownSuite() {
	if s.namespace != nil {
		DeleteNamespace(s.namespace.GetName())
	}
}

func (s *TestSuite) Test1NewNamespace() {
	res, err := NewNamespace(NewNamespaceArgs{
		Name: "test",
		Labels: map[string]string{
			"foo": "bar",
		},
		Annotations: map[string]string{
			"foo": "bar",
		},
	})
	s.NoError(err)
	s.namespace = res
}
func (s *TestSuite) Test2GetNamespaceByName() {
	namespace, err := GetNamespaceByName(s.namespace.GetName())
	s.NoError(err)
	s.Equal(namespace.GetName(), s.namespace.GetName())
}

func (s *TestSuite) Test3DeleteNamespace() {
	err := DeleteNamespace(s.namespace.GetName())
	s.NoError(err)
}
