package source

import (
	"fmt"
	"sync"
)

type memoryStore struct {
	sync.RWMutex
	resources map[string]*Resource
}

var store *memoryStore
var once sync.Once

func GetStore() *memoryStore {
	once.Do(func() {
		store = &memoryStore{
			resources: make(map[string]*Resource),
		}
	})

	return store
}

func (s *memoryStore) Add(resource *Resource) {
	s.Lock()
	defer s.Unlock()

	_, ok := s.resources[resource.ResourceID]
	if !ok {
		s.resources[resource.ResourceID] = resource
	}
}

func (s *memoryStore) Update(resource *Resource) error {
	s.Lock()
	defer s.Unlock()

	_, ok := s.resources[resource.ResourceID]
	if !ok {
		return fmt.Errorf("the resource %s does not exist", resource.ResourceID)
	}

	s.resources[resource.ResourceID] = resource
	return nil
}

func (s *memoryStore) UpdateStatus(resource *Resource) error {
	s.Lock()
	defer s.Unlock()

	last, ok := s.resources[resource.ResourceID]
	if !ok {
		return fmt.Errorf("the resource %s does not exist", resource.ResourceID)
	}

	last.Status = resource.Status
	s.resources[resource.ResourceID] = last
	return nil
}

func (s *memoryStore) Delete(resourceID string) {
	s.Lock()
	defer s.Unlock()

	delete(s.resources, resourceID)
}

func (s *memoryStore) Get(resourceID string) (*Resource, error) {
	s.RLock()
	defer s.RUnlock()

	resource, ok := s.resources[resourceID]
	if !ok {
		return nil, fmt.Errorf("failed to find resource %s", resourceID)
	}

	return resource, nil
}

func (s *memoryStore) List(namespace string) []*Resource {
	s.RLock()
	defer s.RUnlock()

	resources := []*Resource{}
	for _, res := range s.resources {
		if res.Namespace != namespace {
			continue
		}

		resources = append(resources, res)
	}
	return resources
}
