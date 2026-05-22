package constants

const (
	SongStatusListed    = 1
	SongStatusUnlisted  = 0
	SongStatusReviewing = 2

	PlaylistStatusPublic  = 1
	PlaylistStatusPrivate = 0
	PlaylistStatusDeleted = 2
)

const (
	SongsKeyExpire    = 72 * ONE_HOUR
	PlaylistKeyExpire = 72 * ONE_HOUR
)
