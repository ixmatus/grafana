package usertest

import (
	"context"

	"github.com/grafana/grafana/pkg/services/user"
)

type FakeUserService struct {
	ExpectedUser             *user.User
	ExpectedSignedInUser     *user.SignedInUser
	ExpectedError            error
	ExpectedSetUsingOrgError error
	ExpectedSearchUsers      user.SearchUserQueryResult
	ExpectedUserProfileDTO   *user.UserProfileDTO
	ExpectedUserProfileDTOs  []*user.UserProfileDTO
	ExpectedUsageStats       map[string]any

	UpdateFn            func(ctx context.Context, cmd *user.UpdateUserCommand) error
	GetSignedInUserFn   func(ctx context.Context, query *user.GetSignedInUserQuery) (*user.SignedInUser, error)
	CreateFn            func(ctx context.Context, cmd *user.CreateUserCommand) (*user.User, error)
	BatchDisableUsersFn func(ctx context.Context, cmd *user.BatchDisableUsersCommand) error

	counter int
}

func NewUserServiceFake() *FakeUserService {
	return &FakeUserService{}
}

func (f FakeUserService) GetUsageStats(ctx context.Context) map[string]any {
	return f.ExpectedUsageStats
}

func (f *FakeUserService) Create(ctx context.Context, cmd *user.CreateUserCommand) (*user.User, error) {
	if f.CreateFn != nil {
		return f.CreateFn(ctx, cmd)
	}

	return f.ExpectedUser, f.ExpectedError
}

func (f *FakeUserService) CreateServiceAccount(ctx context.Context, cmd *user.CreateUserCommand) (*user.User, error) {
	return f.ExpectedUser, f.ExpectedError
}

func (f *FakeUserService) Delete(ctx context.Context, cmd *user.DeleteUserCommand) error {
	return f.ExpectedError
}

func (f *FakeUserService) GetByID(ctx context.Context, query *user.GetUserByIDQuery) (*user.User, error) {
	return f.ExpectedUser, f.ExpectedError
}

func (f *FakeUserService) GetByLogin(ctx context.Context, query *user.GetUserByLoginQuery) (*user.User, error) {
	return f.ExpectedUser, f.ExpectedError
}

func (f *FakeUserService) GetByEmail(ctx context.Context, query *user.GetUserByEmailQuery) (*user.User, error) {
	return f.ExpectedUser, f.ExpectedError
}

func (f *FakeUserService) Update(ctx context.Context, cmd *user.UpdateUserCommand) error {
	if f.UpdateFn != nil {
		return f.UpdateFn(ctx, cmd)
	}
	return f.ExpectedError
}

func (f *FakeUserService) ChangePassword(ctx context.Context, cmd *user.ChangeUserPasswordCommand) error {
	return f.ExpectedError
}

func (f *FakeUserService) UpdateLastSeenAt(ctx context.Context, cmd *user.UpdateUserLastSeenAtCommand) error {
	return f.ExpectedError
}

func (f *FakeUserService) SetUsingOrg(ctx context.Context, cmd *user.SetUsingOrgCommand) error {
	return f.ExpectedSetUsingOrgError
}

func (f *FakeUserService) GetSignedInUserWithCacheCtx(ctx context.Context, query *user.GetSignedInUserQuery) (*user.SignedInUser, error) {
	return f.GetSignedInUser(ctx, query)
}

func (f *FakeUserService) GetSignedInUser(ctx context.Context, query *user.GetSignedInUserQuery) (*user.SignedInUser, error) {
	if f.GetSignedInUserFn != nil {
		return f.GetSignedInUserFn(ctx, query)
	}
	if f.ExpectedSignedInUser == nil {
		return &user.SignedInUser{}, f.ExpectedError
	}
	return f.ExpectedSignedInUser, f.ExpectedError
}

func (f *FakeUserService) NewAnonymousSignedInUser(ctx context.Context) (*user.SignedInUser, error) {
	return f.ExpectedSignedInUser, f.ExpectedError
}

func (f *FakeUserService) Search(ctx context.Context, query *user.SearchUsersQuery) (*user.SearchUserQueryResult, error) {
	return &f.ExpectedSearchUsers, f.ExpectedError
}

func (f *FakeUserService) BatchDisableUsers(ctx context.Context, cmd *user.BatchDisableUsersCommand) error {
	if f.BatchDisableUsersFn != nil {
		return f.BatchDisableUsersFn(ctx, cmd)
	}
	return f.ExpectedError
}

func (f *FakeUserService) SetUserHelpFlag(ctx context.Context, cmd *user.SetUserHelpFlagCommand) error {
	return f.ExpectedError
}

func (f *FakeUserService) GetProfile(ctx context.Context, query *user.GetUserProfileQuery) (*user.UserProfileDTO, error) {
	if f.ExpectedUserProfileDTO != nil {
		return f.ExpectedUserProfileDTO, f.ExpectedError
	}

	if f.ExpectedUserProfileDTOs == nil {
		return nil, f.ExpectedError
	}

	f.counter++
	return f.ExpectedUserProfileDTOs[f.counter-1], f.ExpectedError
}
