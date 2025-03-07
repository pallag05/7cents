package models

import (
	"math"
	"sort"
)

// UserPair represents a matched pair of users with their similarity score
type UserPair struct {
	User1      User    `json:"user1"`
	User2      User    `json:"user2"`
	Similarity float64 `json:"similarity"`
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
