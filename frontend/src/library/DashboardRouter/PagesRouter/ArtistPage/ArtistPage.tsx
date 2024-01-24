import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { Albums, Artist, RelatedArtists, TopTracks } from "./types.ts";
import { HashLink } from "react-router-hash-link";
import { csrfTokenAtom, spotifyAtom, spotifyIDAtom } from "Atoms";
import { useAtom } from "jotai";
import { Playlist, Playlists } from "@/library/DashboardRouter/AllPlaylists/types.ts";

const noImageURL = "https://static.thenounproject.com/png/1554489-200.png";
let TrackID: string;

function ArtistPage() {
  const navigate = useNavigate();
  const params = new URLSearchParams(window.location.search);
  const hookedParams = useParams();

  const [artist, setArtist] = useState<Artist | null>(null);
  const [albums, setAlbums] = useState<Albums | null>(null);
  const [topTracks, setTopTracks] = useState<TopTracks | null>(null);
  const [relatedArtists, setRelatedArtists] = useState<RelatedArtists | null>(null);

  const [playlists, setPlaylists] = useState<Playlists | null>(null);
  const [modal, setModal] = useState<boolean>(false);
  const [spotify] = useAtom(spotifyAtom);
  const [spotifyID] = useAtom(spotifyIDAtom);
  const [csrf] = useAtom(csrfTokenAtom);

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
        navigate("/dashboard/404");
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
    displayArtist().then(
      () => window.scrollTo({ top: 0, behavior: "smooth" })
    );
  }, [hookedParams]);

  const swapArtist = async (id: string) => {
    // update url (required for back button since react server doesn't care for query params)
    window.history.pushState({}, "", `${window.location.pathname}?${params.toString()}`);
    params.set("id", id);
    navigate(`?${params.toString()}`, { replace: true });

    setAlbums(null);
    setTopTracks(null);
    setRelatedArtists(null);
    setArtist(null);
  }

  const displayModal = async (trackID: string) => {
    if (!playlists) {
      const resp = await fetch("/api/spotify/playlists");
      if (resp.status === 500) {
        console.error("Internal Server Error Fetching Playlists");
        return;
      }

      const playlists = await resp.json() as Playlists;
      setPlaylists(playlists);
    }

    TrackID = trackID;
    setModal(true);
  }

  const addTrackToPlaylist = async (playlist: Playlist) => {
    const resp = await fetch(`/api/spotify/playlists/${playlist.id}/track?id=${TrackID}`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        csrf_: csrf
      })
    });
    if (resp.status !== 201) {
      console.error("Error Adding Track To Playlist:", resp.status, resp.statusText, await resp.text());
      alert("Error Adding Track To Playlist");
      setModal(false);
      return;
    }

    for (let i = 0; i < playlists!.items.length; i++) {
      if (playlists!.items[i].id === playlist.id) {
        playlists!.items[i].tracks.total++;
        break;
      }
    }
    setPlaylists(playlists);
    setModal(false);
  }

  return (
    <>
      <div className="w-full flex justify-center content-center">
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

            <img className="lg:w-[25%] w-[40%] h-auto rounded-md border-black border" src={artist?.images[0]?.url || noImageURL} alt="Artist Image" />
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
                <div key={i} className="flex gap-2 group">
                  <img className="sm:w-20 w-16 sm:h-20 h-16 rounded-md hover:cursor-pointer border-black border" src={track.album.images[0]?.url || noImageURL} alt="Album Image" />
                  <div className="flex justify-between w-full">
                    <div className="flex flex-col">
                      <HashLink
                        to={`/dashboard/pages/track?id=${track.id}#top`}
                        className="font-bold sm:text-xl text-base hover:cursor-pointer hover:underline"
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
                  {spotify &&
                    <button
                      onClick={() => displayModal(track.id)}
                      className="bg-white text-black px-2 py-1 rounded-md BOBorder font-bold border-2 hidden group-hover:block
                    drop-shadow-md h-[35px] w-[35px] self-center hover:outline-none hover:ring-2 focus:border-black hover:ring-opacity-50"
                    >
                      <i className="fas fa-plus"></i>
                    </button>
                  }
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
                  <img className="sm:w-20 sm:h-20 w-16 h-16 rounded-md hover:cursor-pointer border-black border" src={album.images[0]?.url || noImageURL} alt="Album Image" />
                  <div className="flex flex-col">
                    <HashLink
                      to={`/dashboard/pages/album?id=${album.id}#top`}
                      className="font-bold sm:text-xl text-base hover:cursor-pointer hover:underline"
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
      </div>

      <div className={modal ? "block" : "hidden"}> {/* Modal */}
        <div className="fixed z-10 inset-0 overflow-y-auto">
          <div className="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
            {/* Background overlay, show/hide based on modal state. */}
            <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" aria-hidden="true"></div>

            {/* This element is to trick the browser into centering the modal contents. */}
            <span className="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>

            {/* Modal panel, show/hide based on modal state. */}
            <div className="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
              <div className="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
                {/* display all playlists */}
                <button onClick={() => setModal(false)}>
                  <i className="fas fa-times float-right hover:cursor-pointer"></i>
                </button>
                <br />

                <div className="grid grid-cols-2 gap-4 w-[95%]"> {/* Playlists */}
                  {playlists?.items.map((playlist, i) => (
                    playlist.owner.id === spotifyID &&
                    <div key={i} className="flex gap-2 ">
                      <img
                        onClick={() => addTrackToPlaylist(playlist)}
                        className="w-14 h-14 rounded-md hover:cursor-pointer border-black border"
                        src={playlist.images[0]?.url || noImageURL}
                        alt="playlist image"
                      />
                      <div className="flex flex-col">
                        <p
                          onClick={() => addTrackToPlaylist(playlist)}
                          className="font-bold text-base hover:cursor-pointer hover:underline"
                        >
                          {playlist.name}
                        </p>
                        <p className="text-sm font-semibold italic text-gray-700">
                          {playlist.tracks.total} Tracks
                        </p>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}

function commaFormat(number: number | undefined) {
  if (!number) return "0";
  let numStr = number.toString();
  numStr = numStr.replace(/\B(?=(\d{3})+(?!\d))/g, ",");
  return numStr;
}

export default ArtistPage;