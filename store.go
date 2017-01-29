package hydrocarbon

import (
	"errors"
	"time"

	"github.com/fortytw2/abdi"
)

// Errors
var (
	ErrUserExists = errors.New("user already exists")
)

// Store is responsible for persistent (or not) data storage and retrieval
// and abstracting that into business-logic level functions
type Store struct {
	Users        UserStore
	Sessions     SessionStore
	Feeds        FeedStore
	Folders      FolderStore
	Posts        PostStore
	ReadStatuses ReadStatusStore

	EncryptionKey []byte
}

// PrimitiveStore encapsulates all primitive store types
type PrimitiveStore interface {
	UserStore
	SessionStore
	FeedStore
	PostStore
	FolderStore
	ReadStatusStore
}

// NewStore builds a data storage layer out of the persistence primitives
// It automatically sets and maintains all annotations such as "CreatedAt",
// "UpdatedAt", etc, but the underlying PrimitiveStore is equally allowed to
func NewStore(ps PrimitiveStore, encryptionKey []byte) (*Store, error) {
	return &Store{
		Users:         ps,
		Sessions:      ps,
		Feeds:         ps,
		Folders:       ps,
		Posts:         ps,
		ReadStatuses:  ps,
		EncryptionKey: encryptionKey,
	}, nil
}

// CreateUser creates a new user from an email and password
func (s *Store) CreateUser(email, password string, analytics bool) (*User, error) {
	encPass, err := abdi.Hash(password, s.EncryptionKey)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	u, err := s.Users.CreateUser(&User{
		Email:             email,
		EncryptedPassword: *encPass,
		Confirmed:         false,
		TokenCreatedAt:    now,
		Analytics:         analytics,
	})
	if err != nil {
		return nil, err
	}

	f, err := s.Folders.CreateFolder(&Folder{
		Name: "default",
	})
	if err != nil {
		return nil, err
	}

	err = s.Users.AddFolder(u.ID, f.ID)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// GetUserByToken returns the user with the given access token
func (s *Store) GetUserByToken(token string) (*User, error) {
	sesh, err := s.Sessions.GetSessionByAccessToken(token)
	if err != nil {
		return nil, err
	}

	user, err := s.Users.GetUser(sesh.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
