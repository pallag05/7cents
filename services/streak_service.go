package services

import (
	"allen_hackathon/models"
	"allen_hackathon/storage"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
)

type StreakService struct {
	store *storage.MemoryStore
}

// RatingBreakdown provides detailed information about rating calculation
type RatingBreakdown struct {
	BaseScore           float64 `json:"base_score"`
	StreakMultiplier    float64 `json:"streak_multiplier"`
	TypeMultiplier      float64 `json:"type_multiplier"`
	PenaltyPoints       float64 `json:"penalty_points"`
	FinalScore          float64 `json:"final_score"`
	CurrentRating       string  `json:"current_rating"`
	StreakCount         int     `json:"streak_count"`
	StreakType          string  `json:"streak_type"`
	LastStreakUpdated   string  `json:"last_streak_updated"`
	DaysSinceLastStreak int     `json:"days_since_last_streak"`
}

// LeaderboardEntry represents a user's position in the leaderboard
type LeaderboardEntry struct {
	Rank        int     `json:"rank"`
	UserID      string  `json:"user_id"`
	UserName    string  `json:"user_name"`
	Score       float64 `json:"score"`
	Rating      string  `json:"rating"`
	StreakCount int     `json:"streak_count"`
	BatchID     string  `json:"batch_id"`
}

// LeaderboardResponse represents the complete leaderboard response
type LeaderboardResponse struct {
	BatchID     string             `json:"batch_id"`
	TotalUsers  int                `json:"total_users"`
	Entries     []LeaderboardEntry `json:"entries"`
	LastUpdated string             `json:"last_updated"`
}

// Cache keys
const (
	cacheKeyBatchLeaderboard = "batch_leaderboard_%s_%s_%s_%s_%s"
	cacheKeyTopPerformers    = "top_performers_%s_%s_%s_%s_%s"
	cacheDuration            = 5 * time.Minute
)

// Cache structure
type leaderboardCache struct {
	entries     []LeaderboardEntry
	lastUpdated time.Time
}

var (
	batchLeaderboardCache = make(map[string]*leaderboardCache)
	topPerformersCache    = make(map[string]*leaderboardCache)
	cacheMutex            sync.RWMutex
)

// LeaderboardStats represents statistics for the leaderboard
type LeaderboardStats struct {
	TotalUsers    int     `json:"total_users"`
	AverageScore  float64 `json:"average_score"`
	HighestScore  float64 `json:"highest_score"`
	LowestScore   float64 `json:"lowest_score"`
	AverageStreak int     `json:"average_streak"`
	HighestStreak int     `json:"highest_streak"`
	LowestStreak  int     `json:"lowest_streak"`
	LastUpdated   string  `json:"last_updated"`
}

// RatingDistribution represents the distribution of ratings
type RatingDistribution struct {
	Novice       int    `json:"novice"`
	Beginner     int    `json:"beginner"`
	Intermediate int    `json:"intermediate"`
	Advanced     int    `json:"advanced"`
	Expert       int    `json:"expert"`
	LastUpdated  string `json:"last_updated"`
}

// StreakDistribution represents the distribution of streaks
type StreakDistribution struct {
	ZeroToFive     int    `json:"zero_to_five"`
	SixToTen       int    `json:"six_to_ten"`
	ElevenToTwenty int    `json:"eleven_to_twenty"`
	TwentyOnePlus  int    `json:"twenty_one_plus"`
	LastUpdated    string `json:"last_updated"`
}

func NewStreakService() *StreakService {
	return &StreakService{
		store: storage.GetStore(),
	}
}

// CalculateUserRating calculates the user's rating based on their streak performance
func (s *StreakService) CalculateUserRating(userID string) (string, error) {
	breakdown, err := s.GetRatingBreakdown(userID)
	if err != nil {
		return "", err
	}

	return breakdown.CurrentRating, nil
}

// GetRatingBreakdown provides detailed information about the rating calculation
func (s *StreakService) GetRatingBreakdown(userID string) (*RatingBreakdown, error) {
	streakToUser, exists := s.store.GetStreakToUser(userID)
	if !exists {
		return nil, errors.New("user not found")
	}

	streak, exists := s.store.GetStreak(streakToUser.CurrentStreakID)
	if !exists {
		return nil, errors.New("streak not found")
	}

	// Calculate base score (0-100)
	baseScore := calculateBaseScore(streakToUser.StreakCount)

	// Calculate streak multiplier (1.0-2.0)
	streakMultiplier := calculateStreakMultiplier(streakToUser.StreakCount)

	// Calculate type multiplier based on streak type
	typeMultiplier := calculateTypeMultiplier(streak.Type)

	// Calculate penalty points for missed streaks
	penaltyPoints := calculatePenaltyPoints(streakToUser.LastStreakUpdated)

	// Calculate final score
	finalScore := (baseScore * streakMultiplier * typeMultiplier) - penaltyPoints

	// Ensure final score is within bounds
	if finalScore < 0 {
		finalScore = 0
	}
	if finalScore > 100 {
		finalScore = 100
	}

	// Calculate current rating based on final score
	currentRating := calculateRatingFromScore(finalScore)

	// Calculate days since last streak
	daysSinceLastStreak := int(time.Since(streakToUser.LastStreakUpdated).Hours() / 24)

	return &RatingBreakdown{
		BaseScore:           baseScore,
		StreakMultiplier:    streakMultiplier,
		TypeMultiplier:      typeMultiplier,
		PenaltyPoints:       penaltyPoints,
		FinalScore:          finalScore,
		CurrentRating:       currentRating,
		StreakCount:         streakToUser.StreakCount,
		StreakType:          string(streak.Type),
		LastStreakUpdated:   streakToUser.LastStreakUpdated.Format(time.RFC3339),
		DaysSinceLastStreak: daysSinceLastStreak,
	}, nil
}

// Helper functions for rating calculation
func calculateBaseScore(streakCount int) float64 {
	// Base score increases with streak count, but with diminishing returns
	if streakCount <= 0 {
		return 0
	}
	if streakCount >= 30 {
		return 100
	}
	return float64(streakCount) * 3.33 // Linear scaling up to 30 days
}

func calculateStreakMultiplier(streakCount int) float64 {
	// Multiplier increases with streak count, but with diminishing returns
	if streakCount <= 0 {
		return 1.0
	}
	if streakCount >= 30 {
		return 2.0
	}
	return 1.0 + (float64(streakCount) * 0.033) // Linear scaling up to 30 days
}

func calculateTypeMultiplier(streakType models.StreakType) float64 {
	switch streakType {
	case models.StreakTypeBeginner:
		return 1.0
	case models.StreakTypeIntermediate:
		return 1.2
	case models.StreakTypeAdvanced:
		return 1.5
	default:
		return 1.0
	}
}

func calculatePenaltyPoints(lastStreakUpdated time.Time) float64 {
	daysSinceLastStreak := time.Since(lastStreakUpdated).Hours() / 24
	if daysSinceLastStreak <= 1 {
		return 0
	}

	// Exponential penalty for missed streaks
	penalty := 0.0
	for i := 1.0; i <= daysSinceLastStreak; i++ {
		penalty += i * 2 // Each day adds more penalty
	}
	return penalty
}

func calculateRatingFromScore(score float64) string {
	switch {
	case score >= 90:
		return "expert"
	case score >= 75:
		return "advanced"
	case score >= 50:
		return "intermediate"
	case score >= 25:
		return "beginner"
	default:
		return "novice"
	}
}

// RecordActivity records a learning activity and updates the user's streak
func (s *StreakService) RecordActivity(userID string, activityType models.StreakItemType) error {
	// 1. Get user's current streak info
	streakToUser, err := s.getUserStreak(userID)
	if err != nil {
		return err
	}

	// 2. Check if streak is frozen
	if streakToUser.IsFrozen {
		// If freeze has expired, unfreeze the streak
		if time.Now().After(streakToUser.FreezeEndTime) {
			freezeService := NewFreezeService()
			if err := freezeService.UnfreezeStreak(userID); err != nil {
				return err
			}
			streakToUser.IsFrozen = false
			streakToUser.FreezeEndTime = time.Time{}
		} else {
			// Streak is still frozen, don't update it
			return nil
		}
	}

	// 3. Check if the activity is within the same day
	if !isSameDay(streakToUser.LastStreakUpdated, time.Now()) {
		// 4. Check if streak should be broken
		if shouldBreakStreak(streakToUser.LastStreakUpdated) {
			streakToUser.StreakCount = 0
			streakToUser.CurrentStreakID = ""
		}
	}

	// 5. Get or create streak based on user's level
	streak, err := s.getOrCreateStreak(userID)
	if err != nil {
		return err
	}

	// 6. Create streak item
	streakItem := &models.StreakItem{
		ID:       generateID(),
		Type:     activityType,
		StreakID: streak.ID,
	}

	// 7. Update streak count and last updated time
	streakToUser.StreakCount++
	streakToUser.CurrentStreakID = streak.ID
	streakToUser.LastStreakUpdated = time.Now()

	// 8. Update user's rating if needed
	oldRating := streakToUser.CurrentRating
	if shouldUpdateRating(streakToUser) {
		streakToUser.CurrentRating = calculateNewRating(streakToUser)
		if streakToUser.CurrentRating > streakToUser.MaxRating {
			streakToUser.MaxRating = streakToUser.CurrentRating
		}

		// 9. Check for rewards if rating has changed
		if oldRating != streakToUser.CurrentRating {
			rewardService := NewRewardService()
			if err := rewardService.CheckAndAwardRewards(userID, streakToUser.CurrentRating); err != nil {
				// Log the error but don't fail the activity recording
				fmt.Printf("Error checking rewards: %v\n", err)
			}
		}
	}

	// 10. Save all changes
	return s.saveChanges(streakToUser, streakItem)
}

// GetUserStreakInfo returns the user's current streak information
func (s *StreakService) GetUserStreakInfo(userID string) (*models.StreakToUser, error) {
	return s.getUserStreak(userID)
}

// Helper functions
func (s *StreakService) getUserStreak(userID string) (*models.StreakToUser, error) {
	streakToUser, exists := s.store.GetStreakToUser(userID)
	if !exists {
		// Create new streak for user if not exists
		streakToUser = &models.StreakToUser{
			UserID:            userID,
			StreakCount:       0,
			CurrentStreakID:   "",
			CurrentRating:     "beginner",
			MaxRating:         "beginner",
			LastStreakUpdated: time.Now(),
		}
		s.store.SaveStreakToUser(streakToUser)
	}
	return streakToUser, nil
}

func (s *StreakService) getOrCreateStreak(userID string) (*models.Streak, error) {
	streak, exists := s.store.GetStreak(userID)
	if !exists {
		// Create new streak with default settings
		streak = &models.Streak{
			ID:                generateID(),
			Type:              models.StreakTypeBeginner,
			ThresholdDuration: 30, // 30 minutes default
		}
		s.store.SaveStreak(streak)
	}
	return streak, nil
}

func (s *StreakService) saveChanges(streakToUser *models.StreakToUser, streakItem *models.StreakItem) error {
	// Save streak item
	s.store.SaveStreakItem(streakItem)

	// Update streak to user
	s.store.SaveStreakToUser(streakToUser)

	return nil
}

func isSameDay(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}

func shouldBreakStreak(lastUpdated time.Time) bool {
	return time.Since(lastUpdated) > 24*time.Hour
}

func shouldUpdateRating(streakToUser *models.StreakToUser) bool {
	return streakToUser.StreakCount%5 == 0
}

func calculateNewRating(streakToUser *models.StreakToUser) string {
	switch streakToUser.CurrentRating {
	case "beginner":
		if streakToUser.StreakCount >= 10 {
			return "intermediate"
		}
	case "intermediate":
		if streakToUser.StreakCount >= 20 {
			return "advanced"
		}
	}
	return streakToUser.CurrentRating
}

func generateID() string {
	return uuid.New().String()
}

// GetBatchLeaderboard returns the leaderboard for a specific batch/class
func (s *StreakService) GetBatchLeaderboard(batchID, limitStr, offsetStr, rating string, startDate, endDate time.Time) (*LeaderboardResponse, error) {
	// Generate cache key
	cacheKey := fmt.Sprintf(cacheKeyBatchLeaderboard, batchID, limitStr, offsetStr, rating, startDate.Format(time.RFC3339))

	// Check cache first
	cacheMutex.RLock()
	if cache, exists := batchLeaderboardCache[cacheKey]; exists && time.Since(cache.lastUpdated) < cacheDuration {
		cacheMutex.RUnlock()
		return &LeaderboardResponse{
			BatchID:     batchID,
			TotalUsers:  len(cache.entries),
			Entries:     cache.entries,
			LastUpdated: cache.lastUpdated.Format(time.RFC3339),
		}, nil
	}
	cacheMutex.RUnlock()

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	// Get all users in the batch
	users := s.store.GetUsersByBatch(batchID)
	if len(users) == 0 {
		return nil, errors.New("no users found in batch")
	}

	// Calculate scores for all users
	var entries []LeaderboardEntry
	for _, user := range users {
		breakdown, err := s.GetRatingBreakdown(user.ID)
		if err != nil {
			continue // Skip users with no rating
		}

		// Apply filters
		if rating != "" && breakdown.CurrentRating != rating {
			continue
		}

		// Check date range if provided
		lastUpdated, err := time.Parse(time.RFC3339, breakdown.LastStreakUpdated)
		if err != nil {
			continue // Skip if we can't parse the date
		}
		if !startDate.IsZero() && lastUpdated.Before(startDate) {
			continue
		}
		if !endDate.IsZero() && lastUpdated.After(endDate) {
			continue
		}

		entries = append(entries, LeaderboardEntry{
			UserID:      user.ID,
			UserName:    user.Name,
			Score:       breakdown.FinalScore,
			Rating:      breakdown.CurrentRating,
			StreakCount: breakdown.StreakCount,
			BatchID:     batchID,
		})
	}

	// Sort entries by score in descending order
	sortLeaderboardEntries(entries)

	// Add ranks
	for i := range entries {
		entries[i].Rank = i + 1
	}

	// Apply pagination
	start := offset
	end := offset + limit
	if start >= len(entries) {
		start = 0
		end = 0
	}
	if end > len(entries) {
		end = len(entries)
	}

	// Update cache
	cacheMutex.Lock()
	batchLeaderboardCache[cacheKey] = &leaderboardCache{
		entries:     entries[start:end],
		lastUpdated: time.Now(),
	}
	cacheMutex.Unlock()

	return &LeaderboardResponse{
		BatchID:     batchID,
		TotalUsers:  len(entries),
		Entries:     entries[start:end],
		LastUpdated: time.Now().Format(time.RFC3339),
	}, nil
}

// GetTopPerformers returns the top performers across all batches
func (s *StreakService) GetTopPerformers(limitStr, offsetStr, rating string, startDate, endDate time.Time) ([]LeaderboardEntry, error) {
	// Generate cache key
	cacheKey := fmt.Sprintf(cacheKeyTopPerformers, limitStr, offsetStr, rating, startDate.Format(time.RFC3339))

	// Check cache first
	cacheMutex.RLock()
	if cache, exists := topPerformersCache[cacheKey]; exists && time.Since(cache.lastUpdated) < cacheDuration {
		cacheMutex.RUnlock()
		return cache.entries, nil
	}
	cacheMutex.RUnlock()

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	// Get all users
	users := s.store.GetAllUsers()
	if len(users) == 0 {
		return nil, errors.New("no users found")
	}

	// Calculate scores for all users
	var entries []LeaderboardEntry
	for _, user := range users {
		breakdown, err := s.GetRatingBreakdown(user.ID)
		if err != nil {
			continue // Skip users with no rating
		}

		// Apply filters
		if rating != "" && breakdown.CurrentRating != rating {
			continue
		}

		// Check date range if provided
		lastUpdated, err := time.Parse(time.RFC3339, breakdown.LastStreakUpdated)
		if err != nil {
			continue // Skip if we can't parse the date
		}
		if !startDate.IsZero() && lastUpdated.Before(startDate) {
			continue
		}
		if !endDate.IsZero() && lastUpdated.After(endDate) {
			continue
		}

		entries = append(entries, LeaderboardEntry{
			UserID:      user.ID,
			UserName:    user.Name,
			Score:       breakdown.FinalScore,
			Rating:      breakdown.CurrentRating,
			StreakCount: breakdown.StreakCount,
			BatchID:     user.BatchID,
		})
	}

	// Sort entries by score in descending order
	sortLeaderboardEntries(entries)

	// Add ranks
	for i := range entries {
		entries[i].Rank = i + 1
	}

	// Apply pagination
	start := offset
	end := offset + limit
	if start >= len(entries) {
		start = 0
		end = 0
	}
	if end > len(entries) {
		end = len(entries)
	}

	// Update cache
	cacheMutex.Lock()
	topPerformersCache[cacheKey] = &leaderboardCache{
		entries:     entries[start:end],
		lastUpdated: time.Now(),
	}
	cacheMutex.Unlock()

	return entries[start:end], nil
}

// Helper function to sort leaderboard entries
func sortLeaderboardEntries(entries []LeaderboardEntry) {
	// Sort by score in descending order
	for i := 0; i < len(entries)-1; i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[i].Score < entries[j].Score {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}
}

// GetLeaderboardStats returns statistics for the leaderboard
func (s *StreakService) GetLeaderboardStats(batchID string) (*LeaderboardStats, error) {
	users := s.store.GetUsersByBatch(batchID)
	if len(users) == 0 {
		return nil, errors.New("no users found in batch")
	}

	var totalScore float64
	var totalStreak int
	highestScore := -1.0
	lowestScore := 101.0
	highestStreak := 0
	lowestStreak := 999999

	for _, user := range users {
		breakdown, err := s.GetRatingBreakdown(user.ID)
		if err != nil {
			continue
		}

		score := breakdown.FinalScore
		streak := breakdown.StreakCount

		totalScore += score
		totalStreak += streak

		if score > highestScore {
			highestScore = score
		}
		if score < lowestScore {
			lowestScore = score
		}
		if streak > highestStreak {
			highestStreak = streak
		}
		if streak < lowestStreak {
			lowestStreak = streak
		}
	}

	return &LeaderboardStats{
		TotalUsers:    len(users),
		AverageScore:  totalScore / float64(len(users)),
		HighestScore:  highestScore,
		LowestScore:   lowestScore,
		AverageStreak: totalStreak / len(users),
		HighestStreak: highestStreak,
		LowestStreak:  lowestStreak,
		LastUpdated:   time.Now().Format(time.RFC3339),
	}, nil
}

// GetRatingDistribution returns the distribution of ratings in the leaderboard
func (s *StreakService) GetRatingDistribution(batchID string) (*RatingDistribution, error) {
	users := s.store.GetUsersByBatch(batchID)
	if len(users) == 0 {
		return nil, errors.New("no users found in batch")
	}

	distribution := &RatingDistribution{}

	for _, user := range users {
		breakdown, err := s.GetRatingBreakdown(user.ID)
		if err != nil {
			continue
		}

		switch breakdown.CurrentRating {
		case "novice":
			distribution.Novice++
		case "beginner":
			distribution.Beginner++
		case "intermediate":
			distribution.Intermediate++
		case "advanced":
			distribution.Advanced++
		case "expert":
			distribution.Expert++
		}
	}

	distribution.LastUpdated = time.Now().Format(time.RFC3339)
	return distribution, nil
}

// GetStreakDistribution returns the distribution of streaks in the leaderboard
func (s *StreakService) GetStreakDistribution(batchID string) (*StreakDistribution, error) {
	users := s.store.GetUsersByBatch(batchID)
	if len(users) == 0 {
		return nil, errors.New("no users found in batch")
	}

	distribution := &StreakDistribution{}

	for _, user := range users {
		breakdown, err := s.GetRatingBreakdown(user.ID)
		if err != nil {
			continue
		}

		streak := breakdown.StreakCount
		switch {
		case streak <= 5:
			distribution.ZeroToFive++
		case streak <= 10:
			distribution.SixToTen++
		case streak <= 20:
			distribution.ElevenToTwenty++
		default:
			distribution.TwentyOnePlus++
		}
	}

	distribution.LastUpdated = time.Now().Format(time.RFC3339)
	return distribution, nil
}

// GetStreak returns a streak by ID
func (s *StreakService) GetStreak(streakID string) (*models.Streak, bool) {
	return s.store.GetStreak(streakID)
}

// GetFreezeConfig returns the freeze configuration
func (s *StreakService) GetFreezeConfig() (*models.FreezeConfig, bool) {
	return s.store.GetFreezeConfig()
}
