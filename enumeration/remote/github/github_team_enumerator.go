package github

import (
	"fmt"
	remoteerror "github.com/snyk/driftctl/enumeration/remote/error"

	"github.com/snyk/driftctl/enumeration/resource"
	"github.com/snyk/driftctl/enumeration/resource/github"
)

type GithubTeamEnumerator struct {
	repository GithubRepository
	factory    resource.ResourceFactory
}

func NewGithubTeamEnumerator(repo GithubRepository, factory resource.ResourceFactory) *GithubTeamEnumerator {
	return &GithubTeamEnumerator{
		repository: repo,
		factory:    factory,
	}
}

func (g *GithubTeamEnumerator) SupportedType() resource.ResourceType {
	return github.GithubTeamResourceType
}

func (g *GithubTeamEnumerator) Enumerate() ([]*resource.Resource, error) {
	resourceList, err := g.repository.ListTeams()
	if err != nil {
		return nil, remoteerror.NewResourceListingError(err, string(g.SupportedType()))
	}

	results := make([]*resource.Resource, 0, len(resourceList))

	for _, team := range resourceList {
		results = append(
			results,
			g.factory.CreateAbstractResource(
				string(g.SupportedType()),
				fmt.Sprintf("%d", team.DatabaseId),
				map[string]interface{}{},
			),
		)
	}

	return results, err
}
