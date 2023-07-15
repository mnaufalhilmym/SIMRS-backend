package patient

type gender string

const (
	GENDER_MALE   gender = "MALE"
	GENDER_FEMALE gender = "FEMALE"
)

type relationshipInFamily string

const (
	RELATION_HEADOFFAMILY relationshipInFamily = "HEAD_OF_FAMILY"
	RELATION_WIFE         relationshipInFamily = "WIFE"
	RELATION_CHILD        relationshipInFamily = "CHILD"
)
