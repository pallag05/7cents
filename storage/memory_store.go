package storage

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"allen_hackathon/models"

	"github.com/google/uuid"
)

var questions = []models.Question{
	{
		ID:        uuid.New().String(),
		Content:   "What is the value of g (acceleration due to gravity) on Earth?",
		Options:   []string{"9.8 m/s²", "8.9 m/s²", "10.2 m/s²", "7.8 m/s²"},
		Timestamp: time.Now(),
	},
	{
		ID:        uuid.New().String(),
		Content:   "Which of these is a noble gas?",
		Options:   []string{"Helium", "Oxygen", "Nitrogen", "Carbon"},
		Timestamp: time.Now(),
	},
	{
		ID:        uuid.New().String(),
		Content:   "What is the derivative of sin(x)?",
		Options:   []string{"cos(x)", "-sin(x)", "tan(x)", "-cos(x)"},
		Timestamp: time.Now(),
	},
	{
		ID:        uuid.New().String(),
		Content:   "What is the first law of thermodynamics?",
		Options:   []string{"Energy cannot be created or destroyed", "Heat flows from hot to cold", "Entropy always increases", "Work equals force times distance"},
		Timestamp: time.Now(),
	},
	{
		ID:        uuid.New().String(),
		Content:   "What is the pH of a neutral solution?",
		Options:   []string{"7", "0", "14", "1"},
		Timestamp: time.Now(),
	},
}

type MemoryStore struct {
	users      map[string]*models.User
	groups     map[string]*models.Group
	userGroups map[string]*models.UserGroup
	matches    map[string]*models.UserPair // key: match ID
}

func NewMemoryStore() *MemoryStore {
	store := &MemoryStore{
		users:      make(map[string]*models.User),
		groups:     make(map[string]*models.Group),
		userGroups: make(map[string]*models.UserGroup),
		matches:    make(map[string]*models.UserPair),
	}

	// Add dummy questions

	// Add dummy users
	dummyUsers := []struct {
		email    string
		subjects []string
		scores   []int
		id       string
		name     string
	}{
		{
			email:    "alice.smith@example.com",
			subjects: []string{"physics", "chemistry", "maths"},
			scores:   []int{95, 92, 90},
			id:       "1",
			name:     "Alice Smith",
		},
		{
			email:    "bob.jones@example.com",
			subjects: []string{"physics", "chemistry", "maths"},
			scores:   []int{75, 78, 72},
			id:       "2",
			name:     "Bob Jones",
		},
		{
			email:    "carol.wilson@example.com",
			subjects: []string{"physics", "chemistry", "maths"},
			scores:   []int{85, 45, 90},
			id:       "3",
			name:     "Carol Wilson",
		},
		{
			email:    "david.brown@example.com",
			subjects: []string{"physics", "chemistry", "maths"},
			scores:   []int{88, 82, 86},
			id:       "4",
			name:     "David Brown",
		},
		{
			email:    "emma.davis@example.com",
			subjects: []string{"physics", "chemistry", "maths"},
			scores:   []int{92, 85, 78},
			id:       "5",
			name:     "Emma Davis",
		},
	}

	// Map to store users by email for easy lookup when creating matches
	usersByEmail := make(map[string]*models.User)

	for _, du := range dummyUsers {
		var scores []models.Score
		for i, subject := range du.subjects {
			scores = append(scores, models.Score{
				Subject: subject,
				Score:   du.scores[i],
			})
		}

		user := &models.User{
			ID:    du.id,
			Email: du.email,
			Score: scores,
			Name:  du.name,
		}
		store.users[user.ID] = user
		usersByEmail[user.Email] = user

		// Create UserGroup for each user
		userGroup := &models.UserGroup{
			ID:                du.id + "group",
			UserID:            user.ID,
			ActiveGroups:      []string{},
			RecommendedGroups: []string{},
		}
		store.userGroups[user.ID] = userGroup
	}

	// Add dummy matches with predefined pairs
	dummyMatches := []struct {
		email1     string
		email2     string
		similarity float64
		reason     string
		subject    string
		tag        string // Primary subject for the pair
	}{
		{
			email1:     "alice.smith@example.com",
			email2:     "emma.davis@example.com",
			similarity: 0.95,
			reason:     "Both peers are high performers across all subjects",
			subject:    "physics",
			tag:        "High Performers",
		},
		{
			email1:     "david.brown@example.com",
			email2:     "emma.davis@example.com",
			similarity: 0.90,
			reason:     "Both peers show a similar consistent performance pattern",
			subject:    "chemistry",
			tag:        "Consistent Performers",
		},
		{
			email1:     "alice.smith@example.com",
			email2:     "david.brown@example.com",
			similarity: 0.88,
			reason:     "Both peers are strong in physics and are overall consistent",
			subject:    "physics",
			tag:        "Overall Consistent",
		},
		{
			email1:     "bob.jones@example.com",
			email2:     "carol.wilson@example.com",
			similarity: 0.85,
			reason:     "Both peers have complementary strengths in different subjects",
			subject:    "maths",
			tag:        "Complementary Strengths",
		},
	}

	// Create the matches and pair study groups
	for _, dm := range dummyMatches {
		user1 := usersByEmail[dm.email1]
		user2 := usersByEmail[dm.email2]
		if user1 != nil && user2 != nil {
			// Create match
			match := &models.UserPair{
				User1:      *user1,
				User2:      *user2,
				Similarity: dm.similarity,
			}
			matchID := uuid.New().String()
			store.matches[matchID] = match

			// Create a private study group for the pair
			pairGroup := &models.Group{
				ID:                   user1.ID + user2.ID + "group",
				Title:                fmt.Sprintf("Connect for %s", dm.subject),
				Description:          fmt.Sprintf("Private study group for matched pair (%.0f%% similarity)", dm.similarity*100),
				Tag:                  dm.subject,
				Type:                 "Pair Study",
				Private:              true,
				Messages:             []models.Message{},
				CreateBy:             user1.ID,
				Capacity:             2,
				ActivityScore:        int(dm.similarity * 100),
				RecommendationReason: dm.reason,
				RecommendationTag:    dm.tag,
			}

			// Add welcome message
			welcomeMsg := models.Message{
				ID:        uuid.New().String(),
				Content:   fmt.Sprintf("Welcome to your paired study group! %s", dm.reason),
				SenderId:  pairGroup.CreateBy,
				Timestamp: time.Now(),
			}
			pairGroup.Messages = append(pairGroup.Messages, welcomeMsg)
			// Add group to store
			store.groups[pairGroup.ID] = pairGroup

			welcomeMsg.Content = fmt.Sprintf("You were matched based on weak performance in " + dm.subject)

			// Add group to recommended groups for both users
			if userGroup1, exists := store.userGroups[user1.ID]; exists {
				userGroup1.RecommendedGroups = append(userGroup1.RecommendedGroups, pairGroup.ID)
			}
			if userGroup2, exists := store.userGroups[user2.ID]; exists {
				userGroup2.RecommendedGroups = append(userGroup2.RecommendedGroups, pairGroup.ID)
			}

		}
	}

	pairGroup1 := &models.Group{
		ID:          "15" + "group" + "2",
		Title:       fmt.Sprintf("Connect for Thermodynamics"),
		Description: fmt.Sprintf("Public study group for topic weakness Thermodynamics"),
		Tag:         "Thermodynamics",
		Type:        "Topic Weakness",
		Private:     false,
		Messages: []models.Message{{
			ID:        uuid.New().String(),
			Content:   fmt.Sprintf("Welcome to your paired study group! You were matched based on weak performance in Physics"),
			SenderId:  "1",
			Timestamp: time.Now(),
		}},
		CreateBy:             "1",
		Capacity:             10,
		ActivityScore:        85,
		RecommendationReason: "You were matched based on weak performance in Physics",
		RecommendationTag:    "Weak Performance",
	}

	store.groups[pairGroup1.ID] = pairGroup1
	// Add group to recommended groups for both users
	if userGroup1, exists := store.userGroups["1"]; exists {
		userGroup1.RecommendedGroups = append(userGroup1.RecommendedGroups, pairGroup1.ID)
	}

	// Add dummy groups
	subjects := []string{"physics", "chemistry", "maths"}
	for i, subject := range subjects {
		group := &models.Group{
			ID:            "3" + strconv.Itoa(i) + "group",
			Title:         subject + " Study Group",
			Description:   "A group for studying " + subject,
			Members:       []string{},
			Tag:           subject,
			Type:          "study",
			Private:       false,
			Messages:      []models.Message{},
			CreateBy:      uuid.New().String(),
			Capacity:      10,
			ActivityScore: 100 - (i * 25), // 100, 75, 50
		}
		store.groups[group.ID] = group
	}

	// Add one more physics group with different activity score
	physicsGroup2 := &models.Group{
		ID:            "41" + "group",
		Title:         "Advanced Physics Group",
		Description:   "Advanced physics study group",
		Members:       []string{},
		Tag:           "physics",
		Type:          "study",
		Private:       false,
		Messages:      []models.Message{},
		CreateBy:      uuid.New().String(),
		Capacity:      10,
		ActivityScore: 85,
	}
	store.groups[physicsGroup2.ID] = physicsGroup2

	// Add some dummy messages to groups
	for _, group := range store.groups {
		message := models.Message{
			ID:        uuid.New().String(),
			Content:   "Welcome to " + group.Title,
			SenderId:  group.CreateBy,
			Timestamp: time.Now(),
		}
		group.Messages = append(group.Messages, message)
	}

	// Assign users to groups based on their scores
	/*for _, user := range store.users {
		for _, group := range store.groups {
			// Add users to groups if they have a good score in the subject
			var userScore int
			for _, score := range user.Score {
				if score.Subject == group.Tag {
					userScore = score.Score
					break
				}
			}
			if userScore >= 80 {
				group.Members = append(group.Members, user.ID)
				if userGroup, exists := store.userGroups[user.ID]; exists {
					userGroup.ActiveGroups = append(userGroup.ActiveGroups, group.ID)
				}
			}
		}
	}

	// Generate matches between users based on similar scores
	store.generateInitialMatches()*/

	return store
}

// generateInitialMatches creates initial matches between users based on score similarity
func (s *MemoryStore) generateInitialMatches() {
	users := make([]*models.User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}

	// Compare each user with every other user
	for i := 0; i < len(users); i++ {
		for j := i + 1; j < len(users); j++ {
			similarity := s.calculateSimilarity(users[i], users[j])
			if similarity >= 0.8 { // Only create matches for users with high similarity
				match := &models.UserPair{
					User1:      *users[i],
					User2:      *users[j],
					Similarity: similarity,
				}
				matchID := uuid.New().String()
				s.matches[matchID] = match
			}
		}
	}
}

// calculateSimilarity computes how similar two users are based on their scores
func (s *MemoryStore) calculateSimilarity(user1, user2 *models.User) float64 {
	// Create maps of subject to score for easier comparison
	scores1 := make(map[string]int)
	scores2 := make(map[string]int)

	for _, score := range user1.Score {
		scores1[score.Subject] = score.Score
	}
	for _, score := range user2.Score {
		scores2[score.Subject] = score.Score
	}

	// Check if they have the same subjects
	if len(scores1) != len(scores2) {
		return 0
	}

	// Calculate similarity using normalized score differences
	var totalDiff float64
	var maxPossibleDiff float64
	for subject, score1 := range scores1 {
		if score2, exists := scores2[subject]; exists {
			diff := float64(abs(score1 - score2))
			totalDiff += diff
			maxPossibleDiff += 100 // Maximum possible difference in scores
		} else {
			return 0 // Different subjects
		}
	}

	// Convert to similarity score (1 is most similar, 0 is least similar)
	if maxPossibleDiff == 0 {
		return 0
	}
	return 1 - (totalDiff / maxPossibleDiff)
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// GetMatches returns all matches for a specific user
func (s *MemoryStore) GetMatches(userID string) []*models.UserPair {
	var userMatches []*models.UserPair
	for _, match := range s.matches {
		if match.User1.ID == userID || match.User2.ID == userID {
			userMatches = append(userMatches, match)
		}
	}

	// Sort matches by similarity score (highest first)
	sort.Slice(userMatches, func(i, j int) bool {
		return userMatches[i].Similarity > userMatches[j].Similarity
	})

	return userMatches
}

// GetAllMatches returns all matches in the system
func (s *MemoryStore) GetAllMatches() []*models.UserPair {
	matches := make([]*models.UserPair, 0, len(s.matches))
	for _, match := range s.matches {
		matches = append(matches, match)
	}

	// Sort by similarity score
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Similarity > matches[j].Similarity
	})

	return matches
}

// CreateMatch creates a new match between two users
func (s *MemoryStore) CreateMatch(user1ID, user2ID string) (*models.UserPair, error) {
	user1, exists1 := s.users[user1ID]
	user2, exists2 := s.users[user2ID]
	if !exists1 || !exists2 {
		return nil, nil
	}

	similarity := s.calculateSimilarity(user1, user2)
	match := &models.UserPair{
		User1:      *user1,
		User2:      *user2,
		Similarity: similarity,
	}

	matchID := uuid.New().String()
	s.matches[matchID] = match
	return match, nil
}

// DeleteMatch removes a match from the system
func (s *MemoryStore) DeleteMatch(matchID string) error {
	delete(s.matches, matchID)
	return nil
}

// User operations
func (s *MemoryStore) GetUser(id string) (*models.User, error) {
	if user, exists := s.users[id]; exists {
		return user, nil
	}
	return nil, nil
}

func (s *MemoryStore) CreateUser(user *models.User) error {
	s.users[user.ID] = user
	return nil
}

func (s *MemoryStore) UpdateUser(user *models.User) error {
	s.users[user.ID] = user
	return nil
}

func (s *MemoryStore) DeleteUser(id string) error {
	delete(s.users, id)
	return nil
}

// Group operations
func (s *MemoryStore) GetGroup(id string) (*models.Group, error) {
	if group, exists := s.groups[id]; exists {
		group.Questions = questions
		return group, nil
	}
	return nil, nil
}

func (s *MemoryStore) CreateGroup(group *models.Group) error {
	s.groups[group.ID] = group
	return nil
}

func (s *MemoryStore) UpdateGroup(group *models.Group) error {
	s.groups[group.ID] = group
	return nil
}

func (s *MemoryStore) DeleteGroup(id string) error {
	delete(s.groups, id)
	return nil
}

func (s *MemoryStore) GetGroupsByUser(userID string) ([]*models.Group, error) {
	var userGroups []*models.Group
	for _, group := range s.groups {
		for _, memberID := range group.Members {
			if memberID == userID {
				userGroups = append(userGroups, group)
				break
			}
		}
	}
	return userGroups, nil
}

func (s *MemoryStore) AddMemberToGroup(groupID string, userID string) error {
	group, err := s.GetGroup(groupID)
	if err != nil {
		return err
	}
	if group == nil {
		return nil
	}

	// Check if user is already a member
	for _, memberID := range group.Members {
		if memberID == userID {
			return nil
		}
	}

	group.Members = append(group.Members, userID)
	return s.UpdateGroup(group)
}

func (s *MemoryStore) RemoveMemberFromGroup(groupID string, userID string) error {
	group, err := s.GetGroup(groupID)
	if err != nil {
		return err
	}
	if group == nil {
		return nil
	}

	// Remove user from members list
	for i, memberID := range group.Members {
		if memberID == userID {
			group.Members = append(group.Members[:i], group.Members[i+1:]...)
			return s.UpdateGroup(group)
		}
	}
	return nil
}

func (s *MemoryStore) AddMessageToGroup(groupID string, message *models.Message) error {
	group, err := s.GetGroup(groupID)
	if err != nil {
		return err
	}
	if group == nil {
		return nil
	}

	group.Messages = append(group.Messages, *message)
	return s.UpdateGroup(group)
}

// UserGroup operations
func (s *MemoryStore) GetUserGroup(userID string) (*models.UserGroup, error) {
	if userGroup, exists := s.userGroups[userID]; exists {
		return userGroup, nil
	}
	return nil, nil
}

func (s *MemoryStore) CreateUserGroup(userGroup *models.UserGroup) error {
	s.userGroups[userGroup.UserID] = userGroup
	return nil
}

func (s *MemoryStore) UpdateUserGroup(userGroup *models.UserGroup) error {
	s.userGroups[userGroup.UserID] = userGroup
	return nil
}

func (s *MemoryStore) GetGroupsByIDs(groupIDs []string) ([]*models.Group, error) {
	var groups []*models.Group
	for _, id := range groupIDs {
		if group, exists := s.groups[id]; exists {
			groups = append(groups, group)
		}
	}
	return groups, nil
}

func (s *MemoryStore) AddActionToGroup(groupID string, action *models.Action) error {
	group, err := s.GetGroup(groupID)
	if err != nil {
		return err
	}
	if group == nil {
		return nil
	}

	group.Actions = append(group.Actions, *action)
	return s.UpdateGroup(group)
}

func (s *MemoryStore) SearchGroupsByTag(tag string) []*models.Group {
	var matchingGroups []*models.Group
	for _, group := range s.groups {
		if group != nil && group.Tag == tag && group.Private == false {
			matchingGroups = append(matchingGroups, group)
		}
	}

	// Sort by activity score in descending order
	sort.Slice(matchingGroups, func(i, j int) bool {
		return matchingGroups[i].ActivityScore > matchingGroups[j].ActivityScore
	})

	return matchingGroups
}
