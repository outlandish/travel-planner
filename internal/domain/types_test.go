package domain

import "testing"

func TestHolidayParamsValidate(t *testing.T) {
	tests := []struct {
		name    string
		params  HolidayParams
		wantErr bool
	}{
		{
			name:    "valid",
			params:  HolidayParams{Budget: 1000, HolidayType: HolidayTypeBeach, Nature: HolidayNatureSea},
			wantErr: false,
		},
		{
			name:    "zero budget",
			params:  HolidayParams{Budget: 0, HolidayType: HolidayTypeBeach, Nature: HolidayNatureSea},
			wantErr: true,
		},
		{
			name:    "negative budget",
			params:  HolidayParams{Budget: -5, HolidayType: HolidayTypeBeach, Nature: HolidayNatureSea},
			wantErr: true,
		},
		{
			name:    "unknown holiday type",
			params:  HolidayParams{Budget: 1000, HolidayType: HolidayType("space"), Nature: HolidayNatureSea},
			wantErr: true,
		},
		{
			name:    "unknown nature",
			params:  HolidayParams{Budget: 1000, HolidayType: HolidayTypeBeach, Nature: HolidayNature("volcano")},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
