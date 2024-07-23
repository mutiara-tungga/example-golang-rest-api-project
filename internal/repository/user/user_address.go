package user

type IUserAddressRepo interface{}

type UserAddressRepo struct{}

func NewUserAddress() *UserAddressRepo {
	return &UserAddressRepo{}
}
