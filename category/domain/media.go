package domain

type (
	// Media provides the category media interface
	Media interface {
		Type() string
		MimeType() string
		Usage() string
		Title() string
		Reference() string
	}

	// Medias defines the category media slice
	Medias []Media

	// MediaData defines the default domain category media data model
	MediaData struct {
		MediaType      string
		MediaMimeType  string
		MediaTitle     string
		MediaReference string
		MediaUsage     string
	}
)

// Media usage constants
const (
	MediaUsageTeaser = "teaser"
	MediaUsageDetail = "detail"
)

var _ Media = (*MediaData)(nil)

// MimeType gets the media mime type
func (m MediaData) MimeType() string {
	return m.MediaMimeType
}

// Title gets the media title
func (m MediaData) Title() string {
	return m.MediaTitle
}

// Reference gets the media reference
func (m MediaData) Reference() string {
	return m.MediaReference
}

// Type gets the media type
func (m MediaData) Type() string {
	return m.MediaType
}

// Usage gets the media usage
func (m MediaData) Usage() string {
	return m.MediaUsage
}

// Has checks if the Medias contain a media with `usage`
func (m Medias) Has(usage string) bool {
	for _, media := range m {
		if media.Usage() == usage {
			return true
		}
	}

	return false
}

// Get returns a *Media with `usage`, empty media result if none was found
func (m Medias) Get(usage string) Media {
	for _, media := range m {
		if media.Usage() == usage {
			return media
		}
	}

	return MediaData{}
}
