package server

import (
	"GrooveGuru/server/handlers"
	"GrooveGuru/server/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupEndpoints(app *fiber.App, handlers *handlers.Handlers, mw *middleware.Middlewares) {
	/** api endpoints **/
	api := app.Group("/api")
	api.Get("/health", handlers.Health)

	/** session endpoints **/
	api.Post("/register", mw.RedirectAuthorized, handlers.Register)
	api.Post("/login", mw.RedirectAuthorized, handlers.Login)
	api.Post("/logout", mw.CheckCSRF, handlers.Logout)
	api.Post("/authenticate", mw.CheckCSRF, handlers.Authenticate)

	/** spotify-link endpoints **/
	spotify := api.Group("/spotify")
	spotify.Post("/link", mw.CheckCSRF, mw.RedirectLinked, handlers.LinkSpotify)
	spotify.Get("/callback", mw.AuthorizeAny, handlers.SpotifyCallback)
	spotify.Post("/unlink", mw.CheckCSRF, handlers.UnlinkSpotify)
	spotify.Get("/me", mw.AuthorizeLinked, mw.SetAccess, handlers.GetCurrentUser)

	/** spotify-artist endpoints **/
	artists := spotify.Group("/artists")
	artists.Get("/:id", mw.AuthorizeAny, mw.SetAccess, handlers.GetArtist)
	artists.Get("/:id/related-artists", mw.AuthorizeAny, mw.SetAccess, handlers.GetRelatedArtists)
	artists.Get("/:id/top-tracks", mw.AuthorizeAny, mw.SetAccess, handlers.GetArtistTopTracks)
	artists.Get("/:id/albums", mw.AuthorizeAny, mw.SetAccess, handlers.GetArtistAlbums)

	/** spotify-album endpoints **/
	albums := spotify.Group("/albums")
	albums.Get("/:id", mw.AuthorizeAny, mw.SetAccess, handlers.GetAlbum)
	albums.Get("/:id/tracks", mw.AuthorizeAny, mw.SetAccess, handlers.GetAlbumTracks)

	/** spotify-tracks endpoints **/
	tracks := spotify.Group("/tracks")
	tracks.Get("/:id", mw.AuthorizeAny, mw.SetAccess, handlers.GetTrack)

	/** spotify-playlist endpoints **/
	playlists := spotify.Group("/playlists")
	playlists.Get("/", mw.AuthorizeLinked, mw.SetAccess, handlers.GetAllPlaylists)
	playlists.Get("/:id", mw.AuthorizeLinked, mw.SetAccess, handlers.GetPlaylistWithTracks)
	playlists.Get("/:id/load-more", mw.AuthorizeLinked, mw.SetAccess, handlers.GetMorePlaylistTracks)
	playlists.Post("/:id/track", mw.CheckCSRF, mw.AuthorizeLinked, mw.SetAccess, handlers.AddTrackToPlaylist)
	playlists.Delete("/:id/track", mw.CheckCSRF, mw.AuthorizeLinked, mw.SetAccess, handlers.RemoveTrackFromPlaylist)

	/** spotify-search endpoints **/
	search := spotify.Group("/search")
	search.Get("/:query", mw.AuthorizeAny, mw.SetAccess, handlers.Search)
}
