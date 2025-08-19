package gbf

import (
	"fmt"
	"strings"
	"time"
)

// BattleType represents different types of battles in GBF
type BattleType string

const (
	// Raid battles
	BattleTypeHL       BattleType = "hl"        // High Level raids
	BattleTypeFaaHL    BattleType = "faa_hl"    // Lucilius Hard
	BattleTypeBahaHL   BattleType = "baha_hl"   // Bahamut Hard
	BattleTypeUBahaHL  BattleType = "ubaha_hl"  // Ultimate Bahamut Hard
	BattleTypeAkashaHL BattleType = "akasha_hl" // Akasha Hard
	BattleTypeLuciHL   BattleType = "luci_hl"   // Lucifer Hard
	
	// Guild War
	BattleTypeGW       BattleType = "gw"        // Guild War
	BattleTypeGWNM     BattleType = "gw_nm"     // Guild War Nightmare
	
	// Events
	BattleTypeEvent    BattleType = "event"     // Event raids
	BattleTypeEventHL  BattleType = "event_hl"  // Event High Level
	
	// Others
	BattleTypeTrain    BattleType = "train"     // Training rooms
	BattleTypeCustom   BattleType = "custom"    // Custom battles
)

// BattleInfo represents information about a specific battle
type BattleInfo struct {
	ID          string
	Name        string
	Type        BattleType
	Level       int
	MinRank     int
	MaxPlayers  int
	Description string
	IsActive    bool
	CreatedAt   time.Time
}

// BattleManager manages battle types and information
type BattleManager struct {
	battles map[string]*BattleInfo
}

// NewBattleManager creates a new battle manager
func NewBattleManager() *BattleManager {
	bm := &BattleManager{
		battles: make(map[string]*BattleInfo),
	}
	
	// Initialize default battles
	bm.initializeDefaultBattles()
	
	return bm
}

// initializeDefaultBattles sets up default battle configurations
func (bm *BattleManager) initializeDefaultBattles() {
	defaultBattles := []*BattleInfo{
		{
			ID:          "faa_hl",
			Name:        "Lucilius (Hard)",
			Type:        BattleTypeFaaHL,
			Level:       200,
			MinRank:     150,
			MaxPlayers:  6,
			Description: "Dark Rapture Hard mode raid",
			IsActive:    true,
		},
		{
			ID:          "baha_hl",
			Name:        "Proto Bahamut (Hard)",
			Type:        BattleTypeBahaHL,
			Level:       150,
			MinRank:     101,
			MaxPlayers:  18,
			Description: "Proto Bahamut Hard mode raid",
			IsActive:    true,
		},
		{
			ID:          "ubaha_hl",
			Name:        "Ultimate Bahamut (Hard)",
			Type:        BattleTypeUBahaHL,
			Level:       200,
			MinRank:     120,
			MaxPlayers:  30,
			Description: "Ultimate Bahamut Hard mode raid",
			IsActive:    true,
		},
		{
			ID:          "akasha_hl",
			Name:        "Akasha (Hard)",
			Type:        BattleTypeAkashaHL,
			Level:       200,
			MinRank:     120,
			MaxPlayers:  30,
			Description: "Akasha Hard mode raid",
			IsActive:    true,
		},
		{
			ID:          "luci_hl",
			Name:        "Lucifer (Hard)",
			Type:        BattleTypeEventHL,
			Level:       200,
			MinRank:     120,
			MaxPlayers:  30,
			Description: "Lucifer Hard mode raid",
			IsActive:    true,
		},
		{
			ID:          "gw_nm95",
			Name:        "Guild War NM95",
			Type:        BattleTypeGWNM,
			Level:       95,
			MinRank:     80,
			MaxPlayers:  30,
			Description: "Guild War Nightmare 95 raid",
			IsActive:    false, // Activated during GW periods
		},
		{
			ID:          "gw_nm150",
			Name:        "Guild War NM150",
			Type:        BattleTypeGWNM,
			Level:       150,
			MinRank:     120,
			MaxPlayers:  30,
			Description: "Guild War Nightmare 150 raid",
			IsActive:    false, // Activated during GW periods
		},
	}
	
	for _, battle := range defaultBattles {
		battle.CreatedAt = time.Now()
		bm.battles[battle.ID] = battle
	}
}

// GetBattle retrieves battle information by ID
func (bm *BattleManager) GetBattle(id string) (*BattleInfo, error) {
	battle, exists := bm.battles[strings.ToLower(id)]
	if !exists {
		return nil, fmt.Errorf("battle not found: %s", id)
	}
	return battle, nil
}

// GetActiveBattles returns all currently active battles
func (bm *BattleManager) GetActiveBattles() []*BattleInfo {
	var activeBattles []*BattleInfo
	for _, battle := range bm.battles {
		if battle.IsActive {
			activeBattles = append(activeBattles, battle)
		}
	}
	return activeBattles
}

// GetBattlesByType returns battles of a specific type
func (bm *BattleManager) GetBattlesByType(battleType BattleType) []*BattleInfo {
	var typeBattles []*BattleInfo
	for _, battle := range bm.battles {
		if battle.Type == battleType {
			typeBattles = append(typeBattles, battle)
		}
	}
	return typeBattles
}

// AddBattle adds a new battle type
func (bm *BattleManager) AddBattle(battle *BattleInfo) error {
	if battle.ID == "" {
		return fmt.Errorf("battle ID cannot be empty")
	}
	
	battle.CreatedAt = time.Now()
	bm.battles[strings.ToLower(battle.ID)] = battle
	return nil
}

// UpdateBattle updates an existing battle
func (bm *BattleManager) UpdateBattle(id string, battle *BattleInfo) error {
	if _, exists := bm.battles[strings.ToLower(id)]; !exists {
		return fmt.Errorf("battle not found: %s", id)
	}
	
	battle.ID = strings.ToLower(id)
	bm.battles[battle.ID] = battle
	return nil
}

// RemoveBattle removes a battle type
func (bm *BattleManager) RemoveBattle(id string) error {
	if _, exists := bm.battles[strings.ToLower(id)]; !exists {
		return fmt.Errorf("battle not found: %s", id)
	}
	
	delete(bm.battles, strings.ToLower(id))
	return nil
}

// SetBattleActive sets the active status of a battle
func (bm *BattleManager) SetBattleActive(id string, active bool) error {
	battle, exists := bm.battles[strings.ToLower(id)]
	if !exists {
		return fmt.Errorf("battle not found: %s", id)
	}
	
	battle.IsActive = active
	return nil
}

// IsValidBattleType checks if a battle type is valid
func IsValidBattleType(battleType BattleType) bool {
	validTypes := []BattleType{
		BattleTypeHL, BattleTypeFaaHL, BattleTypeBahaHL, BattleTypeUBahaHL,
		BattleTypeAkashaHL, BattleTypeLuciHL, BattleTypeGW, BattleTypeGWNM,
		BattleTypeEvent, BattleTypeEventHL, BattleTypeTrain, BattleTypeCustom,
	}
	
	for _, validType := range validTypes {
		if battleType == validType {
			return true
		}
	}
	return false
}

// FormatBattleInfo formats battle information for display
func FormatBattleInfo(battle *BattleInfo) string {
	status := "🔴 Inactive"
	if battle.IsActive {
		status = "🟢 Active"
	}
	
	return fmt.Sprintf("**%s** (Level %d)\n"+
		"Type: %s | Min Rank: %d | Max Players: %d\n"+
		"Status: %s\n"+
		"Description: %s",
		battle.Name, battle.Level, battle.Type, battle.MinRank,
		battle.MaxPlayers, status, battle.Description)
}