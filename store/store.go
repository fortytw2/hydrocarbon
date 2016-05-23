package store

import (
	"github.com/fortytw2/kiasu"
)

type Store interface {
	UserStore
	FeedStore
}

type UserStore interface {
	CreateUser(*kiasu.User) (*kiasu.User, error)

	GetFeeds(*kiasu.User) ([]*kiasu.Feed, error)
}

type FeedStore interface {
	CreateFeed(*kiasu.Feed) (*kiasu.Feed, error)

	MarkUpdate(*kiasu.Feed) error
	GetFeedsToUpdate() ([]*kiasu.Feed, error)

	CreateArticle(*kiasu.Article, *kiasu.Feed) error
	GetArticles(*kiasu.User, *kiasu.Feed) ([]*kiasu.Article, error)
}
