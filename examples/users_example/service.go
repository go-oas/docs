package users

type service struct {
	repo Repository
}

type Repository interface {
	CreateUser() error
}

func newService(repo Repository) service {
	return service{
		repo,
	}
}
