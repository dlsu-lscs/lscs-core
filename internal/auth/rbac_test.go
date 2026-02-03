package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPositionLevel(t *testing.T) {
	tests := []struct {
		name       string
		positionID string
		want       int
	}{
		{"PRES has highest level", "PRES", 7},
		{"EVP has second highest", "EVP", 6},
		{"VP has level 5", "VP", 5},
		{"AVP has level 4", "AVP", 4},
		{"CT has level 3", "CT", 3},
		{"JO has level 2", "JO", 2},
		{"MEM has level 1", "MEM", 1},
		{"unknown position has level 0", "UNKNOWN", 0},
		{"empty position has level 0", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetPositionLevel(tt.positionID)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsHigherPosition(t *testing.T) {
	tests := []struct {
		name      string
		position1 string
		position2 string
		want      bool
	}{
		{"PRES higher than EVP", "PRES", "EVP", true},
		{"EVP higher than VP", "EVP", "VP", true},
		{"VP higher than AVP", "VP", "AVP", true},
		{"AVP higher than CT", "AVP", "CT", true},
		{"CT higher than JO", "CT", "JO", true},
		{"JO higher than MEM", "JO", "MEM", true},

		{"EVP not higher than PRES", "EVP", "PRES", false},
		{"VP not higher than EVP", "VP", "EVP", false},
		{"MEM not higher than anyone", "MEM", "JO", false},

		{"same position not higher", "VP", "VP", false},
		{"unknown not higher than MEM", "UNKNOWN", "MEM", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsHigherPosition(tt.position1, tt.position2)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsHigherOrEqualPosition(t *testing.T) {
	tests := []struct {
		name      string
		position1 string
		position2 string
		want      bool
	}{
		{"PRES >= EVP", "PRES", "EVP", true},
		{"PRES >= PRES", "PRES", "PRES", true},
		{"VP >= VP", "VP", "VP", true},
		{"VP >= AVP", "VP", "AVP", true},
		{"AVP >= MEM", "AVP", "MEM", true},

		{"EVP not >= PRES", "EVP", "PRES", false},
		{"MEM not >= JO", "MEM", "JO", false},
		{"CT not >= AVP", "CT", "AVP", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsHigherOrEqualPosition(tt.position1, tt.position2)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEditableFields(t *testing.T) {
	t.Run("self-editable fields are defined", func(t *testing.T) {
		expectedFields := []EditableField{
			FieldNickname,
			FieldTelegram,
			FieldDiscord,
			FieldInterests,
			FieldContactNumber,
			FieldFbLink,
		}

		for _, field := range expectedFields {
			assert.True(t, selfEditableFields[field], "field %s should be self-editable", field)
		}
	})

	t.Run("authorized-editable fields are defined", func(t *testing.T) {
		expectedFields := []EditableField{
			FieldFullName,
			FieldEmail,
			FieldPositionID,
			FieldCommitteeID,
			FieldCollege,
			FieldProgram,
			FieldHouseID,
		}

		for _, field := range expectedFields {
			assert.True(t, authorizedEditableFields[field], "field %s should be authorized-editable", field)
		}
	})

	t.Run("self-editable fields should not overlap with authorized-editable", func(t *testing.T) {
		for field := range selfEditableFields {
			assert.False(t, authorizedEditableFields[field], "field %s should not be in both maps", field)
		}
	})
}
