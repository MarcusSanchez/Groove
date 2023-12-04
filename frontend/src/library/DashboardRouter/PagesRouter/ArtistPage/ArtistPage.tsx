import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { Albums, Artist, RelatedArtists, TopTracks } from "./types.ts";
import { HashLink } from "react-router-hash-link";

const noImageURL = "https://static.thenounproject.com/png/1554489-200.png";

function ArtistPage() {
  const params = new URLSearchParams(window.location.search);
  const navigate = useNavigate();

  const [artist, setArtist] = useState<Artist | null>(null);
  const [albums, setAlbums] = useState<Albums | null>(null);
  const [topTracks, setTopTracks] = useState<TopTracks | null>(null);
  const [relatedArtists, setRelatedArtists] = useState<RelatedArtists | null>(null);

  const displayArtist = async () => {
    const id = params.get("id");
    if (!id) return;

    // retrieved cached data
    let localData = JSON.parse(localStorage.getItem(`artistJSON-${id}`) || "{}");
    if (localData.id === id && localData.expires > Date.now()) {
      setArtist(localData.artist);
      setAlbums(localData.albums);
      setTopTracks(localData.topTracks);
      setRelatedArtists(localData.relatedArtists);
      return;
    }

    let resp = await fetch(`/api/spotify/artists/${id}`);
    switch (resp.status) {
      case 400:
      case 404:
        navigate("/404");
        return;
      case 500:
        console.error("Internal Server Error Fetching Artist");
        return;
    }
    const artist = await resp.json() as Artist;
    setArtist(artist);

    resp = await fetch(`/api/spotify/artists/${id}/related-artists`);
    const relatedArtists = await resp.json() as RelatedArtists;
    setRelatedArtists(relatedArtists);

    resp = await fetch(`/api/spotify/artists/${id}/top-tracks`);
    const topTracks = await resp.json() as TopTracks;
    setTopTracks(topTracks);

    resp = await fetch(`/api/spotify/artists/${id}/albums`);
    const albums = await resp.json() as Albums;
    setAlbums(albums);

    // cache data to not spam spotify api
    localStorage.setItem(`artistJSON-${id}`, JSON.stringify({
      id: id,
      artist: artist,
      albums: albums,
      topTracks: topTracks,
      relatedArtists: relatedArtists,
      expires: Date.now() + (1000 * 60 * 20) // 1 hour
    }))
  }

  useEffect(() => {
    displayArtist().then();
  }, []);

  const swapArtist = async (id: string) => {
    params.set("id", id);

    setAlbums(null);
    setTopTracks(null);
    setRelatedArtists(null);
    setArtist(null);

    await displayArtist();
    window.scrollTo({ top: 0, behavior: "smooth" });
  }

  return (
    <div className="m-5 sm:p-8 p-4 rounded-lg Shadow BOBorder border-2 w-[95%] h-min">
      <div className="flex gap-4 justify-between"> {/* Artist Name and Image */}
        <div className="flex flex-col w-[50%]">
          <h1 className="font-bold lg:text-4xl md:text-2xl text-xl mb-2">{artist?.name}</h1>
          <p className="mb-1 md:text-base text-sm">
            <span className="text-base font-semibold italic text-gray-700">Genres: </span>
            {artist?.genres.join(", ")}
          </p>
          <p className="mb-1 md:text-base text-sm">
            <span className="font-semibold italic text-gray-700">Followers: </span>
            {commaFormat(artist?.followers.total)}
          </p>
          <p className="mb-1 md:text-base text-sm">
            <span className="font-semibold italic text-gray-700">Popularity: </span>
            {artist?.popularity}
          </p>
        </div>

        <img className="lg:w-[25%] w-[40%] h-min rounded-md border-black border" src={artist?.images[0].url || noImageURL} alt="Artist Image" />
      </div>

      <hr className="my-4 border-black" />

      <div className="flex flex-col mb-4"> {/* Related Artists */}
        <h1 className="font-bold text-xl mb-1">Related Artists: </h1>
        <div className="flex flex-wrap gap-2">
          {relatedArtists?.artists.slice(0, 5).map((artist, i) => (
            <p
              key={i}
              onClick={() => swapArtist(artist.id)}
              className="underline hover:cursor-pointer md:text-base text-sm"
            >
              {artist.name}
            </p>
          ))}
        </div>
      </div>

      <div className="flex flex-col"> {/* Top Tracks */}
        <h1 className="font-bold lg:text-4xl md:text-2xl text-xl mb-2">Top Tracks</h1>
        <hr className="border border-black mb-4" />
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 w-[95%]">
          {topTracks?.tracks.map((track, i) => (
            <div key={i} className="flex gap-2 ">
              <img className="w-20 h-20 rounded-md hover:cursor-pointer border-black border" src={track.album.images[0].url || noImageURL} alt="Album Image" />
              <div className="flex flex-col">
                <HashLink
                  to={`/dashboard/pages/track?id=${track.id}#top`}
                  className="font-bold text-xl hover:cursor-pointer hover:underline"
                >
                  {track.name}
                </HashLink>
                <HashLink
                  to={`/dashboard/pages/album?id=${track.album.id}#top`}
                  className="hover:underline hover:cursor-pointer text-sm font-semibold italic text-gray-700"
                >
                  {track.album.name}
                </HashLink>
              </div>
            </div>
          ))}
        </div>
      </div>

      <div className="flex flex-col mt-10"> {/* Albums */}
        <h1 className="font-bold lg:text-4xl md:text-2xl text-xl mb-2">Albums</h1>
        <hr className="border border-black mb-4" />
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 w-[95%]">
          {albums?.items.map((album, i) => (
            <div key={i} className="flex gap-2 ">
              <img className="w-20 h-20 rounded-md hover:cursor-pointer border-black border" src={album.images[0].url || noImageURL} alt="Album Image" />
              <div className="flex flex-col">
                <HashLink
                  to={`/dashboard/pages/album?id=${album.id}`}
                  className="font-bold text-xl hover:cursor-pointer hover:underline"
                >
                  {album.name}
                </HashLink>
                <p className="text-sm font-semibold italic text-gray-700">
                  {album.release_date.slice(0, 4)}
                </p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

function commaFormat(number: number | undefined) {
  if (!number) return "0";
  let numStr = number.toString();
  numStr = numStr.replace(/\B(?=(\d{3})+(?!\d))/g, ",");
  return numStr;
}

export default ArtistPage;