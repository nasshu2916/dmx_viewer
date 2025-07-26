package infrastructure

import (
	"sync"

	"github.com/nasshu2916/dmx_viewer/internal/domain/model"
)

type ArtNetNodeRepositoryImpl struct {
	mu    sync.RWMutex
	nodes map[string]*model.ArtNetNode
}

func NewArtNetNodeRepository() *ArtNetNodeRepositoryImpl {
	return &ArtNetNodeRepositoryImpl{
		nodes: make(map[string]*model.ArtNetNode),
	}
}

func (r *ArtNetNodeRepositoryImpl) Save(node *model.ArtNetNode) {
	ip := node.IPAddress.String()
	r.mu.Lock()
	r.nodes[ip] = node
	r.mu.Unlock()
}

func (r *ArtNetNodeRepositoryImpl) All() []*model.ArtNetNode {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*model.ArtNetNode, 0, len(r.nodes))
	for _, n := range r.nodes {
		result = append(result, n)
	}
	return result
}
