package egov

import (
	"context"
	"fmt"
	"strconv"
	"time"

	errs "github.com/IceMAN2377/kaspitest/internal/errors"
	"github.com/IceMAN2377/kaspitest/internal/models"
	"github.com/IceMAN2377/kaspitest/internal/repository"
	"github.com/IceMAN2377/kaspitest/internal/service"
)

type egov struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) service.Service {
	return &egov{
		repo: repo,
	}
}

func (e *egov) CheckIIN(ctx context.Context, iin string) (*models.IINInfo, error) {
	if len(iin) != 12 {
		return nil, fmt.Errorf("IIN must be of length 12")
	}

	datePart := iin[0:6]
	centuryDigit := iin[6:7]
	checkSumDigitStr := iin[11:12]

	dayStr := datePart[4:6]
	monthStr := datePart[2:4]
	yearShortStr := datePart[0:2]

	day, err := strconv.Atoi(dayStr)
	if err != nil {
		return nil, fmt.Errorf("error converting day of birth: %w", err)
	}
	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return nil, fmt.Errorf("error converting month of birth: %w", err)
	}
	yearShort, err := strconv.Atoi(yearShortStr)
	if err != nil {
		return nil, fmt.Errorf("error converting year of birth: %w", err)
	}

	var fullYear int
	century, err := strconv.Atoi(centuryDigit)
	if err != nil {
		return nil, fmt.Errorf("error converting century: %w", err)
	}

	var sex string
	switch century {
	case 1, 3, 5:
		sex = "male"
	case 2, 4, 6:
		sex = "female"
	default:
		return nil, fmt.Errorf("incorrect century value: %s", centuryDigit)
	}

	switch century {
	case 1, 2:
		fullYear = 1800 + yearShort
	case 3, 4:
		fullYear = 1900 + yearShort
	case 5, 6:
		fullYear = 2000 + yearShort
	default:
		return nil, fmt.Errorf("incorrect century value: %s", centuryDigit)
	}

	// Проверка корректности даты и форматирование
	dateOfBirthTime := time.Date(fullYear, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	if dateOfBirthTime.Year() != fullYear || int(dateOfBirthTime.Month()) != month || dateOfBirthTime.Day() != day {
		return nil, fmt.Errorf("incorrect date of birth")
	}
	dateOfBirth := dateOfBirthTime.Format("02.01.2006")

	// Проверка контрольной суммы
	checkSumDigit, err := strconv.Atoi(checkSumDigitStr)
	if err != nil {
		return nil, fmt.Errorf("error converting control value: %w", err)
	}

	weights1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	sum1 := 0
	for i := 0; i < 11; i++ {
		digit, err := strconv.Atoi(string(iin[i]))
		if err != nil {
			return nil, fmt.Errorf("error converting IIN for control sum: %w", err)
		}
		sum1 += digit * weights1[i]
	}
	remainder1 := sum1 % 11

	if remainder1 == checkSumDigit {
		return &models.IINInfo{
			Correct:   true,
			Gender:    sex,
			Birthdate: dateOfBirth,
		}, nil
	}

	if remainder1 == 10 {
		weights2 := []int{3, 4, 5, 6, 7, 8, 9, 10, 11, 1, 2}
		sum2 := 0
		for i := 0; i < 11; i++ {
			digit, err := strconv.Atoi(string(iin[i]))
			if err != nil {
				return nil, fmt.Errorf("error converting IIN for control sum (second attempt): %w", err)
			}
			sum2 += digit * weights2[i]
		}
		remainder2 := sum2 % 11
		if remainder2 == checkSumDigit {
			return &models.IINInfo{
				Correct:   true,
				Gender:    sex,
				Birthdate: dateOfBirth,
			}, nil
		}
	}

	return nil, fmt.Errorf("incorrect control sum of IIN")
}

func (e *egov) GetByIIN(ctx context.Context, iin string) (*models.User, error) {
	return e.repo.GetByIIN(ctx, iin)
}

func (e *egov) GetBySearch(ctx context.Context, search string) ([]models.User, error) {
	return e.repo.GetBySearch(ctx, search)
}

func (e *egov) CreatePerson(ctx context.Context, person *models.User) error {
	if _, err := e.CheckIIN(ctx, person.IIN); err != nil {
		return errs.ErrIncorrectData
	}

	return e.repo.CreatePerson(ctx, person)
}
