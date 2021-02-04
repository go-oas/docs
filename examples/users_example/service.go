package users

type service struct {
	repo exampleRepo
}

type exampleRepo interface {
	CreateUser() error
}

func newService(repo exampleRepo) service {
	return service{
		repo,
	}
}
