package actions

import (
	"github.com/gofiber/fiber/v2"
	. "groove/pkgs/util"
)

// GetArtist returns the artist with the given id.
func (*Actions) GetArtist(c *fiber.Ctx, artistID string) error {
	proxy := &Proxy{
		Endpoint: "/artists/" + artistID,
		Access:   c.Locals("access").(string),
	}
	return proxy.ArtistRequest(c)
}

// GetRelatedArtists returns the artists related to the artist with the given id.
func (*Actions) GetRelatedArtists(c *fiber.Ctx, artistID string) error {
	proxy := &Proxy{
		Endpoint: "/artists/" + artistID + "/related-artists",
		Access:   c.Locals("access").(string),
	}
	return proxy.ArtistRequest(c)
}

// GetArtistTopTracks returns the top tracks of the artist with the given id.
func (*Actions) GetArtistTopTracks(c *fiber.Ctx, artistID string) error {
	proxy := &Proxy{
		Endpoint: "/artists/" + artistID + "/top-tracks?market=US",
		Access:   c.Locals("access").(string),
	}
	return proxy.ArtistRequest(c)
}

// GetArtistAlbums returns the albums of the artist with the given id.
func (*Actions) GetArtistAlbums(c *fiber.Ctx, artistID string) error {
	proxy := &Proxy{
		Endpoint: "/artists/" + artistID + "/albums?market=US&limit=50&include_groups=album",
		Access:   c.Locals("access").(string),
	}
	return proxy.ArtistRequest(c)
}
