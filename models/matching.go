package models

import (
	"math"
	"sort"

	"github.com/google/uuid"
)

// UserPair represents a matched pair of users with their similarity score
type UserPair struct {
	User1      User    `json:"user1"`
	User2      User    `json:"user2"`
	Similarity float64 `json:"similarity"`
}

// GenerateExampleUsers creates 10 example users with different score patterns
func GenerateExampleUsers() []User {
	subjects := []string{"Math", "Physics", "Chemistry", "Biology", "English"}

	// Define different score patterns
	patterns := [][]int{
		{95, 92, 88, 85, 90}, // High performer, consistent
		{75, 78, 72, 76, 74}, // Average, consistent
		{65, 85, 45, 90, 55}, // Highly variable
		{88, 85, 82, 86, 84}, // Good, very consistent
		{92, 88, 85, 80, 78}, // Declining pattern
		{75, 80, 85, 88, 90}, // Improving pattern
		{95, 65, 90, 70, 85}, // Alternating high-low
		{70, 72, 68, 71, 73}, // Low-average, consistent
		{85, 83, 87, 82, 86}, // Above average, consistent
		{90, 92, 88, 91, 89}, // High performer, slight variation
	}

	var users []User
	for i, pattern := range patterns {
		var scores []Score
		for j, score := range pattern {
			scores = append(scores, Score{
				Subject: subjects[j],
				Score:   score,
			})
		}

		users = append(users, User{
			ID:    uuid.New().String(),
			Email: generateEmail(i),
			Score: scores,
		})
	}

	return users
}

// Helper function to generate example emails
func generateEmail(i int) string {
	emails := []string{
		"alice.smith@example.com",
		"bob.jones@example.com",
		"carol.wilson@example.com",
		"david.brown@example.com",
		"emma.davis@example.com",
		"frank.miller@example.com",
		"grace.taylor@example.com",
		"henry.white@example.com",
		"isabel.clark@example.com",
		"jack.martin@example.com",
	}
	return emails[i]
}

// normalizeScores converts raw scores to relative scores (z-scores)
func normalizeScores(scores []Score) map[string]float64 {
	// Create a map of subject to score for easier processing
	scoreMap := make(map[string]float64)
	var total float64
	n := float64(len(scores))

	// Calculate mean
	for _, score := range scores {
		scoreMap[score.Subject] = float64(score.Score)
		total += float64(score.Score)
	}
	mean := total / n

	// Calculate standard deviation
	var sumSquaredDiff float64
	for _, score := range scores {
		diff := float64(score.Score) - mean
		sumSquaredDiff += diff * diff
	}
	stdDev := math.Sqrt(sumSquaredDiff / n)

	// If stdDev is 0, return raw scores divided by max score to avoid division by zero
	if stdDev == 0 {
		maxScore := 1.0
		for _, score := range scores {
			if float64(score.Score) > maxScore {
				maxScore = float64(score.Score)
			}
		}
		for subject, score := range scoreMap {
			scoreMap[subject] = score / maxScore
		}
		return scoreMap
	}

	// Calculate z-scores
	for subject, score := range scoreMap {
		scoreMap[subject] = (score - mean) / stdDev
	}

	return scoreMap
}

// calculateSimilarity calculates the similarity between two users based on their normalized scores
func calculateSimilarity(user1, user2 User) float64 {
	scores1 := normalizeScores(user1.Score)
	scores2 := normalizeScores(user2.Score)

	// Get all unique subjects
	subjects := make(map[string]bool)
	for _, score := range user1.Score {
		subjects[score.Subject] = true
	}
	for _, score := range user2.Score {
		subjects[score.Subject] = true
	}

	// If users don't have the same subjects, they're not a good match
	if len(scores1) != len(scores2) {
		return 0
	}

	// Calculate Euclidean distance between normalized scores
	var sumSquaredDiff float64
	for subject := range subjects {
		score1, ok1 := scores1[subject]
		score2, ok2 := scores2[subject]
		if !ok1 || !ok2 {
			return 0 // If either user doesn't have a score for this subject, they're not a match
		}
		diff := score1 - score2
		sumSquaredDiff += diff * diff
	}

	// Convert distance to similarity (1 / (1 + distance))
	// This gives us a similarity score between 0 and 1, where 1 is perfect match
	similarity := 1 / (1 + math.Sqrt(sumSquaredDiff))
	return similarity
}

// FindMatches finds the best matches for all users
func FindMatches(users []User, minSimilarity float64) []UserPair {
	var pairs []UserPair

	// Compare each user with every other user
	for i := 0; i < len(users); i++ {
		for j := i + 1; j < len(users); j++ {
			similarity := calculateSimilarity(users[i], users[j])
			if similarity >= minSimilarity {
				pairs = append(pairs, UserPair{
					User1:      users[i],
					User2:      users[j],
					Similarity: similarity,
				})
			}
		}
	}

	// Sort pairs by similarity (highest first)
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Similarity > pairs[j].Similarity
	})

	return pairs
}
