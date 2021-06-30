module github.com/crockerio/user-service

go 1.16

require (
	github.com/crockerio/cservice v0.0.0
	gorm.io/gorm v1.21.11
)

replace github.com/crockerio/cservice => ./../cservice
