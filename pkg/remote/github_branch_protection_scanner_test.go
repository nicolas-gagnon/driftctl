package remote

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/snyk/driftctl/mocks"
	"github.com/snyk/driftctl/pkg/filter"
	"github.com/snyk/driftctl/pkg/remote/alerts"
	"github.com/snyk/driftctl/pkg/remote/cache"
	"github.com/snyk/driftctl/pkg/remote/common"
	remoteerr "github.com/snyk/driftctl/pkg/remote/error"
	"github.com/snyk/driftctl/pkg/remote/github"
	githubres "github.com/snyk/driftctl/pkg/resource/github"
	"github.com/snyk/driftctl/pkg/terraform"
	testresource "github.com/snyk/driftctl/test/resource"
	tftest "github.com/snyk/driftctl/test/terraform"
	"github.com/stretchr/testify/mock"

	"github.com/snyk/driftctl/pkg/resource"
	"github.com/snyk/driftctl/test"
	"github.com/snyk/driftctl/test/goldenfile"
	"github.com/stretchr/testify/assert"
)

func TestScanGithubBranchProtection(t *testing.T) {

	cases := []struct {
		test    string
		dirName string
		mocks   func(*github.MockGithubRepository, *mocks.AlerterInterface)
		err     error
	}{
		{
			test:    "no branch protection",
			dirName: "github_branch_protection_empty",
			mocks: func(client *github.MockGithubRepository, alerter *mocks.AlerterInterface) {
				client.On("ListBranchProtection").Return([]string{}, nil)
			},
			err: nil,
		},
		{
			test:    "Multiple branch protections",
			dirName: "github_branch_protection_multiples",
			mocks: func(client *github.MockGithubRepository, alerter *mocks.AlerterInterface) {
				client.On("ListBranchProtection").Return([]string{
					"MDIwOkJyYW5jaFByb3RlY3Rpb25SdWxlMTk1NDg0NzI=", //"repo0:main"
					"MDIwOkJyYW5jaFByb3RlY3Rpb25SdWxlMTk1NDg0Nzg=", //"repo0:toto"
					"MDIwOkJyYW5jaFByb3RlY3Rpb25SdWxlMTk1NDg0NzQ=", //"repo1:main"
					"MDIwOkJyYW5jaFByb3RlY3Rpb25SdWxlMTk1NDg0ODA=", //"repo1:toto"
					"MDIwOkJyYW5jaFByb3RlY3Rpb25SdWxlMTk1NDg0NzE=", //"repo2:main"
					"MDIwOkJyYW5jaFByb3RlY3Rpb25SdWxlMTk1NDg0Nzc=", //"repo2:toto"
				}, nil)
			},
			err: nil,
		},
		{
			test:    "cannot list branch protections",
			dirName: "github_branch_protection_empty",
			mocks: func(client *github.MockGithubRepository, alerter *mocks.AlerterInterface) {
				client.On("ListBranchProtection").Return(nil, errors.New("Your token has not been granted the required scopes to execute this query."))

				alerter.On("SendAlert", githubres.GithubBranchProtectionResourceType, alerts.NewRemoteAccessDeniedAlert(common.RemoteGithubTerraform, remoteerr.NewResourceListingErrorWithType(errors.New("Your token has not been granted the required scopes to execute this query."), githubres.GithubBranchProtectionResourceType, githubres.GithubBranchProtectionResourceType), alerts.EnumerationPhase)).Return()
			},
			err: nil,
		},
	}

	schemaRepository := testresource.InitFakeSchemaRepository("github", "4.4.0")
	githubres.InitResourcesMetadata(schemaRepository)
	factory := terraform.NewTerraformResourceFactory(schemaRepository)
	deserializer := resource.NewDeserializer(factory)

	for _, c := range cases {
		t.Run(c.test, func(tt *testing.T) {
			shouldUpdate := c.dirName == *goldenfile.Update

			scanOptions := ScannerOptions{Deep: true}

			providerLibrary := terraform.NewProviderLibrary()
			remoteLibrary := common.NewRemoteLibrary()

			// Initialize mocks
			alerter := &mocks.AlerterInterface{}
			mockedRepo := github.MockGithubRepository{}
			c.mocks(&mockedRepo, alerter)

			var repo github.GithubRepository = &mockedRepo

			realProvider, err := tftest.InitTestGithubProvider(providerLibrary, "4.4.0")
			if err != nil {
				t.Fatal(err)
			}
			provider := tftest.NewFakeTerraformProvider(realProvider)
			provider.WithResponse(c.dirName)

			if shouldUpdate {
				err := realProvider.Init()
				if err != nil {
					t.Fatal(err)
				}
				provider.ShouldUpdate()
				repo = github.NewGithubRepository(realProvider.GetConfig(), cache.New(0))
			}

			remoteLibrary.AddEnumerator(github.NewGithubBranchProtectionEnumerator(repo, factory))
			remoteLibrary.AddDetailsFetcher(githubres.GithubBranchProtectionResourceType, common.NewGenericDetailsFetcher(githubres.GithubBranchProtectionResourceType, provider, deserializer))

			testFilter := &filter.MockFilter{}
			testFilter.On("IsTypeIgnored", mock.Anything).Return(false)

			s := NewScanner(remoteLibrary, alerter, scanOptions, testFilter)
			got, err := s.Resources()
			assert.Equal(tt, err, c.err)
			if err != nil {
				return
			}
			test.TestAgainstGoldenFile(got, githubres.GithubBranchProtectionResourceType, c.dirName, provider, deserializer, shouldUpdate, tt)
			mockedRepo.AssertExpectations(tt)
			alerter.AssertExpectations(tt)
		})
	}
}
