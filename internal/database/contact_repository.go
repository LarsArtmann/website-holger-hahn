// Package database provides database repository implementations using sqlc generated code.
// This file implements the ContactRepository interface with SQLite backend.
package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"holger-hahn-website/internal/domain"
)

// ContactRepository implements domain.ContactRepository using sqlc generated code.
type ContactRepository struct {
	queries *Queries
}

// NewContactRepository creates a new database contact repository.
func NewContactRepository(queries *Queries) *ContactRepository {
	return &ContactRepository{
		queries: queries,
	}
}

// Save stores a contact submission.
func (r *ContactRepository) Save(ctx context.Context, contact *domain.Contact) error {
	if contact == nil {
		return fmt.Errorf("%w", domain.ErrContactNil)
	}

	if err := contact.IsValid(); err != nil {
		return fmt.Errorf("%w: %w", domain.ErrInvalidContact, err)
	}

	params := CreateContactParams{
		Name:    contact.Name,
		Email:   contact.Email,
		Company: nullStringFromString(contact.Company),
		Message: contact.Message,
		Subject: nullStringFromString(contact.Subject),
		Source:  nullStringFromString(contact.Source),
	}

	created, err := r.queries.CreateContact(ctx, params)
	if err != nil {
		return fmt.Errorf("%w: %v", domain.ErrSaveContact, err)
	}

	// Update the contact with the generated ID and timestamps
	contact.ID = created.ID
	if created.CreatedAt.Valid {
		contact.SubmittedAt = created.CreatedAt.Time
	}

	return nil
}

// FindByID retrieves a contact by ID.
func (r *ContactRepository) FindByID(ctx context.Context, id string) (*domain.Contact, error) {
	if id == "" {
		return nil, fmt.Errorf("%w", domain.ErrIDEmpty)
	}

	dbContact, err := r.queries.GetContact(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w", domain.ErrContactNotFound)
		}
		return nil, fmt.Errorf("%w: %v", domain.ErrFindContact, err)
	}

	return r.toDomainContact(dbContact), nil
}

// FindAll retrieves all contacts with optional filtering.
func (r *ContactRepository) FindAll(ctx context.Context, status domain.ContactStatus, limit, offset int) ([]*domain.Contact, error) {
	var dbContacts []Contact
	var err error

	if status == "" {
		// Get all contacts
		params := ListContactsParams{
			Limit:  int64(limit),
			Offset: int64(offset),
		}
		dbContacts, err = r.queries.ListContacts(ctx, params)
	} else {
		// Filter by status
		params := ListContactsByStatusParams{
			Status: sql.NullString{
				String: string(status),
				Valid:  true,
			},
			Limit:  int64(limit),
			Offset: int64(offset),
		}
		dbContacts, err = r.queries.ListContactsByStatus(ctx, params)
	}

	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrListContacts, err)
	}

	contacts := make([]*domain.Contact, len(dbContacts))
	for i, dbContact := range dbContacts {
		contacts[i] = r.toDomainContact(dbContact)
	}

	return contacts, nil
}

// Update updates an existing contact.
func (r *ContactRepository) Update(ctx context.Context, contact *domain.Contact) error {
	if contact == nil {
		return fmt.Errorf("%w", domain.ErrContactNil)
	}

	if err := contact.IsValid(); err != nil {
		return fmt.Errorf("%w: %w", domain.ErrInvalidContact, err)
	}

	// First check if contact exists
	_, err := r.queries.GetContact(ctx, contact.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("%w", domain.ErrContactNotFound)
		}
		return fmt.Errorf("%w: %v", domain.ErrFindContact, err)
	}

	// Update status
	params := UpdateContactStatusParams{
		Status: sql.NullString{
			String: contact.Status,
			Valid:  true,
		},
		ID: contact.ID,
	}

	updated, err := r.queries.UpdateContactStatus(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to update contact: %w", err)
	}

	// Update timestamps from database
	if updated.UpdatedAt.Valid {
		contact.ProcessedAt = &updated.UpdatedAt.Time
	}

	return nil
}

// Delete removes a contact (for GDPR compliance).
func (r *ContactRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("%w", domain.ErrIDEmpty)
	}

	// First check if contact exists
	_, err := r.queries.GetContact(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("%w", domain.ErrContactNotFound)
		}
		return fmt.Errorf("%w: %v", domain.ErrFindContact, err)
	}

	if err := r.queries.DeleteContact(ctx, id); err != nil {
		return fmt.Errorf("failed to delete contact: %w", err)
	}

	return nil
}

// Count returns the total number of contacts.
func (r *ContactRepository) Count(ctx context.Context, status domain.ContactStatus) (int, error) {
	var count int64
	var err error

	if status == "" {
		count, err = r.queries.CountContacts(ctx)
	} else {
		statusParam := sql.NullString{
			String: string(status),
			Valid:  true,
		}
		count, err = r.queries.CountContactsByStatus(ctx, statusParam)
	}

	if err != nil {
		return 0, fmt.Errorf("failed to count contacts: %w", err)
	}

	return int(count), nil
}

// toDomainContact converts a database Contact to a domain Contact.
func (r *ContactRepository) toDomainContact(dbContact Contact) *domain.Contact {
	contact := &domain.Contact{
		ID:      dbContact.ID,
		Name:    dbContact.Name,
		Email:   dbContact.Email,
		Message: dbContact.Message,
		Status:  stringFromNullString(dbContact.Status),
		Source:  stringFromNullString(dbContact.Source),
	}

	if dbContact.Company.Valid {
		contact.Company = dbContact.Company.String
	}

	if dbContact.Subject.Valid {
		contact.Subject = dbContact.Subject.String
	}

	if dbContact.CreatedAt.Valid {
		contact.SubmittedAt = dbContact.CreatedAt.Time
	} else {
		contact.SubmittedAt = time.Now().UTC()
	}

	if dbContact.UpdatedAt.Valid {
		contact.ProcessedAt = &dbContact.UpdatedAt.Time
	}

	return contact
}

// Helper functions for sql.Null* types
func nullStringFromString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func stringFromNullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
