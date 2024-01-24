package actions

import (
	"github.com/gofiber/fiber/v2"
	. "groove/pkgs/util"
)

// GetAlbum returns the album with the given id.
func (*Actions) GetAlbum(c *fiber.Ctx, albumID string) error {
	proxy := &Proxy{
		Endpoint: "/albums/" + albumID,
		Access:   c.Locals("access").(string),
	}
	return proxy.AlbumRequest(c)
}

// GetAlbumTracks returns the tracks of the album with the given id.
func (*Actions) GetAlbumTracks(c *fiber.Ctx, albumID string) error {
	proxy := &Proxy{
		Endpoint: "/albums/" + albumID + "/tracks?limit=50&market=US",
		Access:   c.Locals("access").(string),
	}
	return proxy.AlbumRequest(c)
}
