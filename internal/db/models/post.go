package models

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"regexp"
	"strings"
	"time"
)

// PostDB is the database for the post data.
type PostDB struct {
	conn *gorm.DB
}

// NewPostDB creates a new PostDB.
func NewPostDB(conn *gorm.DB) *PostDB {
	return &PostDB{
		conn: conn,
	}
}

// Post is the model for the post-data.
type Post struct {
	ID          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Slug        string `json:"slug" gorm:"uniqueIndex"` // Slug is the URL friendly version of the title
	Title       string `json:"title"`
	Description string `json:"description"`
	Keywords    string `json:"keywords"` // Keywords are comma separated
	Content     string `json:"content"`
	UserID      int    `json:"user_id" gorm:"not null;constraint:OnUpdate:CASCADE;foreignKey:UserID;references:ID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p *Post) Validate() error {
	switch {
	case p.Slug == "":
		return ErrPostURLRequired
	case !regexp.MustCompile(`^[a-z0-9-]+$`).MatchString(p.Slug):
		return ErrPostInvalidSlug
	case p.Title == "":
		return ErrPostTitleRequired
	case p.Description == "":
		return ErrPostDescriptionRequired
	case p.Content == "":
		return ErrPostContentRequired
	case p.Keywords != "":
		keywords := strings.Split(p.Keywords, ",")
		for _, k := range keywords {
			if k == "" {
				return ErrPostWrongKeywordsString
			}
		}
	case p.UserID == 0:
		return ErrPostUserIDRequired
	}

	return nil
}

func (p *Post) BeforeCreate(_ *gorm.DB) error {
	err := p.Validate()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}

	return nil
}

func (p *Post) BeforeUpdate(_ *gorm.DB) error {
	err := p.Validate()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}

	p.UpdatedAt = time.Now()

	return nil
}

// Create creates a new Post.
func (db *PostDB) Create(ctx context.Context, p *Post) error {
	err := db.conn.WithContext(ctx).Create(p).Error
	if err != nil {
		return mapGormError(err)
	}

	return nil
}

// GetBySlug finds a Post by its URL Slug.
func (db *PostDB) GetBySlug(ctx context.Context, slug string) (*Post, error) {
	var p Post
	err := db.conn.WithContext(ctx).Where("slug = ?", slug).First(&p).Error
	if err != nil {
		return nil, mapGormError(err)
	}

	return &p, nil
}

// FindAll returns all the posts with pagination, sorted by the created time.
// Selects only the necessary fields to reduce the payload - slug, title, description, keywords, created_at.
func (db *PostDB) FindAll(ctx context.Context, page, perPage int) ([]*Post, error) {
	var posts []*Post
	err := db.conn.
		WithContext(ctx).
		Select("slug, title, description, keywords, created_at").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Order("created_at desc").
		Find(&posts).Error
	if err != nil {
		return nil, mapGormError(err)
	}

	return posts, nil
}

// Update updates the Post.
func (db *PostDB) Update(ctx context.Context, p *Post) error {
	err := db.conn.WithContext(ctx).Save(p).Error
	if err != nil {
		return mapGormError(err)
	}

	return nil
}
