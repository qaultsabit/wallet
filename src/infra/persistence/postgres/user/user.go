package user

import (
	"fmt"

	dto "github.com/qaultsabit/wallet/src/app/dto/user"

	"log"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	"errors"

	helper "github.com/qaultsabit/wallet/src/infra/helper"
)

type UserRepository interface {
	Register(data *dto.RegisterReqDTO) (*dto.RegisterRespDTO, error)
	Login(data *dto.LoginReqDTO) (*dto.RegisterRespDTO, error)
}

const (
	Register = `insert into public.users (username, password)
	values ($1, $2) returning id, username`

	Login = `select u.id, u.username, u.password, w.id as wallet_id
	from public.users u
	join public.wallets w
	on u.id = w.user_id
	where u.username = $1`

	CreateWallet = `insert into public.wallets (user.id)
	values ($1) returning id as wallet_id`
)

var statement PreparedStatement

type PreparedStatement struct {
	login *sqlx.Stmt
}

type userRepo struct {
	Connection *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	repo := &userRepo{
		Connection: db,
	}
	InitPreparedStatement(repo)
	return repo
}

func (p *userRepo) Preparex(query string) *sqlx.Stmt {
	statement, err := p.Connection.Preparex(query)
	if err != nil {
		log.Fatalf("Failed to preparex query: %s. Error: %s", query, err.Error())
	}

	return statement
}

func InitPreparedStatement(m *userRepo) {
	statement = PreparedStatement{
		login: m.Preparex(Login),
	}
}

func (p *userRepo) Register(data *dto.RegisterReqDTO) (resp *dto.RegisterRespDTO, err error) {
	var resultData dto.RegisterModel

	pwd, err := hashPassword(data.Password)
	if err != nil {
		return nil, err
	}

	tx, err := p.Connection.Beginx()
	if err != nil {
		log.Println("Failed Begin Tx Register  : ", err.Error())
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			log.Println("Recovered in Register:", p)
			err = fmt.Errorf("panic occurred: %v", p)
		} else if err != nil {
			tx.Rollback()
			log.Println("Rolling back transaction due to:", err)
		} else {
			err = tx.Commit()
			if err != nil {
				log.Println("Failed to commit transaction:", err.Error())
			}
		}
	}()

	if err = tx.QueryRow(Register, data.UserName, pwd).Scan(&resultData.ID, &resultData.UserName); err != nil {
		log.Println("Failed query register: ", err.Error())
		return nil, err
	}

	if err = tx.QueryRow(CreateWallet, resultData.ID).Scan(&resultData.WalletID); err != nil {
		log.Println("Failed query create wallet: ", err.Error())
		return nil, err
	}

	resp = &dto.RegisterRespDTO{}
	if resp.Token, err = helper.GenerateToken(&resultData); err != nil {
		return nil, err
	}

	return resp, nil
}

func (p *userRepo) Login(data *dto.LoginReqDTO) (*dto.RegisterRespDTO, error) {
	var resultData []*dto.RegisterModel
	var resp dto.RegisterRespDTO

	// Execute the login query
	if err := statement.login.Select(&resultData, data.UserName); err != nil {
		return nil, err
	}

	// Check if no rows were returned from the query
	if len(resultData) < 1 {
		return nil, errors.New("no rows returned from the query")
	}

	// Verify the password
	if err := verifyPassword(resultData[0].Password, data.Password); err != nil {
		return nil, err
	}

	// Generate token
	token, err := helper.GenerateToken(resultData[0])
	if err != nil {
		return nil, err
	}
	resp.Token = token

	// Return the response object if everything is successful
	return &resp, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func verifyPassword(hashedPassword, inputPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
}
