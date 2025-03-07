package models

type UserGroup struct {
	ID                string   `json:"id"`
	UserID            string   `json:"user_id"`
	ActiveGroups      []string `json:"active_groups"`
	RecommendedGroups []string `json:"recommended_groups"`
}

type GroupsPageResponse struct {
	SystemRecommendedGroups []Group `json:"system_recommended_groups"`
	UserActiveGroups        []Group `json:"user_active_groups"`
	User                    User    `json:"user"`
}
