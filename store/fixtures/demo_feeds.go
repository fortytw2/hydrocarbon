package fixtures

import (
    "github.com/fortytw2/kiasu"
    "github.com/fortytw2/kiasu/store"
)

// Feeds are the default feeds loaded by Kiasu
var Feeds = []kiasu.Feed{




}

// SyncFeeds adds the boilerplate feeds to the datastore
func SyncFeeds(s store.FeedStore) error {
    return nil
}