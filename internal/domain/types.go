package domain

import (
	"errors"
	"fmt"
)

type HolidayParams struct {
	Budget      int
	HolidayType HolidayType
	Nature      HolidayNature
}

func (params HolidayParams) Validate() error {
	if err := validateBudget(params.Budget); err != nil {
		return err
	}
	if err := validateHolidayType(params.HolidayType); err != nil {
		return err
	}
	if err := validateHolidayNature(params.Nature); err != nil {
		return err
	}
	return nil
}

func validateBudget(budget int) error {
	if budget <= 0 {
		if budget == 0 {
			return errors.New("budget is required")
		}
		return errors.New("budget should be more than 0")
	}
	return nil
}

func validateHolidayType(holidayType HolidayType) error {
	switch holidayType {
	case "":
		return errors.New("holiday type is required")
	case HolidayTypeActive, HolidayTypeAttractions, HolidayTypeBeach:
		return nil
	default:
		return fmt.Errorf("holiday type is not supported: %q", holidayType)
	}
}

func validateHolidayNature(holidayNature HolidayNature) error {
	switch holidayNature {
	case "":
		return errors.New("holiday nature is required")
	case HolidayNatureMountains, HolidayNatureCity, HolidayNatureSea:
		return nil
	default:
		return fmt.Errorf("holiday nature is not supported: %q", holidayNature)
	}
}

type HolidayType string

const (
	HolidayTypeActive      HolidayType = "active"
	HolidayTypeBeach       HolidayType = "beach"
	HolidayTypeAttractions HolidayType = "attractions"
)

type HolidayNature string

const (
	HolidayNatureSea       HolidayNature = "sea"
	HolidayNatureMountains HolidayNature = "mountains"
	HolidayNatureCity      HolidayNature = "city"
)
