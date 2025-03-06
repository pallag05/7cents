package models

// LeaderboardStats represents statistics for a leaderboard
type LeaderboardStats struct {
	TotalUsers         int            `json:"total_users"`
	AverageRating      string         `json:"average_rating"`
	MaxRating          string         `json:"max_rating"`
	MinRating          string         `json:"min_rating"`
	AverageStreak      int            `json:"average_streak"`
	MaxStreak          int            `json:"max_streak"`
	MinStreak          int            `json:"min_streak"`
	RatingDistribution map[string]int `json:"rating_distribution"`
	StreakDistribution map[string]int `json:"streak_distribution"`
}

// LeaderboardDistribution represents the distribution of ratings and streaks
type LeaderboardDistribution struct {
	RatingDistribution map[string]int     `json:"rating_distribution"`
	StreakDistribution map[string]int     `json:"streak_distribution"`
	RatingPercentiles  map[string]float64 `json:"rating_percentiles"`
	StreakPercentiles  map[string]float64 `json:"streak_percentiles"`
}
