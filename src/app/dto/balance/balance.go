package balance

import validation "github.com/go-ozzo/ozzo-validation"

type BalanceDTOInterface interface {
	Validate() error
}

type TopUpReqDTO struct {
	Amount   float64 `json:"amount"`
	WalletID int64   `json:"wallet_id"`
}

func (dto *TopUpReqDTO) Validate() error {
	if err := validation.ValidateStruct(
		dto,
		validation.Field(
			&dto.Amount,
			validation.Required,
			validation.Min(10000.00),
			validation.Max(9999999.99),
		),
	); err != nil {
		return err
	}
	return nil
}

type BalanceReqDTO struct {
	WalletID int64 `json:"wallet_id"`
}

type BalanceRespDTO struct {
	Balance float64 `json:"balance" dto:"balance"`
}
