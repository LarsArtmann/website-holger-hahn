// Package infrastructure provides concrete implementations of external service interfaces.
// It contains in-memory contact repository implementation for development and testing
// with thread-safe operations and CRUD functionality for contact management.
package infrastructure

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"holger-hahn-website/internal/domain"
)

// MemoryContactRepository is an in-memory implementation of ContactRepository.
type MemoryContactRepository struct {
	contacts map[string]*domain.Contact
	mu       sync.RWMutex
}

// NewMemoryContactRepository creates a new in-memory contact repository.
func NewMemoryContactRepository() *MemoryContactRepository {
	return &MemoryContactRepository{
		contacts: make(map[string]*domain.Contact),
	}
}

// Save stores a contact submission.
func (r *MemoryContactRepository) Save(ctx context.Context, contact *domain.Contact) error {
	if contact == nil {
		return fmt.Errorf("%w", domain.ErrContactNil)
	}

	if err := contact.IsValid(); err != nil {
		return fmt.Errorf("%w: %w", domain.ErrInvalidContact, err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Create a copy to avoid external mutations
	contactCopy := *contact
	r.contacts[contact.ID] = &contactCopy

	return nil
}

// FindByID retrieves a contact by ID.
func (r *MemoryContactRepository) FindByID(ctx context.Context, id string) (*domain.Contact, error) {
	if id == "" {
		return nil, fmt.Errorf("%w", domain.ErrIDEmpty)
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	contact, exists := r.contacts[id]
	if !exists {
		return nil, fmt.Errorf("%w", domain.ErrContactNotFound)
	}

	// Return a copy to avoid external mutations
	contactCopy := *contact

	return &contactCopy, nil
}

// FindAll retrieves all contacts with optional filtering.
func (r *MemoryContactRepository) FindAll(ctx context.Context, status domain.ContactStatus, limit, offset int) ([]*domain.Contact, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*domain.Contact

	for _, contact := range r.contacts {
		if status == "" || domain.ContactStatus(contact.Status) == status {
			// Create a copy to avoid external mutations
			contactCopy := *contact
			filtered = append(filtered, &contactCopy)
		}
	}

	// Sort by submission time (newest first)
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].SubmittedAt.After(filtered[j].SubmittedAt)
	})

	// Apply pagination
	start := offset
	if start < 0 {
		start = 0
	}

	if start >= len(filtered) {
		return []*domain.Contact{}, nil
	}

	end := start + limit
	if limit <= 0 || end > len(filtered) {
		end = len(filtered)
	}

	return filtered[start:end], nil
}

// Update updates an existing contact.
func (r *MemoryContactRepository) Update(ctx context.Context, contact *domain.Contact) error {
	if contact == nil {
		return fmt.Errorf("%w", domain.ErrContactNil)
	}

	if err := contact.IsValid(); err != nil {
		return fmt.Errorf("%w: %w", domain.ErrInvalidContact, err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.contacts[contact.ID]; !exists {
		return fmt.Errorf("%w", domain.ErrContactNotFound)
	}

	// Create a copy to avoid external mutations
	contactCopy := *contact
	r.contacts[contact.ID] = &contactCopy

	return nil
}

// Delete removes a contact (for GDPR compliance).
func (r *MemoryContactRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("%w", domain.ErrIDEmpty)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.contacts[id]; !exists {
		return fmt.Errorf("%w", domain.ErrContactNotFound)
	}

	delete(r.contacts, id)

	return nil
}

// Count returns the total number of contacts.
func (r *MemoryContactRepository) Count(ctx context.Context, status domain.ContactStatus) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if status == "" {
		return len(r.contacts), nil
	}

	count := 0

	for _, contact := range r.contacts {
		if domain.ContactStatus(contact.Status) == status {
			count++
		}
	}

	return count, nil
}
