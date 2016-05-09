package store

import (
    "github.com/fortytw2/kiasu"
)

type UserStore interface {
    CreateUser(*kiasu.User) (*kiasu.User, error)


    GetFeeds(*kiasu.User) ([]*kiasu.Feed, error)
}

type FeedStore interface {
    CreateFeed(*kiasu.Feed) (*kiasu.Feed, error)

    MarkUpdate(*kiasu.Feed) error
    GetFeedsToUpdate() ([]*kiasu.Feed, error)

    GetUnreadArticles(*kiasu.User, *kiasu.Feed) ([]*kiasu.Article, error)
}

type ArticleStore interface {
}
