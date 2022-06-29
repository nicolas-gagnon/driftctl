package repository

import (
	"fmt"
	"github.com/snyk/driftctl/enumeration/remote/cache"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	awstest "github.com/snyk/driftctl/test/aws"

	"github.com/stretchr/testify/mock"

	"github.com/r3labs/diff/v2"
	"github.com/stretchr/testify/assert"
)

func Test_IAMRepository_ListAllAccessKeys(t *testing.T) {
	tests := []struct {
		name    string
		users   []*iam.User
		mocks   func(client *awstest.MockFakeIAM)
		want    []*iam.AccessKeyMetadata
		wantErr error
	}{
		{
			name: "List only access keys with multiple pages",
			users: []*iam.User{
				{
					UserName: aws.String("test-driftctl"),
				},
				{
					UserName: aws.String("test-driftctl2"),
				},
			},
			mocks: func(client *awstest.MockFakeIAM) {

				client.On("ListAccessKeysPages",
					&iam.ListAccessKeysInput{
						UserName: aws.String("test-driftctl"),
					},
					mock.MatchedBy(func(callback func(res *iam.ListAccessKeysOutput, lastPage bool) bool) bool {
						callback(&iam.ListAccessKeysOutput{AccessKeyMetadata: []*iam.AccessKeyMetadata{
							{
								AccessKeyId: aws.String("AKIA5QYBVVD223VWU32A"),
								UserName:    aws.String("test-driftctl"),
							},
						}}, false)
						callback(&iam.ListAccessKeysOutput{AccessKeyMetadata: []*iam.AccessKeyMetadata{
							{
								AccessKeyId: aws.String("AKIA5QYBVVD2QYI36UZP"),
								UserName:    aws.String("test-driftctl"),
							},
						}}, true)
						return true
					})).Return(nil).Once()
				client.On("ListAccessKeysPages",
					&iam.ListAccessKeysInput{
						UserName: aws.String("test-driftctl2"),
					},
					mock.MatchedBy(func(callback func(res *iam.ListAccessKeysOutput, lastPage bool) bool) bool {
						callback(&iam.ListAccessKeysOutput{AccessKeyMetadata: []*iam.AccessKeyMetadata{
							{
								AccessKeyId: aws.String("AKIA5QYBVVD26EJME25D"),
								UserName:    aws.String("test-driftctl2"),
							},
						}}, false)
						callback(&iam.ListAccessKeysOutput{AccessKeyMetadata: []*iam.AccessKeyMetadata{
							{
								AccessKeyId: aws.String("AKIA5QYBVVD2SWDFVVMG"),
								UserName:    aws.String("test-driftctl2"),
							},
						}}, true)
						return true
					})).Return(nil).Once()
			},
			want: []*iam.AccessKeyMetadata{
				{
					AccessKeyId: aws.String("AKIA5QYBVVD223VWU32A"),
					UserName:    aws.String("test-driftctl"),
				},
				{
					AccessKeyId: aws.String("AKIA5QYBVVD2QYI36UZP"),
					UserName:    aws.String("test-driftctl"),
				},
				{
					AccessKeyId: aws.String("AKIA5QYBVVD223VWU32A"),
					UserName:    aws.String("test-driftctl"),
				},
				{
					AccessKeyId: aws.String("AKIA5QYBVVD2QYI36UZP"),
					UserName:    aws.String("test-driftctl"),
				},
				{
					AccessKeyId: aws.String("AKIA5QYBVVD26EJME25D"),
					UserName:    aws.String("test-driftctl2"),
				},
				{
					AccessKeyId: aws.String("AKIA5QYBVVD2SWDFVVMG"),
					UserName:    aws.String("test-driftctl2"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := cache.New(2)
			client := &awstest.MockFakeIAM{}
			tt.mocks(client)
			r := &iamRepository{
				client: client,
				cache:  store,
			}
			got, err := r.ListAllAccessKeys(tt.users)
			assert.Equal(t, tt.wantErr, err)

			if err == nil {
				// Check that results were cached
				cachedData, err := r.ListAllAccessKeys(tt.users)
				assert.NoError(t, err)
				assert.Equal(t, got, cachedData)
				for _, user := range tt.users {
					assert.IsType(t, []*iam.AccessKeyMetadata{}, store.Get(fmt.Sprintf("iamListAllAccessKeys_user_%s", *user.UserName)))
				}
			}

			changelog, err := diff.Diff(got, tt.want)
			assert.Nil(t, err)
			if len(changelog) > 0 {
				for _, change := range changelog {
					t.Errorf("%s: %v -> %v", strings.Join(change.Path, "."), change.From, change.To)
				}
				t.Fail()
			}
		})
	}
}

func Test_IAMRepository_ListAllUsers(t *testing.T) {
	tests := []struct {
		name    string
		mocks   func(client *awstest.MockFakeIAM)
		want    []*iam.User
		wantErr error
	}{
		{
			name: "List only users with multiple pages",
			mocks: func(client *awstest.MockFakeIAM) {

				client.On("ListUsersPages",
					&iam.ListUsersInput{},
					mock.MatchedBy(func(callback func(res *iam.ListUsersOutput, lastPage bool) bool) bool {
						callback(&iam.ListUsersOutput{Users: []*iam.User{
							{
								UserName: aws.String("test-driftctl"),
							},
							{
								UserName: aws.String("test-driftctl2"),
							},
						}}, false)
						callback(&iam.ListUsersOutput{Users: []*iam.User{
							{
								UserName: aws.String("test-driftctl3"),
							},
							{
								UserName: aws.String("test-driftctl4"),
							},
						}}, true)
						return true
					})).Return(nil).Once()
			},
			want: []*iam.User{
				{
					UserName: aws.String("test-driftctl"),
				},
				{
					UserName: aws.String("test-driftctl2"),
				},
				{
					UserName: aws.String("test-driftctl3"),
				},
				{
					UserName: aws.String("test-driftctl4"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := cache.New(1)
			client := &awstest.MockFakeIAM{}
			tt.mocks(client)
			r := &iamRepository{
				client: client,
				cache:  store,
			}
			got, err := r.ListAllUsers()
			assert.Equal(t, tt.wantErr, err)

			if err == nil {
				// Check that results were cached
				cachedData, err := r.ListAllUsers()
				assert.NoError(t, err)
				assert.Equal(t, got, cachedData)
				assert.IsType(t, []*iam.User{}, store.Get("iamListAllUsers"))
			}

			changelog, err := diff.Diff(got, tt.want)
			assert.Nil(t, err)
			if len(changelog) > 0 {
				for _, change := range changelog {
					t.Errorf("%s: %v -> %v", strings.Join(change.Path, "."), change.From, change.To)
				}
				t.Fail()
			}
		})
	}
}

func Test_IAMRepository_ListAllPolicies(t *testing.T) {
	tests := []struct {
		name    string
		mocks   func(client *awstest.MockFakeIAM)
		want    []*iam.Policy
		wantErr error
	}{
		{
			name: "List only policies with multiple pages",
			mocks: func(client *awstest.MockFakeIAM) {

				client.On("ListPoliciesPages",
					&iam.ListPoliciesInput{Scope: aws.String(iam.PolicyScopeTypeLocal)},
					mock.MatchedBy(func(callback func(res *iam.ListPoliciesOutput, lastPage bool) bool) bool {
						callback(&iam.ListPoliciesOutput{Policies: []*iam.Policy{
							{
								PolicyName: aws.String("test-driftctl"),
							},
							{
								PolicyName: aws.String("test-driftctl2"),
							},
						}}, false)
						callback(&iam.ListPoliciesOutput{Policies: []*iam.Policy{
							{
								PolicyName: aws.String("test-driftctl3"),
							},
							{
								PolicyName: aws.String("test-driftctl4"),
							},
						}}, true)
						return true
					})).Return(nil).Once()
			},
			want: []*iam.Policy{
				{
					PolicyName: aws.String("test-driftctl"),
				},
				{
					PolicyName: aws.String("test-driftctl2"),
				},
				{
					PolicyName: aws.String("test-driftctl3"),
				},
				{
					PolicyName: aws.String("test-driftctl4"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := cache.New(1)
			client := &awstest.MockFakeIAM{}
			tt.mocks(client)
			r := &iamRepository{
				client: client,
				cache:  store,
			}
			got, err := r.ListAllPolicies()
			assert.Equal(t, tt.wantErr, err)

			if err == nil {
				// Check that results were cached
				cachedData, err := r.ListAllPolicies()
				assert.NoError(t, err)
				assert.Equal(t, got, cachedData)
				assert.IsType(t, []*iam.Policy{}, store.Get("iamListAllPolicies"))
			}

			changelog, err := diff.Diff(got, tt.want)
			assert.Nil(t, err)
			if len(changelog) > 0 {
				for _, change := range changelog {
					t.Errorf("%s: %v -> %v", strings.Join(change.Path, "."), change.From, change.To)
				}
				t.Fail()
			}
		})
	}
}

func Test_IAMRepository_ListAllRoles(t *testing.T) {
	tests := []struct {
		name    string
		mocks   func(client *awstest.MockFakeIAM)
		want    []*iam.Role
		wantErr error
	}{
		{
			name: "List only roles with multiple pages",
			mocks: func(client *awstest.MockFakeIAM) {

				client.On("ListRolesPages",
					&iam.ListRolesInput{},
					mock.MatchedBy(func(callback func(res *iam.ListRolesOutput, lastPage bool) bool) bool {
						callback(&iam.ListRolesOutput{Roles: []*iam.Role{
							{
								RoleName: aws.String("test-driftctl"),
							},
							{
								RoleName: aws.String("test-driftctl2"),
							},
						}}, false)
						callback(&iam.ListRolesOutput{Roles: []*iam.Role{
							{
								RoleName: aws.String("test-driftctl3"),
							},
							{
								RoleName: aws.String("test-driftctl4"),
							},
						}}, true)
						return true
					})).Return(nil).Once()
			},
			want: []*iam.Role{
				{
					RoleName: aws.String("test-driftctl"),
				},
				{
					RoleName: aws.String("test-driftctl2"),
				},
				{
					RoleName: aws.String("test-driftctl3"),
				},
				{
					RoleName: aws.String("test-driftctl4"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := cache.New(1)
			client := &awstest.MockFakeIAM{}
			tt.mocks(client)
			r := &iamRepository{
				client: client,
				cache:  store,
			}
			got, err := r.ListAllRoles()
			assert.Equal(t, tt.wantErr, err)

			if err == nil {
				// Check that results were cached
				cachedData, err := r.ListAllRoles()
				assert.NoError(t, err)
				assert.Equal(t, got, cachedData)
				assert.IsType(t, []*iam.Role{}, store.Get("iamListAllRoles"))
			}

			changelog, err := diff.Diff(got, tt.want)
			assert.Nil(t, err)
			if len(changelog) > 0 {
				for _, change := range changelog {
					t.Errorf("%s: %v -> %v", strings.Join(change.Path, "."), change.From, change.To)
				}
				t.Fail()
			}
		})
	}
}

func Test_IAMRepository_ListAllRolePolicyAttachments(t *testing.T) {
	tests := []struct {
		name    string
		roles   []*iam.Role
		mocks   func(client *awstest.MockFakeIAM)
		want    []*AttachedRolePolicy
		wantErr error
	}{
		{
			name: "List only role policy attachments with multiple pages",
			roles: []*iam.Role{
				{
					RoleName: aws.String("test-role"),
				},
				{
					RoleName: aws.String("test-role2"),
				},
			},
			mocks: func(client *awstest.MockFakeIAM) {

				shouldSkipfirst := false
				shouldSkipSecond := false

				client.On("ListAttachedRolePoliciesPages",
					&iam.ListAttachedRolePoliciesInput{
						RoleName: aws.String("test-role"),
					},
					mock.MatchedBy(func(callback func(res *iam.ListAttachedRolePoliciesOutput, lastPage bool) bool) bool {
						if shouldSkipfirst {
							return false
						}
						callback(&iam.ListAttachedRolePoliciesOutput{AttachedPolicies: []*iam.AttachedPolicy{
							{
								PolicyArn:  aws.String("arn:aws:iam::526954929923:policy/test-policy"),
								PolicyName: aws.String("policy"),
							},
							{
								PolicyArn:  aws.String("arn:aws:iam::526954929923:policy/test-policy2"),
								PolicyName: aws.String("policy2"),
							},
						}}, false)
						callback(&iam.ListAttachedRolePoliciesOutput{AttachedPolicies: []*iam.AttachedPolicy{
							{
								PolicyArn:  aws.String("arn:aws:iam::526954929923:policy/test-policy3"),
								PolicyName: aws.String("policy3"),
							},
						}}, true)
						shouldSkipfirst = true
						return true
					})).Return(nil).Once()

				client.On("ListAttachedRolePoliciesPages",
					&iam.ListAttachedRolePoliciesInput{
						RoleName: aws.String("test-role2"),
					},
					mock.MatchedBy(func(callback func(res *iam.ListAttachedRolePoliciesOutput, lastPage bool) bool) bool {
						if shouldSkipSecond {
							return false
						}
						callback(&iam.ListAttachedRolePoliciesOutput{AttachedPolicies: []*iam.AttachedPolicy{
							{
								PolicyArn:  aws.String("arn:aws:iam::526954929923:policy/test-policy"),
								PolicyName: aws.String("policy"),
							},
							{
								PolicyArn:  aws.String("arn:aws:iam::526954929923:policy/test-policy2"),
								PolicyName: aws.String("policy2"),
							},
						}}, false)
						callback(&iam.ListAttachedRolePoliciesOutput{AttachedPolicies: []*iam.AttachedPolicy{
							{
								PolicyArn:  aws.String("arn:aws:iam::526954929923:policy/test-policy3"),
								PolicyName: aws.String("policy3"),
							},
						}}, true)
						shouldSkipSecond = true
						return true
					})).Return(nil).Once()
			},
			want: []*AttachedRolePolicy{
				{
					iam.AttachedPolicy{
						PolicyArn:  aws.String("arn:aws:iam::526954929923:policy/test-policy"),
						PolicyName: aws.String("policy"),
					},
					*aws.String("test-role"),
				},
				{
					iam.AttachedPolicy{
						PolicyArn:  aws.String("arn:aws:iam::526954929923:policy/test-policy2"),
						PolicyName: aws.String("policy2"),
					},
					*aws.String("test-role"),
				},
				{
					iam.AttachedPolicy{
						PolicyArn:  aws.String("arn:aws:iam::526954929923:policy/test-policy3"),
						PolicyName: aws.String("policy3"),
					},
					*aws.String("test-role"),
				},
				{
					iam.AttachedPolicy{
						PolicyArn:  aws.String("arn:aws:iam::526954929923:policy/test-policy"),
						PolicyName: aws.String("policy"),
					},
					*aws.String("test-role2"),
				},
				{
					iam.AttachedPolicy{
						PolicyArn:  aws.String("arn:aws:iam::526954929923:policy/test-policy2"),
						PolicyName: aws.String("policy2"),
					},
					*aws.String("test-role2"),
				},
				{
					iam.AttachedPolicy{
						PolicyArn:  aws.String("arn:aws:iam::526954929923:policy/test-policy3"),
						PolicyName: aws.String("policy3"),
					},
					*aws.String("test-role2"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := cache.New(2)
			client := &awstest.MockFakeIAM{}
			tt.mocks(client)
			r := &iamRepository{
				client: client,
				cache:  store,
			}
			got, err := r.ListAllRolePolicyAttachments(tt.roles)
			assert.Equal(t, tt.wantErr, err)

			if err == nil {
				// Check that results were cached
				cachedData, err := r.ListAllRolePolicyAttachments(tt.roles)
				assert.NoError(t, err)
				assert.Equal(t, got, cachedData)
				for _, role := range tt.roles {
					assert.IsType(t, []*AttachedRolePolicy{}, store.Get(fmt.Sprintf("iamListAllRolePolicyAttachments_role_%s", *role.RoleName)))
				}
			}

			changelog, err := diff.Diff(got, tt.want)
			assert.Nil(t, err)
			if len(changelog) > 0 {
				for _, change := range changelog {
					t.Errorf("%s: %v -> %v", strings.Join(change.Path, "."), change.From, change.To)
				}
				t.Fail()
			}
		})
	}
}

func Test_IAMRepository_ListAllRolePolicies(t *testing.T) {
	tests := []struct {
		name    string
		roles   []*iam.Role
		mocks   func(client *awstest.MockFakeIAM)
		want    []RolePolicy
		wantErr error
	}{
		{
			name: "List only role policies with multiple pages",
			roles: []*iam.Role{
				{
					RoleName: aws.String("test_role_0"),
				},
				{
					RoleName: aws.String("test_role_1"),
				},
			},
			mocks: func(client *awstest.MockFakeIAM) {
				firstMockCalled := false
				client.On("ListRolePoliciesPages",
					&iam.ListRolePoliciesInput{
						RoleName: aws.String("test_role_0"),
					},
					mock.MatchedBy(func(callback func(res *iam.ListRolePoliciesOutput, lastPage bool) bool) bool {
						if firstMockCalled {
							return false
						}
						callback(&iam.ListRolePoliciesOutput{
							PolicyNames: []*string{
								aws.String("policy-role0-0"),
								aws.String("policy-role0-1"),
							},
						}, false)
						callback(&iam.ListRolePoliciesOutput{
							PolicyNames: []*string{
								aws.String("policy-role0-2"),
							},
						}, true)
						firstMockCalled = true
						return true
					})).Once().Return(nil)
				client.On("ListRolePoliciesPages",
					&iam.ListRolePoliciesInput{
						RoleName: aws.String("test_role_1"),
					},
					mock.MatchedBy(func(callback func(res *iam.ListRolePoliciesOutput, lastPage bool) bool) bool {
						callback(&iam.ListRolePoliciesOutput{
							PolicyNames: []*string{
								aws.String("policy-role1-0"),
								aws.String("policy-role1-1"),
							},
						}, false)
						callback(&iam.ListRolePoliciesOutput{
							PolicyNames: []*string{
								aws.String("policy-role1-2"),
							},
						}, true)
						return true
					})).Once().Return(nil)
			},
			want: []RolePolicy{
				{Policy: "policy-role0-0", RoleName: "test_role_0"},
				{Policy: "policy-role0-1", RoleName: "test_role_0"},
				{Policy: "policy-role0-2", RoleName: "test_role_0"},
				{Policy: "policy-role1-0", RoleName: "test_role_1"},
				{Policy: "policy-role1-1", RoleName: "test_role_1"},
				{Policy: "policy-role1-2", RoleName: "test_role_1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := cache.New(2)
			client := &awstest.MockFakeIAM{}
			tt.mocks(client)
			r := &iamRepository{
				client: client,
				cache:  store,
			}
			got, err := r.ListAllRolePolicies(tt.roles)
			assert.Equal(t, tt.wantErr, err)

			if err == nil {
				// Check that results were cached
				cachedData, err := r.ListAllRolePolicies(tt.roles)
				assert.NoError(t, err)
				assert.Equal(t, got, cachedData)
				for _, role := range tt.roles {
					assert.IsType(t, []RolePolicy{}, store.Get(fmt.Sprintf("iamListAllRolePolicies_role_%s", *role.RoleName)))
				}
			}

			changelog, err := diff.Diff(got, tt.want)
			assert.Nil(t, err)
			if len(changelog) > 0 {
				for _, change := range changelog {
					t.Errorf("%s: %v -> %v", strings.Join(change.Path, "."), change.From, change.To)
				}
				t.Fail()
			}
		})
	}
}

func Test_IAMRepository_ListAllUserPolicyAttachments(t *testing.T) {
	tests := []struct {
		name    string
		users   []*iam.User
		mocks   func(client *awstest.MockFakeIAM)
		want    []*AttachedUserPolicy
		wantErr error
	}{
		{
			name: "List only user policy attachments with multiple pages",
			users: []*iam.User{
				{
					UserName: aws.String("loadbalancer"),
				},
				{
					UserName: aws.String("loadbalancer2"),
				},
			},
			mocks: func(client *awstest.MockFakeIAM) {

				client.On("ListAttachedUserPoliciesPages",
					&iam.ListAttachedUserPoliciesInput{
						UserName: aws.String("loadbalancer"),
					},
					mock.MatchedBy(func(callback func(res *iam.ListAttachedUserPoliciesOutput, lastPage bool) bool) bool {
						callback(&iam.ListAttachedUserPoliciesOutput{AttachedPolicies: []*iam.AttachedPolicy{
							{
								PolicyArn:  aws.String("arn:aws:iam::526954929923:policy/test"),
								PolicyName: aws.String("test-attach"),
							},
							{
								PolicyArn:  aws.String("arn:aws:iam::526954929923:policy/test2"),
								PolicyName: aws.String("test-attach2"),
							},
						}}, false)
						callback(&iam.ListAttachedUserPoliciesOutput{AttachedPolicies: []*iam.AttachedPolicy{
							{
								PolicyArn:  aws.String("arn:aws:iam::526954929923:policy/test3"),
								PolicyName: aws.String("test-attach3"),
							},
						}}, true)
						return true
					})).Return(nil).Once()

				client.On("ListAttachedUserPoliciesPages",
					&iam.ListAttachedUserPoliciesInput{
						UserName: aws.String("loadbalancer2"),
					},
					mock.MatchedBy(func(callback func(res *iam.ListAttachedUserPoliciesOutput, lastPage bool) bool) bool {
						callback(&iam.ListAttachedUserPoliciesOutput{AttachedPolicies: []*iam.AttachedPolicy{
							{
								PolicyArn:  aws.String("arn:aws:iam::526954929923:policy/test"),
								PolicyName: aws.String("test-attach"),
							},
							{
								PolicyArn:  aws.String("arn:aws:iam::526954929923:policy/test2"),
								PolicyName: aws.String("test-attach2"),
							},
						}}, false)
						callback(&iam.ListAttachedUserPoliciesOutput{AttachedPolicies: []*iam.AttachedPolicy{
							{
								PolicyArn:  aws.String("arn:aws:iam::526954929923:policy/test3"),
								PolicyName: aws.String("test-attach3"),
							},
						}}, true)
						return true
					})).Return(nil).Once()
			},

			want: []*AttachedUserPolicy{
				{
					iam.AttachedPolicy{
						PolicyArn:  aws.String("arn:aws:iam::526954929923:policy/test"),
						PolicyName: aws.String("test-attach"),
					},
					*aws.String("loadbalancer"),
				},
				{
					iam.AttachedPolicy{
						PolicyArn:  aws.String("arn:aws:iam::526954929923:policy/test2"),
						PolicyName: aws.String("test-attach2"),
					},
					*aws.String("loadbalancer"),
				},
				{
					iam.AttachedPolicy{
						PolicyArn:  aws.String("arn:aws:iam::526954929923:policy/test3"),
						PolicyName: aws.String("test-attach3"),
					},
					*aws.String("loadbalancer"),
				},
				{
					iam.AttachedPolicy{
						PolicyArn:  aws.String("arn:aws:iam::526954929923:policy/test"),
						PolicyName: aws.String("test-attach"),
					},
					*aws.String("loadbalancer2"),
				},
				{
					iam.AttachedPolicy{
						PolicyArn:  aws.String("arn:aws:iam::526954929923:policy/test2"),
						PolicyName: aws.String("test-attach2"),
					},
					*aws.String("loadbalancer2"),
				},
				{
					iam.AttachedPolicy{
						PolicyArn:  aws.String("arn:aws:iam::526954929923:policy/test3"),
						PolicyName: aws.String("test-attach3"),
					},
					*aws.String("loadbalancer2"),
				},
				{
					iam.AttachedPolicy{
						PolicyArn:  aws.String("arn:aws:iam::526954929923:policy/test"),
						PolicyName: aws.String("test-attach"),
					},
					*aws.String("loadbalancer2"),
				},
				{
					iam.AttachedPolicy{
						PolicyArn:  aws.String("arn:aws:iam::526954929923:policy/test2"),
						PolicyName: aws.String("test-attach2"),
					},
					*aws.String("loadbalancer2"),
				},
				{
					iam.AttachedPolicy{
						PolicyArn:  aws.String("arn:aws:iam::526954929923:policy/test3"),
						PolicyName: aws.String("test-attach3"),
					},
					*aws.String("loadbalancer2"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := cache.New(2)
			client := &awstest.MockFakeIAM{}
			tt.mocks(client)
			r := &iamRepository{
				client: client,
				cache:  store,
			}
			got, err := r.ListAllUserPolicyAttachments(tt.users)
			assert.Equal(t, tt.wantErr, err)

			if err == nil {
				// Check that results were cached
				cachedData, err := r.ListAllUserPolicyAttachments(tt.users)
				assert.NoError(t, err)
				assert.Equal(t, got, cachedData)
				for _, user := range tt.users {
					assert.IsType(t, []*AttachedUserPolicy{}, store.Get(fmt.Sprintf("iamListAllUserPolicyAttachments_user_%s", *user.UserName)))
				}
			}

			changelog, err := diff.Diff(got, tt.want)
			assert.Nil(t, err)
			if len(changelog) > 0 {
				for _, change := range changelog {
					t.Errorf("%s: %v -> %v", strings.Join(change.Path, "."), change.From, change.To)
				}
				t.Fail()
			}
		})
	}
}

func Test_IAMRepository_ListAllUserPolicies(t *testing.T) {
	tests := []struct {
		name    string
		users   []*iam.User
		mocks   func(client *awstest.MockFakeIAM)
		want    []string
		wantErr error
	}{
		{
			name: "List only user policies with multiple pages",
			users: []*iam.User{
				{
					UserName: aws.String("loadbalancer"),
				},
				{
					UserName: aws.String("loadbalancer2"),
				},
			},
			mocks: func(client *awstest.MockFakeIAM) {

				client.On("ListUserPoliciesPages",
					&iam.ListUserPoliciesInput{
						UserName: aws.String("loadbalancer"),
					},
					mock.MatchedBy(func(callback func(res *iam.ListUserPoliciesOutput, lastPage bool) bool) bool {
						callback(&iam.ListUserPoliciesOutput{PolicyNames: []*string{
							aws.String("test"),
							aws.String("test2"),
							aws.String("test3"),
						}}, false)
						callback(&iam.ListUserPoliciesOutput{PolicyNames: []*string{
							aws.String("test4"),
						}}, true)
						return true
					})).Return(nil).Once()

				client.On("ListUserPoliciesPages",
					&iam.ListUserPoliciesInput{
						UserName: aws.String("loadbalancer2"),
					},
					mock.MatchedBy(func(callback func(res *iam.ListUserPoliciesOutput, lastPage bool) bool) bool {
						callback(&iam.ListUserPoliciesOutput{PolicyNames: []*string{
							aws.String("test2"),
							aws.String("test22"),
							aws.String("test23"),
						}}, false)
						callback(&iam.ListUserPoliciesOutput{PolicyNames: []*string{
							aws.String("test24"),
						}}, true)
						return true
					})).Return(nil).Once()
			},
			want: []string{
				*aws.String("loadbalancer:test"),
				*aws.String("loadbalancer:test2"),
				*aws.String("loadbalancer:test3"),
				*aws.String("loadbalancer:test4"),
				*aws.String("loadbalancer2:test"),
				*aws.String("loadbalancer2:test2"),
				*aws.String("loadbalancer2:test3"),
				*aws.String("loadbalancer2:test4"),
				*aws.String("loadbalancer2:test2"),
				*aws.String("loadbalancer2:test22"),
				*aws.String("loadbalancer2:test23"),
				*aws.String("loadbalancer2:test24"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := cache.New(2)
			client := &awstest.MockFakeIAM{}
			tt.mocks(client)
			r := &iamRepository{
				client: client,
				cache:  store,
			}
			got, err := r.ListAllUserPolicies(tt.users)
			assert.Equal(t, tt.wantErr, err)

			if err == nil {
				// Check that results were cached
				cachedData, err := r.ListAllUserPolicies(tt.users)
				assert.NoError(t, err)
				assert.Equal(t, got, cachedData)
				for _, user := range tt.users {
					assert.IsType(t, []string{}, store.Get(fmt.Sprintf("iamListAllUserPolicies_user_%s", *user.UserName)))
				}
			}

			changelog, err := diff.Diff(got, tt.want)
			assert.Nil(t, err)
			if len(changelog) > 0 {
				for _, change := range changelog {
					t.Errorf("%s: %v -> %v", strings.Join(change.Path, "."), change.From, change.To)
				}
				t.Fail()
			}
		})
	}
}

func Test_IAMRepository_ListAllGroups(t *testing.T) {
	tests := []struct {
		name    string
		mocks   func(client *awstest.MockFakeIAM)
		want    []*iam.Group
		wantErr error
	}{
		{
			name: "List groups with multiple pages",
			mocks: func(client *awstest.MockFakeIAM) {

				client.On("ListGroupsPages",
					&iam.ListGroupsInput{},
					mock.MatchedBy(func(callback func(res *iam.ListGroupsOutput, lastPage bool) bool) bool {
						callback(&iam.ListGroupsOutput{Groups: []*iam.Group{
							{
								GroupName: aws.String("group1"),
							},
							{
								GroupName: aws.String("group2"),
							},
						}}, false)
						callback(&iam.ListGroupsOutput{Groups: []*iam.Group{
							{
								GroupName: aws.String("group3"),
							},
							{
								GroupName: aws.String("group4"),
							},
						}}, true)
						return true
					})).Return(nil).Once()
			},
			want: []*iam.Group{
				{
					GroupName: aws.String("group1"),
				},
				{
					GroupName: aws.String("group2"),
				},
				{
					GroupName: aws.String("group3"),
				},
				{
					GroupName: aws.String("group4"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := cache.New(1)
			client := &awstest.MockFakeIAM{}
			tt.mocks(client)
			r := &iamRepository{
				client: client,
				cache:  store,
			}
			got, err := r.ListAllGroups()
			assert.Equal(t, tt.wantErr, err)

			if err == nil {
				// Check that results were cached
				cachedData, err := r.ListAllGroups()
				assert.NoError(t, err)
				assert.Equal(t, got, cachedData)
				assert.IsType(t, []*iam.Group{}, store.Get("iamListAllGroups"))
			}

			changelog, err := diff.Diff(got, tt.want)
			assert.Nil(t, err)
			if len(changelog) > 0 {
				for _, change := range changelog {
					t.Errorf("%s: %v -> %v", strings.Join(change.Path, "."), change.From, change.To)
				}
				t.Fail()
			}
		})
	}
}

func Test_IAMRepository_ListAllGroupPolicies(t *testing.T) {
	tests := []struct {
		name    string
		groups  []*iam.Group
		mocks   func(client *awstest.MockFakeIAM)
		want    []string
		wantErr error
	}{
		{
			name: "List only group policies with multiple pages",
			groups: []*iam.Group{
				{
					GroupName: aws.String("group1"),
				},
				{
					GroupName: aws.String("group2"),
				},
			},
			mocks: func(client *awstest.MockFakeIAM) {
				firstMockCalled := false
				client.On("ListGroupPoliciesPages",
					&iam.ListGroupPoliciesInput{
						GroupName: aws.String("group1"),
					},
					mock.MatchedBy(func(callback func(res *iam.ListGroupPoliciesOutput, lastPage bool) bool) bool {
						if firstMockCalled {
							return false
						}
						callback(&iam.ListGroupPoliciesOutput{PolicyNames: []*string{
							aws.String("policy1"),
							aws.String("policy2"),
							aws.String("policy3"),
						}}, false)
						callback(&iam.ListGroupPoliciesOutput{PolicyNames: []*string{
							aws.String("policy4"),
						}}, true)
						firstMockCalled = true
						return true
					})).Return(nil).Once()

				client.On("ListGroupPoliciesPages",
					&iam.ListGroupPoliciesInput{
						GroupName: aws.String("group2"),
					},
					mock.MatchedBy(func(callback func(res *iam.ListGroupPoliciesOutput, lastPage bool) bool) bool {
						callback(&iam.ListGroupPoliciesOutput{PolicyNames: []*string{
							aws.String("policy2"),
							aws.String("policy22"),
							aws.String("policy23"),
						}}, false)
						callback(&iam.ListGroupPoliciesOutput{PolicyNames: []*string{
							aws.String("policy24"),
						}}, true)
						return true
					})).Return(nil).Once()
			},
			want: []string{
				*aws.String("group1:policy1"),
				*aws.String("group1:policy2"),
				*aws.String("group1:policy3"),
				*aws.String("group1:policy4"),
				*aws.String("group2:policy2"),
				*aws.String("group2:policy22"),
				*aws.String("group2:policy23"),
				*aws.String("group2:policy24"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := cache.New(2)
			client := &awstest.MockFakeIAM{}
			tt.mocks(client)
			r := &iamRepository{
				client: client,
				cache:  store,
			}
			got, err := r.ListAllGroupPolicies(tt.groups)
			assert.Equal(t, tt.wantErr, err)

			if err == nil {
				// Check that results were cached
				cachedData, err := r.ListAllGroupPolicies(tt.groups)
				assert.NoError(t, err)
				assert.Equal(t, got, cachedData)
				for _, group := range tt.groups {
					assert.IsType(t, []string{}, store.Get(fmt.Sprintf("iamListAllGroupPolicies_group_%s", *group.GroupName)))
				}
			}

			changelog, err := diff.Diff(got, tt.want)
			assert.Nil(t, err)
			if len(changelog) > 0 {
				for _, change := range changelog {
					t.Errorf("%s: %v -> %v", strings.Join(change.Path, "."), change.From, change.To)
				}
				t.Fail()
			}
		})
	}
}
