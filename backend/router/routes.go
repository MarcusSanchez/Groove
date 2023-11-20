package router

import (
	"GrooveGuru/router/handlers"
	"GrooveGuru/router/middleware"
	"github.com/gofiber/fiber/v2"
)

func Start(app *fiber.App) {
	/** api endpoints **/
	api := app.Group("/api")
	api.Get("/health", handlers.Health)

	/** session endpoints **/
	api.Post("/register", middleware.RedirectAuthorized, handlers.Register)
	api.Post("/login", middleware.RedirectAuthorized, handlers.Login)
	api.Post("/logout", middleware.CheckCSRF, handlers.Logout)
	api.Post("/authenticate", middleware.CheckCSRF, handlers.Authenticate)

	/** spotify-link endpoints **/
	spotify := api.Group("/spotify")
	spotify.Post("/link", middleware.CheckCSRF, middleware.RedirectLinked, handlers.LinkSpotify)
	spotify.Get("/callback", middleware.AuthorizeAny, handlers.SpotifyCallback)
	spotify.Post("/unlink", middleware.CheckCSRF, handlers.UnlinkSpotify)

	/** spotify-artist endpoints **/
	artists := spotify.Group("/artists")
	artists.Get("/:id", middleware.AuthorizeAny, middleware.SetAccess, handlers.GetArtist)
	artists.Get("/:id/related-artists", middleware.AuthorizeAny, middleware.SetAccess, handlers.GetRelatedArtists)
	artists.Get("/:id/top-tracks", middleware.AuthorizeAny, middleware.SetAccess, handlers.GetArtistTopTracks)
	artists.Get("/:id/albums", middleware.AuthorizeAny, middleware.SetAccess, handlers.GetArtistAlbums)

	/** spotify-album endpoints **/
	albums := spotify.Group("/albums")
	albums.Get("/:id", middleware.AuthorizeAny, middleware.SetAccess, handlers.GetAlbum)
	albums.Get("/:id/tracks", middleware.AuthorizeAny, middleware.SetAccess, handlers.GetAlbumTracks)

	/** spotify-tracks endpoints **/
	tracks := spotify.Group("/tracks")
	tracks.Get("/:id", middleware.AuthorizeAny, middleware.SetAccess, handlers.GetTrack)

	/** spotify-search endpoints **/
	search := spotify.Group("/search")
	search.Get("/:query", middleware.AuthorizeAny, middleware.SetAccess, handlers.Search)
}
