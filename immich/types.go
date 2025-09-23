package immich

// Config holds the configuration for the Immich client
type Config struct {
	ServerURL string
	AlbumName string
	APIKey    string
	DryRun    bool
	Verbose   bool
	Debug     bool
}

// Asset represents an Immich asset (photo/video)
type Asset struct {
	ID         string `json:"id"`
	IsFavorite bool   `json:"isFavorite"`
	Type       string `json:"type"`
	OwnerID    string `json:"ownerId"`
}

// Album represents an Immich album
type Album struct {
	ID          string `json:"id"`
	AlbumName   string `json:"albumName"`
	Description string `json:"description,omitempty"`
}

// AlbumUser represents a user in an album
type AlbumUser struct {
	UserID string `json:"userId"`
	Role   string `json:"role"`
}

// CreateAlbumRequest represents the request to create a new album
type CreateAlbumRequest struct {
	AlbumName   string      `json:"albumName"`
	AssetIDs    []string    `json:"assetIds,omitempty"`
	Description string      `json:"description,omitempty"`
	AlbumUsers  []AlbumUser `json:"albumUsers"`
}

// AddAssetsRequest represents the request to add assets to an album
type AddAssetsRequest struct {
	AssetIDs []string `json:"ids"`
}

// SearchMetadataRequest represents a search request for assets with metadata filtering
type SearchMetadataRequest struct {
	IsFavorite bool   `json:"isFavorite"`
	Type       string `json:"type,omitempty"` // "IMAGE", "VIDEO", or empty for both
	Page       int    `json:"page,omitempty"`
	Size       int    `json:"size,omitempty"`
}

// User represents a user in the system
type User struct {
	ID                string `json:"id"`
	Email             string `json:"email"`
	Name              string `json:"name"`
	AvatarColor       string `json:"avatarColor"`
	ProfileImagePath  string `json:"profileImagePath"`
	IsAdmin           bool   `json:"isAdmin"`
	ShouldChangePassword bool `json:"shouldChangePassword"`
	ProfileChangedAt  string `json:"profileChangedAt"`
}

// SearchMetadataResponse represents the response from a search metadata request
type SearchMetadataResponse struct {
	Albums AssetsContainer `json:"albums"`
	Assets AssetsContainer `json:"assets"`
}

// AssetsContainer represents a container for assets with pagination
type AssetsContainer struct {
	Total     int    `json:"total"`
	Count     int    `json:"count"`
	Items     []Asset `json:"items"`
	Facets    []interface{} `json:"facets"`
	NextPage  *string `json:"nextPage"`
}
