package prompt

import (
	"fmt"
	"strings"
	"travelplanner/internal/domain"
)

func Build(params domain.HolidayParams) (string, error) {
	holidayType, err := getHolidayTypeForPrompt(params.HolidayType)
	if err != nil {
		return "", err
	}

	holidayNature, err := getNatureForPrompt(params.Nature)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(
		fmt.Sprintf(`
You are a travel planning assistant.
Suggest 3-5 countries for a holiday based on these parameters:
- Budget: about %d (in US dollars, per person; flights can be ignored)
- Holiday type: %s
- Preferred nature / environment: %s
Response requirements:
1. Only countries that truly match the budget and preferences.
2. For each country: the name and 1-2 sentences on why it fits these parameters.
3. Do not suggest options that clearly contradict the holiday type or nature.
4. Answer as a list, with no introduction or conclusion.
5. Format of each item:
- Country: <name>
Why it fits: <short explanation>`, params.Budget, holidayType, holidayNature)), nil
}

func getNatureForPrompt(nature domain.HolidayNature) (string, error) {
	switch nature {
	case domain.HolidayNatureCity:
		return "city", nil
	case domain.HolidayNatureSea:
		return "sea", nil
	case domain.HolidayNatureMountains:
		return "mountains", nil
	default:
		return "", fmt.Errorf("unknown holiday nature: %q", nature)
	}
}

func getHolidayTypeForPrompt(holidayType domain.HolidayType) (string, error) {
	switch holidayType {
	case domain.HolidayTypeAttractions:
		return "sightseeing", nil
	case domain.HolidayTypeBeach:
		return "beach holiday", nil
	case domain.HolidayTypeActive:
		return "active holiday", nil
	default:
		return "", fmt.Errorf("unknown holiday type: %q", holidayType)
	}
}
