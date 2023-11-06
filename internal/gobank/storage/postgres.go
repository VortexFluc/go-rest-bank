package storage

import (
	"github.com/vortexfluc/gobank/internal/gobank/account"
	"github.com/vortexfluc/gobank/internal/gobank/storage/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresStore struct {
	db *gorm.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: "host=localhost user=postgres password=gobank dbname=postgres port=5432 sslmode=disable TimeZone=Europe/Moscow",
	}), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.db.AutoMigrate(&entity.Account{})
}

func (s *PostgresStore) CreateAccount(acc *account.Account) (int, error) {
	return acc.ID, s.db.Create(acc).Error
}

func (s *PostgresStore) UpdateAccount(acc *account.Account) error {
	e := mapToEntity(acc)
	err := s.db.Save(&e).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	return s.db.Delete(&entity.Account{}, id).Error
}

func (s *PostgresStore) GetAccountByNumber(number int64) (*account.Account, error) {
	var acEntity entity.Account
	err := s.db.First(&acEntity, "number = ?", number).Error
	return mapToModel(&acEntity), err
}

func (s *PostgresStore) GetAccountById(id int) (*account.Account, error) {
	var acEntity entity.Account
	err := s.db.First(&acEntity, id).Error
	return mapToModel(&acEntity), err
}

func (s *PostgresStore) GetAccounts() ([]*account.Account, error) {
	var accounts []*entity.Account
	err := s.db.Find(&accounts).Error
	if err != nil {
		return nil, err
	}

	modelAc := make([]*account.Account, 0)
	for _, e := range accounts {
		modelAc = append(modelAc, mapToModel(e))
	}
	return modelAc, nil
}

func mapToModel(entity *entity.Account) *account.Account {
	return &account.Account{
		ID:                entity.ID,
		FirstName:         entity.FirstName,
		LastName:          entity.LastName,
		Number:            entity.Number,
		EncryptedPassword: entity.EncryptedPassword,
		Balance:           entity.Balance,
		CreatedAt:         entity.CreatedAt,
	}
}

func mapToEntity(model *account.Account) *entity.Account {
	return &entity.Account{
		ID:                model.ID,
		FirstName:         model.FirstName,
		LastName:          model.LastName,
		Number:            model.Number,
		EncryptedPassword: model.EncryptedPassword,
		Balance:           model.Balance,
		Model:             gorm.Model{CreatedAt: model.CreatedAt},
	}
}
