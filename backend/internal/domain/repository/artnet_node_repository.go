package repository

import "github.com/nasshu2916/dmx_viewer/internal/domain/model"

type ArtNetNodeRepository interface {
	Save(node *model.ArtNetNode)
	All() []*model.ArtNetNode
}
