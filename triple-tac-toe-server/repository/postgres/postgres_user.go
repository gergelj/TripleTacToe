package postgres

import (
	"errors"
	"gregvader/triple-tac-toe/database"
	"log"

	"gregvader/triple-tac-toe/domain"

	"gregvader/triple-tac-toe/security/helpers/hash"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type postgresUserRepository struct {
	conn *database.DBConn
}

func (p postgresUserRepository) WithTrx(trxHandle *gorm.DB) domain.UserRepository {
	if trxHandle == nil {
		log.Print("Transaction Database not found")
		return p
	}
	p.conn.DB = trxHandle
	return p
}

func (p postgresUserRepository) FindByUsername(username string) (domain.User, error) {
	var entity domain.User
	result := p.conn.DB.Where(&domain.User{Username: username}).First(&entity)
	return entity, result.Error
}

func (p postgresUserRepository) Create(user domain.User) (domain.User, error) {
	_, valErr := govalidator.ValidateStruct(user)

	if valErr != nil {
		return domain.User{}, errors.New("create user: user validation failed")
	}

	hashedPw, hashErr := hash.HashPassword(user.Password)
	if hashErr != nil {
		return user, errors.New("create user: password hashing error")
	}

	user.Password = hashedPw

	result := p.conn.DB.Create(&user)
	if err := result.Error; err != nil {
		return domain.User{}, errors.New("create user: username is not unique")
	}
	return user, nil
}

func (p postgresUserRepository) DeleteByUsername(Username string) (domain.User, error) {
	var user domain.User
	result := p.conn.DB.Where(&domain.User{Username: Username}).First(&user)

	if result.Error != nil {
		return domain.User{}, errors.New("delete user: user does not exist")
	}
	p.conn.DB.Delete(user)
	return user, nil
}

/*
func (p postgresUserRepository) GetAll() []domain.User {
	var Users []domain.User
	p.conn.DB.Find(&Users)
	return Users
}

func (p postgresUserRepository) Search(username string, page int, pageSize int) []domain.User {
	var Users []domain.User
	//p.conn.DB.Scopes(helpers.Paginate(page, pageSize)).Joins("JOIN registered_users ON users.id = registered_users.user_id").Joins("JOIN profile_configurations ON registered_users.prof_conf_id = profile_configurations.id").Where("Public_Profile = ?", true).Where("upper(Username) LIKE ?", strings.ToUpper("%"+username+"%")).Find(&Users)
	p.conn.DB.Scopes(helpers.Paginate(page, pageSize)).Joins("JOIN registered_users ON users.id = registered_users.user_id").Where("upper(Username) LIKE ?", strings.ToUpper("%"+username+"%")).Find(&Users)
	return Users
}
*/

func NewPostgresUserRepository(conn *database.DBConn) domain.UserRepository {
	return &postgresUserRepository{conn: conn}
}
