package services

import (
	"allen_hackathon/models"
	"allen_hackathon/storage"
	"errors"
	"fmt"
	"sort"
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
	CurrentRating       float64 `json:"current_rating"`
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

// LeaderboardDistribution represents the distribution of ratings and streaks
type LeaderboardDistribution struct {
	RatingDistribution map[string]int     `json:"rating_distribution"`
	StreakDistribution map[string]int     `json:"streak_distribution"`
	RatingPercentiles  map[string]float64 `json:"rating_percentiles"`
	StreakPercentiles  map[string]float64 `json:"streak_percentiles"`
}

func NewStreakService() *StreakService {
	return &StreakService{
		store: storage.GetStore(),
	}
}

// CalculateUserRating calculates a user's rating based on their streak and activity
func (s *StreakService) CalculateUserRating(userID string) (float64, error) {
	streakToUser, exists := s.store.GetStreakToUser(userID)
	if !exists {
		return 0, errors.New("user not found")
	}

	streak, exists := s.store.GetStreak(streakToUser.CurrentStreakID)
	if !exists {
		return 0, errors.New("streak not found")
	}

	// Calculate base score
	baseScore := calculateBaseScore(streakToUser.StreakCount)

	// Calculate multipliers
	streakMultiplier := calculateStreakMultiplier(streakToUser.StreakCount)
	typeMultiplier := calculateTypeMultiplier(streak.Type)

	// Calculate penalty points
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

	return finalScore, nil
}

// GetRatingBreakdown returns detailed information about a user's rating calculation
func (s *StreakService) GetRatingBreakdown(userID string) (*RatingBreakdown, error) {
	streakToUser, exists := s.store.GetStreakToUser(userID)
	if !exists {
		return nil, errors.New("user not found")
	}

	streak, exists := s.store.GetStreak(streakToUser.CurrentStreakID)
	if !exists {
		return nil, errors.New("streak not found")
	}

	// Calculate base score
	baseScore := calculateBaseScore(streakToUser.StreakCount)

	// Calculate multipliers
	streakMultiplier := calculateStreakMultiplier(streakToUser.StreakCount)
	typeMultiplier := calculateTypeMultiplier(streak.Type)

	// Calculate penalty points
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

	// Calculate days since last streak
	lastStreakUpdated, err := time.Parse(time.RFC3339, streakToUser.LastStreakUpdated)
	if err != nil {
		return nil, fmt.Errorf("invalid last streak updated time: %v", err)
	}
	daysSinceLastStreak := int(time.Since(lastStreakUpdated).Hours() / 24)

	return &RatingBreakdown{
		BaseScore:           baseScore,
		StreakMultiplier:    streakMultiplier,
		TypeMultiplier:      typeMultiplier,
		PenaltyPoints:       penaltyPoints,
		FinalScore:          finalScore,
		CurrentRating:       finalScore,
		StreakCount:         streakToUser.StreakCount,
		StreakType:          string(streak.Type),
		LastStreakUpdated:   streakToUser.LastStreakUpdated,
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

func calculatePenaltyPoints(lastStreakUpdated string) float64 {
	lastUpdated, err := time.Parse(time.RFC3339, lastStreakUpdated)
	if err != nil {
		return 0
	}

	daysSinceLastStreak := time.Since(lastUpdated).Hours() / 24
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

// getAppropriateStreak returns the appropriate streak type based on streak count
func (s *StreakService) getAppropriateStreak(streakCount int) *models.Streak {
	var streakType models.StreakType
	switch {
	case streakCount >= 30:
		streakType = models.StreakTypeAdvanced
	case streakCount >= 15:
		streakType = models.StreakTypeIntermediate
	default:
		streakType = models.StreakTypeBeginner
	}

	// Find or create streak with the appropriate type
	for _, streak := range s.store.GetAllStreaks() {
		if streak.Type == streakType {
			return streak
		}
	}

	// Create new streak if not found
	streak := &models.Streak{
		ID:                uuid.New().String(),
		Type:              streakType,
		ThresholdDuration: getThresholdDuration(streakType),
		CreatedAt:         time.Now().Format(time.RFC3339),
		UpdatedAt:         time.Now().Format(time.RFC3339),
	}
	s.store.SaveStreak(streak)
	return streak
}

// getThresholdDuration returns the threshold duration for a streak type
func getThresholdDuration(streakType models.StreakType) int {
	switch streakType {
	case models.StreakTypeBeginner:
		return 30
	case models.StreakTypeIntermediate:
		return 45
	case models.StreakTypeAdvanced:
		return 60
	default:
		return 30
	}
}

// RecordActivity records a user's activity and updates their streak
func (s *StreakService) RecordActivity(userID string, activityType string) error {
	streakToUser, exists := s.store.GetStreakToUser(userID)
	if !exists {
		return errors.New("user not found")
	}

	// Check if streak is frozen
	if streakToUser.IsFrozen {
		freezeEndTime, err := time.Parse(time.RFC3339, streakToUser.FreezeEndTime)
		if err != nil {
			return fmt.Errorf("invalid freeze end time: %v", err)
		}

		if time.Now().After(freezeEndTime) {
			streakToUser.IsFrozen = false
			streakToUser.FreezeEndTime = ""
			streakToUser.UpdatedAt = time.Now().Format(time.RFC3339)
			s.store.SaveStreakToUser(streakToUser)
		} else {
			return errors.New("streak is frozen")
		}
	}

	// Check if activity is on the same day as last streak
	lastStreakUpdated, err := time.Parse(time.RFC3339, streakToUser.LastStreakUpdated)
	if err != nil {
		return fmt.Errorf("invalid last streak updated time: %v", err)
	}

	if !isSameDay(lastStreakUpdated, time.Now()) {
		// Check if streak should be broken
		if shouldBreakStreak(streakToUser.LastStreakUpdated) {
			streakToUser.StreakCount = 0
			streakToUser.CurrentStreakID = ""
			streakToUser.CurrentRating = 0
			streakToUser.MaxRating = 0
			streakToUser.LastStreakUpdated = time.Now().Format(time.RFC3339)
			streakToUser.UpdatedAt = time.Now().Format(time.RFC3339)
			s.store.SaveStreakToUser(streakToUser)
			return nil
		}

		// Update streak count and get appropriate streak type
		streakToUser.StreakCount++
		streak := s.getAppropriateStreak(streakToUser.StreakCount)
		if streak == nil {
			return errors.New("no appropriate streak found")
		}

		streakToUser.CurrentStreakID = streak.ID
		streakToUser.LastStreakUpdated = time.Now().Format(time.RFC3339)
		streakToUser.UpdatedAt = time.Now().Format(time.RFC3339)

		// Calculate and update rating
		rating, err := s.CalculateUserRating(userID)
		if err != nil {
			return err
		}

		streakToUser.CurrentRating = rating
		if rating > streakToUser.MaxRating {
			streakToUser.MaxRating = rating
		}

		s.store.SaveStreakToUser(streakToUser)

		// Check and award rewards based on the new rating
		rewardService := NewRewardService()
		if err := rewardService.CheckAndAwardRewards(userID, rating); err != nil {
			// Log the error but don't fail the streak update
			fmt.Printf("Error checking rewards: %v\n", err)
		}
	}

	return nil
}

// Helper function to check if two times are on the same day
func isSameDay(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}

// Helper function to check if streak should be broken
func shouldBreakStreak(lastUpdated string) bool {
	lastUpdatedTime, err := time.Parse(time.RFC3339, lastUpdated)
	if err != nil {
		return false
	}
	return time.Since(lastUpdatedTime) > 24*time.Hour
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
			CurrentRating:     1, // novice
			MaxRating:         1, // novice
			LastStreakUpdated: time.Now().Format(time.RFC3339),
		}
		s.store.SaveStreakToUser(streakToUser)
	}
	return streakToUser, nil
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
		if rating != "" {
			ratingFloat, err := strconv.ParseFloat(rating, 64)
			if err != nil {
				continue // Skip if rating filter is invalid
			}
			if breakdown.CurrentRating != ratingFloat {
				continue
			}
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
			Rating:      fmt.Sprintf("%.1f", breakdown.CurrentRating),
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
	cacheKey := fmt.Sprintf(cacheKeyTopPerformers, limitStr, offsetStr, rating, startDate.Format(time.RFC3339), endDate.Format(time.RFC3339))

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
		if rating != "" {
			ratingFloat, err := strconv.ParseFloat(rating, 64)
			if err != nil {
				continue // Skip if rating filter is invalid
			}
			if breakdown.CurrentRating != ratingFloat {
				continue
			}
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
			Rating:      fmt.Sprintf("%.1f", breakdown.CurrentRating),
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

		// Convert rating to category based on score
		switch {
		case breakdown.CurrentRating < 25:
			distribution.Novice++
		case breakdown.CurrentRating < 50:
			distribution.Beginner++
		case breakdown.CurrentRating < 75:
			distribution.Intermediate++
		case breakdown.CurrentRating < 90:
			distribution.Advanced++
		default:
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

// GetOverallLeaderboardStats returns overall statistics for all leaderboards
func (s *StreakService) GetOverallLeaderboardStats() (*models.LeaderboardStats, error) {
	allStreaks := s.store.GetAllStreakToUsers()
	if len(allStreaks) == 0 {
		return &models.LeaderboardStats{
			TotalUsers:    0,
			AverageRating: "0",
			MaxRating:     "0",
			MinRating:     "0",
			AverageStreak: 0,
			MaxStreak:     0,
			MinStreak:     0,
		}, nil
	}

	var totalRating float64
	var maxRating float64
	var minRating float64 = 999999
	var totalStreak int
	var maxStreak int
	var minStreak int = 999999

	for _, streak := range allStreaks {
		// Use the CurrentRating directly since it's already a float64
		rating := streak.CurrentRating
		totalRating += rating
		if rating > maxRating {
			maxRating = rating
		}
		if rating < minRating {
			minRating = rating
		}

		if streak.StreakCount > maxStreak {
			maxStreak = streak.StreakCount
		}
		if streak.StreakCount < minStreak {
			minStreak = streak.StreakCount
		}
		totalStreak += streak.StreakCount
	}

	avgRating := totalRating / float64(len(allStreaks))
	avgStreak := totalStreak / len(allStreaks)

	return &models.LeaderboardStats{
		TotalUsers:    len(allStreaks),
		AverageRating: fmt.Sprintf("%.2f", avgRating),
		MaxRating:     fmt.Sprintf("%.2f", maxRating),
		MinRating:     fmt.Sprintf("%.2f", minRating),
		AverageStreak: avgStreak,
		MaxStreak:     maxStreak,
		MinStreak:     minStreak,
	}, nil
}

// getStreakRange returns the range category for a streak count
func getStreakRange(count int) string {
	switch {
	case count == 0:
		return "0"
	case count <= 7:
		return "1-7"
	case count <= 14:
		return "8-14"
	case count <= 30:
		return "15-30"
	default:
		return "30+"
	}
}

// GetLeaderboardDistribution returns the distribution of ratings and streaks
func (s *StreakService) GetLeaderboardDistribution() (*models.LeaderboardDistribution, error) {
	users := s.store.GetAllUsers()
	if users == nil {
		return nil, fmt.Errorf("failed to get users")
	}

	distribution := &models.LeaderboardDistribution{
		RatingDistribution: make(map[string]int),
		StreakDistribution: make(map[string]int),
		RatingPercentiles:  make(map[string]float64),
		StreakPercentiles:  make(map[string]float64),
	}

	var ratings []float64
	var streaks []int

	for _, user := range users {
		streak, exists := s.store.GetStreakToUser(user.ID)
		if !exists {
			continue
		}

		// Use the CurrentRating directly since it's already a float64
		ratings = append(ratings, streak.CurrentRating)
		streaks = append(streaks, streak.StreakCount)

		// Count rating distribution
		rating := fmt.Sprintf("%.1f", streak.CurrentRating)
		distribution.RatingDistribution[rating]++

		// Count streak distribution
		streakRange := getStreakRange(streak.StreakCount)
		distribution.StreakDistribution[streakRange]++
	}

	// Calculate percentiles
	if len(ratings) > 0 {
		sort.Float64s(ratings)
		distribution.RatingPercentiles["25th"] = calculatePercentile(ratings, 25)
		distribution.RatingPercentiles["50th"] = calculatePercentile(ratings, 50)
		distribution.RatingPercentiles["75th"] = calculatePercentile(ratings, 75)
		distribution.RatingPercentiles["90th"] = calculatePercentile(ratings, 90)
	}

	if len(streaks) > 0 {
		sort.Ints(streaks)
		distribution.StreakPercentiles["25th"] = float64(calculateIntPercentile(streaks, 25))
		distribution.StreakPercentiles["50th"] = float64(calculateIntPercentile(streaks, 50))
		distribution.StreakPercentiles["75th"] = float64(calculateIntPercentile(streaks, 75))
		distribution.StreakPercentiles["90th"] = float64(calculateIntPercentile(streaks, 90))
	}

	return distribution, nil
}

// calculatePercentile calculates the percentile value for a slice of float64 values
func calculatePercentile(values []float64, percentile int) float64 {
	if len(values) == 0 {
		return 0
	}

	index := (percentile * (len(values) - 1)) / 100
	return values[index]
}

// calculateIntPercentile calculates the percentile value for a slice of int values
func calculateIntPercentile(values []int, percentile int) int {
	if len(values) == 0 {
		return 0
	}

	index := (percentile * (len(values) - 1)) / 100
	return values[index]
}
