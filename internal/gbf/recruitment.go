package gbf

import (
	"fmt"
	"time"
)

// RecruitmentStatus represents the status of a recruitment
type RecruitmentStatus string

const (
	RecruitmentStatusOpen      RecruitmentStatus = "open"      // Accepting participants
	RecruitmentStatusFull      RecruitmentStatus = "full"      // Full, no more participants
	RecruitmentStatusClosed    RecruitmentStatus = "closed"    // Manually closed
	RecruitmentStatusCompleted RecruitmentStatus = "completed" // Battle completed
	RecruitmentStatusCancelled RecruitmentStatus = "cancelled" // Cancelled
)

// ParticipantRole represents the role of a participant
type ParticipantRole string

const (
	ParticipantRoleHost   ParticipantRole = "host"   // Recruitment host
	ParticipantRoleMember ParticipantRole = "member" // Regular participant
	ParticipantRoleBackup ParticipantRole = "backup" // Backup participant
)

// Participant represents a participant in a recruitment
type Participant struct {
	UserID      string          `json:"user_id"`
	Username    string          `json:"username"`
	Role        ParticipantRole `json:"role"`
	JoinedAt    time.Time       `json:"joined_at"`
	IsConfirmed bool            `json:"is_confirmed"`
}

// Recruitment represents a battle recruitment
type Recruitment struct {
	ID          string            `json:"id"`
	MessageID   string            `json:"message_id"`
	ChannelID   string            `json:"channel_id"`
	GuildID     string            `json:"guild_id"`
	BattleID    string            `json:"battle_id"`
	HostUserID  string            `json:"host_user_id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      RecruitmentStatus `json:"status"`
	MaxPlayers  int               `json:"max_players"`
	MinRank     int               `json:"min_rank"`

	// Participants
	Participants []Participant `json:"participants"`

	// Scheduling
	ScheduledTime *time.Time `json:"scheduled_time,omitempty"`

	// Creation and update times
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// RecruitmentManager manages battle recruitments
type RecruitmentManager struct {
	recruitments  map[string]*Recruitment
	battleManager *BattleManager
}

// NewRecruitmentManager creates a new recruitment manager
func NewRecruitmentManager(battleManager *BattleManager) *RecruitmentManager {
	return &RecruitmentManager{
		recruitments:  make(map[string]*Recruitment),
		battleManager: battleManager,
	}
}

// CreateRecruitment creates a new recruitment
func (rm *RecruitmentManager) CreateRecruitment(req *Recruitment) error {
	if req.ID == "" {
		return fmt.Errorf("recruitment ID cannot be empty")
	}

	// Validate battle exists
	battle, err := rm.battleManager.GetBattle(req.BattleID)
	if err != nil {
		return fmt.Errorf("invalid battle ID: %w", err)
	}

	// Set defaults
	now := time.Now()
	req.CreatedAt = now
	req.UpdatedAt = now
	req.Status = RecruitmentStatusOpen

	// Set expiration time (24 hours from creation)
	if req.ExpiresAt.IsZero() {
		req.ExpiresAt = now.Add(24 * time.Hour)
	}

	// Set max players from battle info if not specified
	if req.MaxPlayers == 0 {
		req.MaxPlayers = battle.MaxPlayers
	}

	// Set min rank from battle info if not specified
	if req.MinRank == 0 {
		req.MinRank = battle.MinRank
	}

	// Add host as first participant
	if req.HostUserID != "" {
		host := Participant{
			UserID:      req.HostUserID,
			Role:        ParticipantRoleHost,
			JoinedAt:    now,
			IsConfirmed: true,
		}
		req.Participants = []Participant{host}
	}

	rm.recruitments[req.ID] = req
	return nil
}

// GetRecruitment retrieves a recruitment by ID
func (rm *RecruitmentManager) GetRecruitment(id string) (*Recruitment, error) {
	recruitment, exists := rm.recruitments[id]
	if !exists {
		return nil, fmt.Errorf("recruitment not found: %s", id)
	}
	return recruitment, nil
}

// GetRecruitmentByMessage retrieves a recruitment by message ID
func (rm *RecruitmentManager) GetRecruitmentByMessage(messageID string) (*Recruitment, error) {
	for _, recruitment := range rm.recruitments {
		if recruitment.MessageID == messageID {
			return recruitment, nil
		}
	}
	return nil, fmt.Errorf("recruitment not found for message: %s", messageID)
}

// GetActiveRecruitments returns all active recruitments
func (rm *RecruitmentManager) GetActiveRecruitments() []*Recruitment {
	var activeRecruitments []*Recruitment
	for _, recruitment := range rm.recruitments {
		if recruitment.Status == RecruitmentStatusOpen || recruitment.Status == RecruitmentStatusFull {
			activeRecruitments = append(activeRecruitments, recruitment)
		}
	}
	return activeRecruitments
}

// GetRecruitmentsByChannel returns recruitments for a specific channel
func (rm *RecruitmentManager) GetRecruitmentsByChannel(channelID string) []*Recruitment {
	var channelRecruitments []*Recruitment
	for _, recruitment := range rm.recruitments {
		if recruitment.ChannelID == channelID {
			channelRecruitments = append(channelRecruitments, recruitment)
		}
	}
	return channelRecruitments
}

// AddParticipant adds a participant to a recruitment
func (rm *RecruitmentManager) AddParticipant(recruitmentID, userID, username string) error {
	recruitment, exists := rm.recruitments[recruitmentID]
	if !exists {
		return fmt.Errorf("recruitment not found: %s", recruitmentID)
	}

	// Check if recruitment is open
	if recruitment.Status != RecruitmentStatusOpen {
		return fmt.Errorf("recruitment is not open for new participants")
	}

	// Check if user is already a participant
	for _, participant := range recruitment.Participants {
		if participant.UserID == userID {
			return fmt.Errorf("user is already a participant")
		}
	}

	// Check if recruitment is full
	if len(recruitment.Participants) >= recruitment.MaxPlayers {
		return fmt.Errorf("recruitment is full")
	}

	// Add participant
	participant := Participant{
		UserID:      userID,
		Username:    username,
		Role:        ParticipantRoleMember,
		JoinedAt:    time.Now(),
		IsConfirmed: false,
	}

	recruitment.Participants = append(recruitment.Participants, participant)
	recruitment.UpdatedAt = time.Now()

	// Update status if full
	if len(recruitment.Participants) >= recruitment.MaxPlayers {
		recruitment.Status = RecruitmentStatusFull
	}

	return nil
}

// RemoveParticipant removes a participant from a recruitment
func (rm *RecruitmentManager) RemoveParticipant(recruitmentID, userID string) error {
	recruitment, exists := rm.recruitments[recruitmentID]
	if !exists {
		return fmt.Errorf("recruitment not found: %s", recruitmentID)
	}

	// Find and remove participant
	for i, participant := range recruitment.Participants {
		if participant.UserID == userID {
			// Don't allow host to leave
			if participant.Role == ParticipantRoleHost {
				return fmt.Errorf("host cannot leave recruitment")
			}

			// Remove participant
			recruitment.Participants = append(recruitment.Participants[:i], recruitment.Participants[i+1:]...)
			recruitment.UpdatedAt = time.Now()

			// Update status if no longer full
			if recruitment.Status == RecruitmentStatusFull && len(recruitment.Participants) < recruitment.MaxPlayers {
				recruitment.Status = RecruitmentStatusOpen
			}

			return nil
		}
	}

	return fmt.Errorf("participant not found: %s", userID)
}

// UpdateRecruitmentStatus updates the status of a recruitment
func (rm *RecruitmentManager) UpdateRecruitmentStatus(recruitmentID string, status RecruitmentStatus) error {
	recruitment, exists := rm.recruitments[recruitmentID]
	if !exists {
		return fmt.Errorf("recruitment not found: %s", recruitmentID)
	}

	recruitment.Status = status
	recruitment.UpdatedAt = time.Now()
	return nil
}

// CleanupExpiredRecruitments removes expired recruitments
func (rm *RecruitmentManager) CleanupExpiredRecruitments() []*Recruitment {
	var expired []*Recruitment
	now := time.Now()

	for id, recruitment := range rm.recruitments {
		if now.After(recruitment.ExpiresAt) &&
			(recruitment.Status == RecruitmentStatusOpen || recruitment.Status == RecruitmentStatusFull) {
			recruitment.Status = RecruitmentStatusCancelled
			recruitment.UpdatedAt = now
			expired = append(expired, recruitment)
			delete(rm.recruitments, id)
		}
	}

	return expired
}

// GetParticipantCount returns the current number of participants
func (r *Recruitment) GetParticipantCount() int {
	return len(r.Participants)
}

// GetConfirmedParticipantCount returns the number of confirmed participants
func (r *Recruitment) GetConfirmedParticipantCount() int {
	count := 0
	for _, participant := range r.Participants {
		if participant.IsConfirmed {
			count++
		}
	}
	return count
}

// IsFull returns true if the recruitment is full
func (r *Recruitment) IsFull() bool {
	return len(r.Participants) >= r.MaxPlayers
}

// IsExpired returns true if the recruitment has expired
func (r *Recruitment) IsExpired() bool {
	return time.Now().After(r.ExpiresAt)
}

// CanJoin returns true if a user can join this recruitment
func (r *Recruitment) CanJoin(userID string) bool {
	// Check if already a participant
	for _, participant := range r.Participants {
		if participant.UserID == userID {
			return false
		}
	}

	// Check if recruitment is open and not full
	return r.Status == RecruitmentStatusOpen && !r.IsFull() && !r.IsExpired()
}

// GetHost returns the host participant
func (r *Recruitment) GetHost() *Participant {
	for i, participant := range r.Participants {
		if participant.Role == ParticipantRoleHost {
			return &r.Participants[i]
		}
	}
	return nil
}

// FormatRecruitmentInfo formats recruitment information for display
func FormatRecruitmentInfo(recruitment *Recruitment) string {
	statusEmoji := map[RecruitmentStatus]string{
		RecruitmentStatusOpen:      "🟢",
		RecruitmentStatusFull:      "🟡",
		RecruitmentStatusClosed:    "🔴",
		RecruitmentStatusCompleted: "✅",
		RecruitmentStatusCancelled: "❌",
	}

	emoji, exists := statusEmoji[recruitment.Status]
	if !exists {
		emoji = "❓"
	}

	return fmt.Sprintf("%s **%s**\n"+
		"Battle: %s | Players: %d/%d\n"+
		"Status: %s | Host: <@%s>\n"+
		"Created: %s",
		emoji, recruitment.Title, recruitment.BattleID,
		len(recruitment.Participants), recruitment.MaxPlayers,
		recruitment.Status, recruitment.HostUserID,
		recruitment.CreatedAt.Format("2006-01-02 15:04"))
}
