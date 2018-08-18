// Copyright (c) 2013-2018 The Gorilla Feeds Authors. All rights reserved.

// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:

//   Redistributions of source code must retain the above copyright notice, this
//   list of conditions and the following disclaimer.

//   Redistributions in binary form must reproduce the above copyright notice,
//   this list of conditions and the following disclaimer in the documentation
//   and/or other materials provided with the distribution.

// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
package jsonfeed

import (
	"encoding/json"
	"time"
)

// JSONAuthor represents the author of the feed or of an individual item
// in the feed
type JSONAuthor struct {
	Name   string `json:"name,omitempty"`
	URL    string `json:"url,omitempty"`
	Avatar string `json:"avatar,omitempty"`
}

// JSONAttachment represents a related resource. Podcasts, for instance, would
// include an attachment that’s an audio or video file.
type JSONAttachment struct {
	URL      string        `json:"url,omitempty"`
	MIMEType string        `json:"mime_type,omitempty"`
	Title    string        `json:"title,omitempty"`
	Size     int32         `json:"size,omitempty"`
	Duration time.Duration `json:"duration_in_seconds,omitempty"`
}

// MarshalJSON implements the json.Marshaler interface.
// The Duration field is marshaled in seconds, all other fields are marshaled
// based upon the definitions in struct tags.
func (a *JSONAttachment) MarshalJSON() ([]byte, error) {
	type EmbeddedJSONAttachment JSONAttachment
	return json.Marshal(&struct {
		Duration float64 `json:"duration_in_seconds,omitempty"`
		*EmbeddedJSONAttachment
	}{
		EmbeddedJSONAttachment: (*EmbeddedJSONAttachment)(a),
		Duration:               a.Duration.Seconds(),
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The Duration field is expected to be in seconds, all other field types
// match the struct definition.
func (a *JSONAttachment) UnmarshalJSON(data []byte) error {
	type EmbeddedJSONAttachment JSONAttachment
	var raw struct {
		Duration float64 `json:"duration_in_seconds,omitempty"`
		*EmbeddedJSONAttachment
	}
	raw.EmbeddedJSONAttachment = (*EmbeddedJSONAttachment)(a)

	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	if raw.Duration > 0 {
		nsec := int64(raw.Duration * float64(time.Second))
		raw.EmbeddedJSONAttachment.Duration = time.Duration(nsec)
	}

	return nil
}

// JSONItem represents a single entry/post for the feed.
type JSONItem struct {
	ID            string           `json:"id"`
	URL           string           `json:"url,omitempty"`
	ExternalURL   string           `json:"external_url,omitempty"`
	Title         string           `json:"title,omitempty"`
	ContentHTML   string           `json:"content_html,omitempty"`
	ContentText   string           `json:"content_text,omitempty"`
	Summary       string           `json:"summary,omitempty"`
	Image         string           `json:"image,omitempty"`
	BannerImage   string           `json:"banner_,omitempty"`
	PublishedDate *time.Time       `json:"date_published,omitempty"`
	ModifiedDate  *time.Time       `json:"date_modified,omitempty"`
	Author        *JSONAuthor      `json:"author,omitempty"`
	Tags          []string         `json:"tags,omitempty"`
	Attachments   []JSONAttachment `json:"attachments,omitempty"`
}

// JSONHub describes an endpoint that can be used to subscribe to real-time
// notifications from the publisher of this feed.
type JSONHub struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

// JSONFeed represents a syndication feed in the JSON Feed Version 1 format.
// Matching the specification found here: https://jsonfeed.org/version/1.
type JSONFeed struct {
	Version     string      `json:"version"`
	Title       string      `json:"title"`
	HomePageURL string      `json:"home_page_url,omitempty"`
	FeedURL     string      `json:"feed_url,omitempty"`
	Description string      `json:"description,omitempty"`
	UserComment string      `json:"user_comment,omitempty"`
	NextURL     string      `json:"next_url,omitempty"`
	Icon        string      `json:"icon,omitempty"`
	Favicon     string      `json:"favicon,omitempty"`
	Author      *JSONAuthor `json:"author,omitempty"`
	Expired     *bool       `json:"expired,omitempty"`
	Hubs        []*JSONItem `json:"hubs,omitempty"`
	Items       []*JSONItem `json:"items,omitempty"`
}
