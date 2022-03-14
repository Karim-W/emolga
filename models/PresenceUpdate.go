package models

type PresenceUpdate struct {
	UserId           string           `json:"userId"`
	NotificationType string           `json:"notificationType"`
	NotifiedEntities []NotifiedEntity `json:"notifiedEntities"`
}
type NotifiedEntity struct {
	EntityType string `json:"entityType"`
	EntityId   string `json:"entityId"`
}
