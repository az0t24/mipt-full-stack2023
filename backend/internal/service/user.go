package service

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slices"
	"scoreboardpro/internal/entity"
	"scoreboardpro/internal/entity/model"
	"scoreboardpro/internal/entity/model/football"
)

type UserService struct {
	userRepo  entity.UserRepository
	tableRepo entity.TableRepository
}

func NewUserService(userRepo entity.UserRepository, tableRepo entity.TableRepository) *UserService {
	return &UserService{
		userRepo:  userRepo,
		tableRepo: tableRepo,
	}
}

func (s *UserService) GetAll() (*[]model.User, error) {
	usersDB, err := s.userRepo.GetAll()
	if err != nil {
		return nil, err
	} else {
		for i := 0; i < len(*usersDB); i++ {
			(*usersDB)[i].OmitPassword()
		}
		return usersDB, nil
	}
}

func (s *UserService) Get(id uint) (*model.User, error) {
	userDB, err := s.userRepo.Get(id)
	if err != nil {
		return nil, err
	} else {
		userDB.OmitPassword()
		return userDB, nil
	}
}

func (s *UserService) GetByEmail(email string) (*model.User, error) {
	userDB, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	} else {
		userDB.OmitPassword()
		return userDB, nil
	}
}

func (s *UserService) Update(user *model.User) error {
	userDB, err := s.userRepo.Get(user.ID)
	if err != nil {
		return err
	}

	var newUserPswdHash string
	if comparePaaswordWithHash(user.Password, userDB.Password) != nil {
		newUserPswdHash, err = generatePasswordHash(user.Password)
		if err != nil {
			return err
		}
	}

	user.Password = newUserPswdHash
	err = s.userRepo.Update(user)
	return err
}

func comparePaaswordWithHash(pswdFromInput, pswdHashFromDB string) error {
	err := bcrypt.CompareHashAndPassword([]byte(pswdHashFromDB), []byte(pswdFromInput))
	return err
}

func generatePasswordHash(pswd string) (string, error) {
	pswdHash, err := bcrypt.GenerateFromPassword([]byte(pswd), bcrypt.DefaultCost)
	return string(pswdHash), err
}

func (s *UserService) Delete(user *model.User) error {
	err := s.userRepo.Delete(user.ID)
	return err
}

func (s *UserService) Register(userReg *model.UserRegister) error {
	if userReg.Email == "" {
		return entity.ErrInvalidEmail
	}

	if len(userReg.Password) < 8 {
		return entity.ErrInvalidPassword
	}

	pswdHash, err := generatePasswordHash(userReg.Password)
	if err != nil {
		return err
	}

	userReg.Password = pswdHash
	user := model.User{UserRegister: *userReg}
	err = s.userRepo.Create(&user)
	return err
}

func (s *UserService) Login(userLogin *model.UserLogin) (uint, error) {
	userDB, err := s.userRepo.GetByEmail(userLogin.Email)
	if err != nil {
		return 0, err
	}

	err = comparePaaswordWithHash(userLogin.Password, userDB.Password)
	if err != nil {
		return 0, err
	}

	return userDB.ID, nil
}

func (s *UserService) MarkAsFavorite(id, tableId uint) error {
	user, err := s.userRepo.Get(id)
	if err != nil {
		return err
	}

	for i, t := range user.FavoriteTables {
		if t.ID == tableId {
			slices.Delete(user.FavoriteTables, i, i)
			break
		} else if i == len(user.FavoriteTables)-1 {
			table, err := s.tableRepo.Get(tableId)
			if err == nil {
				user.FavoriteTables = append(user.FavoriteTables, *table)
			} else {
				return errors.New("table does not exist")
			}
		}
	}

	return nil
}

func (s *UserService) GetFavorites(id uint) ([]football.SeasonTable, error) {
	user, err := s.userRepo.Get(id)
	if err != nil {
		return []football.SeasonTable{}, err
	}

	return user.FavoriteTables, nil
}
