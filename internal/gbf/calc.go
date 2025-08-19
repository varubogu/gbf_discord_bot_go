package gbf

import (
	"fmt"
	"math"
)

// AttackCalculator provides GBF-related damage calculation functions
type AttackCalculator struct{}

// NewAttackCalculator creates a new attack calculator instance
func NewAttackCalculator() *AttackCalculator {
	return &AttackCalculator{}
}

// CalculateBaseDamage calculates basic attack damage
// This is a simplified example for testing purposes
func (c *AttackCalculator) CalculateBaseDamage(attack int, defense int) int {
	if attack <= 0 {
		return 0
	}
	
	// Simple damage calculation with defense reduction
	damage := float64(attack) * (1.0 - float64(defense)/10000.0)
	
	// Ensure minimum damage of 1
	if damage < 1 {
		damage = 1
	}
	
	return int(math.Round(damage))
}

// CalculateCriticalDamage calculates critical hit damage
func (c *AttackCalculator) CalculateCriticalDamage(baseDamage int, critMultiplier float64) int {
	if baseDamage <= 0 || critMultiplier <= 0 {
		return baseDamage
	}
	
	return int(math.Round(float64(baseDamage) * critMultiplier))
}

// CalculateElementalDamage applies elemental advantage/disadvantage
func (c *AttackCalculator) CalculateElementalDamage(baseDamage int, elementalModifier float64) int {
	if baseDamage <= 0 {
		return 0
	}
	
	return int(math.Round(float64(baseDamage) * elementalModifier))
}

// WeaponType represents different weapon types in GBF
type WeaponType string

const (
	WeaponTypeSword   WeaponType = "sword"
	WeaponTypeDagger  WeaponType = "dagger"
	WeaponTypeSpear   WeaponType = "spear"
	WeaponTypeAxe     WeaponType = "axe"
	WeaponTypeStaff   WeaponType = "staff"
	WeaponTypeGun     WeaponType = "gun"
	WeaponTypeMelee   WeaponType = "melee"
	WeaponTypeBow     WeaponType = "bow"
	WeaponTypeHarp    WeaponType = "harp"
	WeaponTypeKatana  WeaponType = "katana"
)

// IsValidWeaponType checks if the given weapon type is valid
func IsValidWeaponType(weaponType WeaponType) bool {
	validTypes := []WeaponType{
		WeaponTypeSword, WeaponTypeDagger, WeaponTypeSpear, WeaponTypeAxe,
		WeaponTypeStaff, WeaponTypeGun, WeaponTypeMelee, WeaponTypeBow,
		WeaponTypeHarp, WeaponTypeKatana,
	}
	
	for _, validType := range validTypes {
		if weaponType == validType {
			return true
		}
	}
	return false
}

// GetWeaponTypeMultiplier returns the weapon type multiplier for specific classes
// This is a simplified example
func GetWeaponTypeMultiplier(weaponType WeaponType, className string) (float64, error) {
	if !IsValidWeaponType(weaponType) {
		return 0, fmt.Errorf("invalid weapon type: %s", weaponType)
	}
	
	// Simplified multiplier system - in real GBF this would be much more complex
	baseMultipliers := map[WeaponType]float64{
		WeaponTypeSword:  1.2,
		WeaponTypeDagger: 1.1,
		WeaponTypeSpear:  1.15,
		WeaponTypeAxe:    1.25,
		WeaponTypeStaff:  1.0,
		WeaponTypeGun:    1.1,
		WeaponTypeMelee:  1.3,
		WeaponTypeBow:    1.05,
		WeaponTypeHarp:   1.0,
		WeaponTypeKatana: 1.2,
	}
	
	multiplier, exists := baseMultipliers[weaponType]
	if !exists {
		return 1.0, nil
	}
	
	return multiplier, nil
}