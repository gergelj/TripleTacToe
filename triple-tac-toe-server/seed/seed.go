package seed

import (
	"log"

	"gregvader/triple-tac-toe/database"
	"gregvader/triple-tac-toe/domain"
	"gregvader/triple-tac-toe/repository/postgres"
)

type UserSeed struct {
	Name string
	Run  func(UserRepo domain.UserRepository) error
}

func allUsers() []UserSeed {
	return []UserSeed{
		{
			Name: "Radovan",
			Run: func(RegisteredUserRepository domain.UserRepository) error {
				_, error := RegisteredUserRepository.Create(domain.User{
					Username: "radovan",
					Password: "radovan1",
				})
				return error
			},
		},
		{
			Name: "Gergo",
			Run: func(RegisteredUserRepository domain.UserRepository) error {
				_, error := RegisteredUserRepository.Create(domain.User{
					Username: "gods_best_daemon",
					Password: "00000000",
				})
				return error
			},
		},
		{
			Name: "Dragic",
			Run: func(RegisteredUserRepository domain.UserRepository) error {
				_, error := RegisteredUserRepository.Create(domain.User{
					Username: "nikoladragicc",
					Password: "00000000",
				})
				return error
			},
		},
		{
			Name: "Hadzi",
			Run: func(RegisteredUserRepository domain.UserRepository) error {
				_, error := RegisteredUserRepository.Create(domain.User{
					Username: "hadzi",
					Password: "hadzi1111",
				})
				return error
			},
		},
	}
}

func Run(dbConn *database.DBConn) {
	userRepo := postgres.NewPostgresUserRepository(dbConn)

	for _, seed := range allUsers() {
		if err := seed.Run(userRepo); err != nil {
			log.Fatalf("Running seed '%s', failed with error: %s", seed.Name, err)
		}
	}
}
