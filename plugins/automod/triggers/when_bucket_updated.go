package triggers

import (
	"errors"

	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/events"
)

type BucketUpdated struct {
}

func (t BucketUpdated) Name() string {
	return "when_bucket_updated"
}

func (t BucketUpdated) Args() int {
	return 1
}

func (t BucketUpdated) NewItem(env *models.Env, args []string) (interfaces.TriggerItemInterface, error) {
	if len(args) < 1 {
		return nil, errors.New("too few arguments")
	}

	return &BucketUpdatedItem{
		Tag: args[0],
	}, nil
}

func (t BucketUpdated) Description() string {
	return "automod.triggers.when_bucket_updated"
}

type BucketUpdatedItem struct {
	Tag string
}

func (t *BucketUpdatedItem) Match(env *models.Env) bool {
	if env.Event == nil {
		return false
	}

	if env.Event.Type != events.CacophonyBucketUpdate {
		return false
	}

	if env.Event.BucketUpdate.Tag != t.Tag {
		return false
	}

	return true
}
