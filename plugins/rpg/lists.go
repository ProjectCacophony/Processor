package rpg

import (
	"gitlab.com/Cacophony/Processor/plugins/rpg/happening"
	"gitlab.com/Cacophony/Processor/plugins/rpg/models"
)

var happenings = []models.Happening{
	&happening.Walk{},
	&happening.MeetUser{},
}
