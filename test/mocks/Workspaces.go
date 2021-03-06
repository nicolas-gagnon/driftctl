// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"
	io "io"

	mock "github.com/stretchr/testify/mock"

	tfe "github.com/hashicorp/go-tfe"
)

// Workspaces is an autogenerated mock type for the Workspaces type
type Workspaces struct {
	mock.Mock
}

// AddRemoteStateConsumers provides a mock function with given fields: ctx, workspaceID, options
func (_m *Workspaces) AddRemoteStateConsumers(ctx context.Context, workspaceID string, options tfe.WorkspaceAddRemoteStateConsumersOptions) error {
	ret := _m.Called(ctx, workspaceID, options)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, tfe.WorkspaceAddRemoteStateConsumersOptions) error); ok {
		r0 = rf(ctx, workspaceID, options)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AddTags provides a mock function with given fields: ctx, workspaceID, options
func (_m *Workspaces) AddTags(ctx context.Context, workspaceID string, options tfe.WorkspaceAddTagsOptions) error {
	ret := _m.Called(ctx, workspaceID, options)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, tfe.WorkspaceAddTagsOptions) error); ok {
		r0 = rf(ctx, workspaceID, options)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AssignSSHKey provides a mock function with given fields: ctx, workspaceID, options
func (_m *Workspaces) AssignSSHKey(ctx context.Context, workspaceID string, options tfe.WorkspaceAssignSSHKeyOptions) (*tfe.Workspace, error) {
	ret := _m.Called(ctx, workspaceID, options)

	var r0 *tfe.Workspace
	if rf, ok := ret.Get(0).(func(context.Context, string, tfe.WorkspaceAssignSSHKeyOptions) *tfe.Workspace); ok {
		r0 = rf(ctx, workspaceID, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tfe.Workspace)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, tfe.WorkspaceAssignSSHKeyOptions) error); ok {
		r1 = rf(ctx, workspaceID, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Create provides a mock function with given fields: ctx, organization, options
func (_m *Workspaces) Create(ctx context.Context, organization string, options tfe.WorkspaceCreateOptions) (*tfe.Workspace, error) {
	ret := _m.Called(ctx, organization, options)

	var r0 *tfe.Workspace
	if rf, ok := ret.Get(0).(func(context.Context, string, tfe.WorkspaceCreateOptions) *tfe.Workspace); ok {
		r0 = rf(ctx, organization, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tfe.Workspace)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, tfe.WorkspaceCreateOptions) error); ok {
		r1 = rf(ctx, organization, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, organization, workspace
func (_m *Workspaces) Delete(ctx context.Context, organization string, workspace string) error {
	ret := _m.Called(ctx, organization, workspace)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, organization, workspace)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteByID provides a mock function with given fields: ctx, workspaceID
func (_m *Workspaces) DeleteByID(ctx context.Context, workspaceID string) error {
	ret := _m.Called(ctx, workspaceID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, workspaceID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ForceUnlock provides a mock function with given fields: ctx, workspaceID
func (_m *Workspaces) ForceUnlock(ctx context.Context, workspaceID string) (*tfe.Workspace, error) {
	ret := _m.Called(ctx, workspaceID)

	var r0 *tfe.Workspace
	if rf, ok := ret.Get(0).(func(context.Context, string) *tfe.Workspace); ok {
		r0 = rf(ctx, workspaceID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tfe.Workspace)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, workspaceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, organization, options
func (_m *Workspaces) List(ctx context.Context, organization string, options tfe.WorkspaceListOptions) (*tfe.WorkspaceList, error) {
	ret := _m.Called(ctx, organization, options)

	var r0 *tfe.WorkspaceList
	if rf, ok := ret.Get(0).(func(context.Context, string, tfe.WorkspaceListOptions) *tfe.WorkspaceList); ok {
		r0 = rf(ctx, organization, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tfe.WorkspaceList)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, tfe.WorkspaceListOptions) error); ok {
		r1 = rf(ctx, organization, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Lock provides a mock function with given fields: ctx, workspaceID, options
func (_m *Workspaces) Lock(ctx context.Context, workspaceID string, options tfe.WorkspaceLockOptions) (*tfe.Workspace, error) {
	ret := _m.Called(ctx, workspaceID, options)

	var r0 *tfe.Workspace
	if rf, ok := ret.Get(0).(func(context.Context, string, tfe.WorkspaceLockOptions) *tfe.Workspace); ok {
		r0 = rf(ctx, workspaceID, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tfe.Workspace)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, tfe.WorkspaceLockOptions) error); ok {
		r1 = rf(ctx, workspaceID, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Read provides a mock function with given fields: ctx, organization, workspace
func (_m *Workspaces) Read(ctx context.Context, organization string, workspace string) (*tfe.Workspace, error) {
	ret := _m.Called(ctx, organization, workspace)

	var r0 *tfe.Workspace
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *tfe.Workspace); ok {
		r0 = rf(ctx, organization, workspace)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tfe.Workspace)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, organization, workspace)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReadByID provides a mock function with given fields: ctx, workspaceID
func (_m *Workspaces) ReadByID(ctx context.Context, workspaceID string) (*tfe.Workspace, error) {
	ret := _m.Called(ctx, workspaceID)

	var r0 *tfe.Workspace
	if rf, ok := ret.Get(0).(func(context.Context, string) *tfe.Workspace); ok {
		r0 = rf(ctx, workspaceID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tfe.Workspace)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, workspaceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReadByIDWithOptions provides a mock function with given fields: ctx, workspaceID, options
func (_m *Workspaces) ReadByIDWithOptions(ctx context.Context, workspaceID string, options *tfe.WorkspaceReadOptions) (*tfe.Workspace, error) {
	ret := _m.Called(ctx, workspaceID, options)

	var r0 *tfe.Workspace
	if rf, ok := ret.Get(0).(func(context.Context, string, *tfe.WorkspaceReadOptions) *tfe.Workspace); ok {
		r0 = rf(ctx, workspaceID, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tfe.Workspace)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, *tfe.WorkspaceReadOptions) error); ok {
		r1 = rf(ctx, workspaceID, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReadWithOptions provides a mock function with given fields: ctx, organization, workspace, options
func (_m *Workspaces) ReadWithOptions(ctx context.Context, organization string, workspace string, options *tfe.WorkspaceReadOptions) (*tfe.Workspace, error) {
	ret := _m.Called(ctx, organization, workspace, options)

	var r0 *tfe.Workspace
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *tfe.WorkspaceReadOptions) *tfe.Workspace); ok {
		r0 = rf(ctx, organization, workspace, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tfe.Workspace)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, *tfe.WorkspaceReadOptions) error); ok {
		r1 = rf(ctx, organization, workspace, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Readme provides a mock function with given fields: ctx, workspaceID
func (_m *Workspaces) Readme(ctx context.Context, workspaceID string) (io.Reader, error) {
	ret := _m.Called(ctx, workspaceID)

	var r0 io.Reader
	if rf, ok := ret.Get(0).(func(context.Context, string) io.Reader); ok {
		r0 = rf(ctx, workspaceID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(io.Reader)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, workspaceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoteStateConsumers provides a mock function with given fields: ctx, workspaceID, options
func (_m *Workspaces) RemoteStateConsumers(ctx context.Context, workspaceID string, options *tfe.RemoteStateConsumersListOptions) (*tfe.WorkspaceList, error) {
	ret := _m.Called(ctx, workspaceID, options)

	var r0 *tfe.WorkspaceList
	if rf, ok := ret.Get(0).(func(context.Context, string, *tfe.RemoteStateConsumersListOptions) *tfe.WorkspaceList); ok {
		r0 = rf(ctx, workspaceID, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tfe.WorkspaceList)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, *tfe.RemoteStateConsumersListOptions) error); ok {
		r1 = rf(ctx, workspaceID, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveRemoteStateConsumers provides a mock function with given fields: ctx, workspaceID, options
func (_m *Workspaces) RemoveRemoteStateConsumers(ctx context.Context, workspaceID string, options tfe.WorkspaceRemoveRemoteStateConsumersOptions) error {
	ret := _m.Called(ctx, workspaceID, options)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, tfe.WorkspaceRemoveRemoteStateConsumersOptions) error); ok {
		r0 = rf(ctx, workspaceID, options)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveTags provides a mock function with given fields: ctx, workspaceID, options
func (_m *Workspaces) RemoveTags(ctx context.Context, workspaceID string, options tfe.WorkspaceRemoveTagsOptions) error {
	ret := _m.Called(ctx, workspaceID, options)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, tfe.WorkspaceRemoveTagsOptions) error); ok {
		r0 = rf(ctx, workspaceID, options)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveVCSConnection provides a mock function with given fields: ctx, organization, workspace
func (_m *Workspaces) RemoveVCSConnection(ctx context.Context, organization string, workspace string) (*tfe.Workspace, error) {
	ret := _m.Called(ctx, organization, workspace)

	var r0 *tfe.Workspace
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *tfe.Workspace); ok {
		r0 = rf(ctx, organization, workspace)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tfe.Workspace)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, organization, workspace)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveVCSConnectionByID provides a mock function with given fields: ctx, workspaceID
func (_m *Workspaces) RemoveVCSConnectionByID(ctx context.Context, workspaceID string) (*tfe.Workspace, error) {
	ret := _m.Called(ctx, workspaceID)

	var r0 *tfe.Workspace
	if rf, ok := ret.Get(0).(func(context.Context, string) *tfe.Workspace); ok {
		r0 = rf(ctx, workspaceID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tfe.Workspace)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, workspaceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Tags provides a mock function with given fields: ctx, workspaceID, options
func (_m *Workspaces) Tags(ctx context.Context, workspaceID string, options tfe.WorkspaceTagListOptions) (*tfe.TagList, error) {
	ret := _m.Called(ctx, workspaceID, options)

	var r0 *tfe.TagList
	if rf, ok := ret.Get(0).(func(context.Context, string, tfe.WorkspaceTagListOptions) *tfe.TagList); ok {
		r0 = rf(ctx, workspaceID, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tfe.TagList)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, tfe.WorkspaceTagListOptions) error); ok {
		r1 = rf(ctx, workspaceID, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UnassignSSHKey provides a mock function with given fields: ctx, workspaceID
func (_m *Workspaces) UnassignSSHKey(ctx context.Context, workspaceID string) (*tfe.Workspace, error) {
	ret := _m.Called(ctx, workspaceID)

	var r0 *tfe.Workspace
	if rf, ok := ret.Get(0).(func(context.Context, string) *tfe.Workspace); ok {
		r0 = rf(ctx, workspaceID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tfe.Workspace)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, workspaceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Unlock provides a mock function with given fields: ctx, workspaceID
func (_m *Workspaces) Unlock(ctx context.Context, workspaceID string) (*tfe.Workspace, error) {
	ret := _m.Called(ctx, workspaceID)

	var r0 *tfe.Workspace
	if rf, ok := ret.Get(0).(func(context.Context, string) *tfe.Workspace); ok {
		r0 = rf(ctx, workspaceID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tfe.Workspace)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, workspaceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, organization, workspace, options
func (_m *Workspaces) Update(ctx context.Context, organization string, workspace string, options tfe.WorkspaceUpdateOptions) (*tfe.Workspace, error) {
	ret := _m.Called(ctx, organization, workspace, options)

	var r0 *tfe.Workspace
	if rf, ok := ret.Get(0).(func(context.Context, string, string, tfe.WorkspaceUpdateOptions) *tfe.Workspace); ok {
		r0 = rf(ctx, organization, workspace, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tfe.Workspace)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, tfe.WorkspaceUpdateOptions) error); ok {
		r1 = rf(ctx, organization, workspace, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateByID provides a mock function with given fields: ctx, workspaceID, options
func (_m *Workspaces) UpdateByID(ctx context.Context, workspaceID string, options tfe.WorkspaceUpdateOptions) (*tfe.Workspace, error) {
	ret := _m.Called(ctx, workspaceID, options)

	var r0 *tfe.Workspace
	if rf, ok := ret.Get(0).(func(context.Context, string, tfe.WorkspaceUpdateOptions) *tfe.Workspace); ok {
		r0 = rf(ctx, workspaceID, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tfe.Workspace)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, tfe.WorkspaceUpdateOptions) error); ok {
		r1 = rf(ctx, workspaceID, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateRemoteStateConsumers provides a mock function with given fields: ctx, workspaceID, options
func (_m *Workspaces) UpdateRemoteStateConsumers(ctx context.Context, workspaceID string, options tfe.WorkspaceUpdateRemoteStateConsumersOptions) error {
	ret := _m.Called(ctx, workspaceID, options)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, tfe.WorkspaceUpdateRemoteStateConsumersOptions) error); ok {
		r0 = rf(ctx, workspaceID, options)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
