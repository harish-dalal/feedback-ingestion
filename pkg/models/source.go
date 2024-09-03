package models

type Source string

const (
	SourceIntercom  Source = "intercom"
	SourcePlaystore Source = "playstore"
	SourceDiscourse Source = "discourse"
)

type SourceType string

const (
	STFeedback     SourceType = "feedback"
	STSurvey       SourceType = "survey"
	STConversation SourceType = "conversation"
	STReviews      SourceType = "reviews"
)
