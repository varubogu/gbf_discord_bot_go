package gbf

import (
	"testing"
)

func TestAttackCalculator_CalculateBaseDamage(t *testing.T) {
	calc := NewAttackCalculator()

	tests := []struct {
		name     string
		attack   int
		defense  int
		expected int
	}{
		{
			name:     "basic damage calculation",
			attack:   1000,
			defense:  500,
			expected: 950, // 1000 * (1 - 500/10000) = 950
		},
		{
			name:     "high defense",
			attack:   1000,
			defense:  9000,
			expected: 100, // 1000 * (1 - 9000/10000) = 100
		},
		{
			name:     "zero attack",
			attack:   0,
			defense:  500,
			expected: 0,
		},
		{
			name:     "negative attack",
			attack:   -100,
			defense:  500,
			expected: 0,
		},
		{
			name:     "very high defense results in minimum damage",
			attack:   100,
			defense:  9900,
			expected: 1, // Should be capped at minimum 1 damage
		},
		{
			name:     "no defense",
			attack:   1000,
			defense:  0,
			expected: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.CalculateBaseDamage(tt.attack, tt.defense)
			if result != tt.expected {
				t.Errorf("CalculateBaseDamage(%d, %d) = %d, expected %d",
					tt.attack, tt.defense, result, tt.expected)
			}
		})
	}
}

func TestAttackCalculator_CalculateCriticalDamage(t *testing.T) {
	calc := NewAttackCalculator()

	tests := []struct {
		name           string
		baseDamage     int
		critMultiplier float64
		expected       int
	}{
		{
			name:           "normal critical hit",
			baseDamage:     1000,
			critMultiplier: 1.5,
			expected:       1500,
		},
		{
			name:           "high critical multiplier",
			baseDamage:     1000,
			critMultiplier: 2.0,
			expected:       2000,
		},
		{
			name:           "zero base damage",
			baseDamage:     0,
			critMultiplier: 1.5,
			expected:       0,
		},
		{
			name:           "negative base damage",
			baseDamage:     -100,
			critMultiplier: 1.5,
			expected:       -100,
		},
		{
			name:           "zero multiplier",
			baseDamage:     1000,
			critMultiplier: 0,
			expected:       1000,
		},
		{
			name:           "negative multiplier",
			baseDamage:     1000,
			critMultiplier: -0.5,
			expected:       1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.CalculateCriticalDamage(tt.baseDamage, tt.critMultiplier)
			if result != tt.expected {
				t.Errorf("CalculateCriticalDamage(%d, %f) = %d, expected %d",
					tt.baseDamage, tt.critMultiplier, result, tt.expected)
			}
		})
	}
}

func TestAttackCalculator_CalculateElementalDamage(t *testing.T) {
	calc := NewAttackCalculator()

	tests := []struct {
		name              string
		baseDamage        int
		elementalModifier float64
		expected          int
	}{
		{
			name:              "elemental advantage",
			baseDamage:        1000,
			elementalModifier: 1.5,
			expected:          1500,
		},
		{
			name:              "elemental disadvantage",
			baseDamage:        1000,
			elementalModifier: 0.5,
			expected:          500,
		},
		{
			name:              "neutral element",
			baseDamage:        1000,
			elementalModifier: 1.0,
			expected:          1000,
		},
		{
			name:              "zero base damage",
			baseDamage:        0,
			elementalModifier: 1.5,
			expected:          0,
		},
		{
			name:              "negative base damage",
			baseDamage:        -100,
			elementalModifier: 1.5,
			expected:          0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.CalculateElementalDamage(tt.baseDamage, tt.elementalModifier)
			if result != tt.expected {
				t.Errorf("CalculateElementalDamage(%d, %f) = %d, expected %d",
					tt.baseDamage, tt.elementalModifier, result, tt.expected)
			}
		})
	}
}

func TestIsValidWeaponType(t *testing.T) {
	tests := []struct {
		name       string
		weaponType WeaponType
		expected   bool
	}{
		{
			name:       "valid sword",
			weaponType: WeaponTypeSword,
			expected:   true,
		},
		{
			name:       "valid dagger",
			weaponType: WeaponTypeDagger,
			expected:   true,
		},
		{
			name:       "valid katana",
			weaponType: WeaponTypeKatana,
			expected:   true,
		},
		{
			name:       "invalid weapon type",
			weaponType: WeaponType("invalid"),
			expected:   false,
		},
		{
			name:       "empty weapon type",
			weaponType: WeaponType(""),
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidWeaponType(tt.weaponType)
			if result != tt.expected {
				t.Errorf("IsValidWeaponType(%s) = %t, expected %t",
					tt.weaponType, result, tt.expected)
			}
		})
	}
}

func TestGetWeaponTypeMultiplier(t *testing.T) {
	tests := []struct {
		name        string
		weaponType  WeaponType
		className   string
		expected    float64
		expectError bool
	}{
		{
			name:        "valid sword",
			weaponType:  WeaponTypeSword,
			className:   "fighter",
			expected:    1.2,
			expectError: false,
		},
		{
			name:        "valid axe",
			weaponType:  WeaponTypeAxe,
			className:   "berserker",
			expected:    1.25,
			expectError: false,
		},
		{
			name:        "valid melee",
			weaponType:  WeaponTypeMelee,
			className:   "monk",
			expected:    1.3,
			expectError: false,
		},
		{
			name:        "invalid weapon type",
			weaponType:  WeaponType("invalid"),
			className:   "fighter",
			expected:    0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetWeaponTypeMultiplier(tt.weaponType, tt.className)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("GetWeaponTypeMultiplier(%s, %s) expected error but got none",
						tt.weaponType, tt.className)
				}
			} else {
				if err != nil {
					t.Errorf("GetWeaponTypeMultiplier(%s, %s) unexpected error: %v",
						tt.weaponType, tt.className, err)
				}
				if result != tt.expected {
					t.Errorf("GetWeaponTypeMultiplier(%s, %s) = %f, expected %f",
						tt.weaponType, tt.className, result, tt.expected)
				}
			}
		})
	}
}

// Benchmark tests for performance
func BenchmarkCalculateBaseDamage(b *testing.B) {
	calc := NewAttackCalculator()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.CalculateBaseDamage(1000, 500)
	}
}

func BenchmarkCalculateCriticalDamage(b *testing.B) {
	calc := NewAttackCalculator()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.CalculateCriticalDamage(1000, 1.5)
	}
}

func BenchmarkIsValidWeaponType(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsValidWeaponType(WeaponTypeSword)
	}
}