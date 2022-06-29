package remote

import (
	"testing"

	"github.com/snyk/driftctl/enumeration"
	"github.com/snyk/driftctl/enumeration/alerter"
	common2 "github.com/snyk/driftctl/enumeration/remote/common"

	"github.com/snyk/driftctl/enumeration/resource"

	"github.com/stretchr/testify/assert"
)

func TestScannerShouldIgnoreType(t *testing.T) {

	// Initialize mocks
	alerter := alerter.NewAlerter()
	fakeEnumerator := &common2.MockEnumerator{}
	fakeEnumerator.On("SupportedType").Return(resource.ResourceType("FakeType"))
	fakeEnumerator.AssertNotCalled(t, "Enumerate")

	remoteLibrary := common2.NewRemoteLibrary()
	remoteLibrary.AddEnumerator(fakeEnumerator)

	testFilter := &enumeration.MockFilter{}
	testFilter.On("IsTypeIgnored", resource.ResourceType("FakeType")).Return(true)

	s := NewScanner(remoteLibrary, alerter, ScannerOptions{}, testFilter)
	_, err := s.Resources()
	assert.Nil(t, err)
	fakeEnumerator.AssertExpectations(t)
}
