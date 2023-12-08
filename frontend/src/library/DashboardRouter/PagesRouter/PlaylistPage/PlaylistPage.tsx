import { useNavigate } from "react-router-dom";
import React, { useEffect, useState } from "react";
import { Playlist, PlaylistTrack } from "./types.ts";
import { HashLink } from "react-router-hash-link";
import { csrfTokenAtom, loadedAtom, spotifyAtom } from "Atoms";
import { useAtom } from "jotai";

const noImageURL = "https://static.thenounproject.com/png/1554489-200.png";

function PlaylistPage() {
  const navigate = useNavigate();
  const params = new URLSearchParams(window.location.search);

  const [spotify] = useAtom(spotifyAtom);
  const [isLoaded] = useAtom(loadedAtom)
  const [csrf] = useAtom(csrfTokenAtom);

  const [tracks, setTracks] = useState<PlaylistTrack[] | null>(null);
  const [playlist, setPlaylist] = useState<Playlist | null>(null);
  const [playlistLength, setPlaylistLength] = useState<number>(0);

  if (isLoaded && !spotify) navigate("/dashboard/profile");

  useEffect(() => {
    (async () => {
      const id = params.get("id");
      if (!id) return;

      let resp = await fetch(`/api/spotify/playlists/${id}`);
      switch (resp.status) {
        case 400:
        case 404:
          navigate("/dashboard/404");
          return;
        case 500:
          console.error("Internal Server Error Fetching Playlist");
          return;
      }

      const playlist = await resp.json() as Playlist;
      setPlaylist(playlist);
      setPlaylistLength(playlist.tracks.total)
      setTracks(playlist.tracks.items)
    })()
  }, [])

  const loadMore = async () => {
    const id = params.get("id");
    if (!id) return;

    let resp = await fetch(`/api/spotify/playlists/${id}/load-more?offset=${tracks?.length}`);
    switch (resp.status) {
      case 400:
      case 404:
        navigate("/dashboard/404");
        return;
      case 500:
        console.error("Internal Server Error Fetching Playlist");
        return;
    }

    const data = await resp.json();
    const extendedTracks = [...tracks!, ...data.items];
    setTracks(extendedTracks);
  }

  const removeTrack = async (trackID: string) => {
    const id = params.get("id");
    if (!id) return;

    let resp = await fetch(`/api/spotify/playlists/${id}/track?id=${trackID}`, {
      method: "DELETE",
      credentials: "include",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        csrf_: csrf
      })
    });

    switch (resp.status) {
      case 500:
        console.error("Internal Server Error Removing Track");
        return;
      case 200:
        break;
      default:
        console.error("Error Removing Track: ", resp.status, await resp.text());
        return;
    }


    // remove this and replace with splice and findIndex
    const filteredTracks = tracks!.filter(track => track.track.id !== trackID);
    setTracks(filteredTracks);
    setPlaylistLength(filteredTracks.length);
    setPlaylist({
      ...playlist!,
      tracks: {
        ...playlist!.tracks,
        items: filteredTracks
      }
    });
  }

  return (
    <div className="w-full flex justify-center content-center">
      <div className="m-5 sm:p-8 p-4 rounded-lg Shadow BOBorder border-2 w-[95%] h-min">
        <div className="flex gap-4 justify-between"> {/* Album name and Image */}
          <div className="flex flex-col">
            <h1 className="sm:text-3xl text-xl font-bold mb-1">{playlist?.name}</h1>
            <h2 className="text-base font-semibold italic text-gray-700">{playlistLength} Tracks</h2>
          </div>
          <img className="sm:w-40 sm:h-40 w-28 h-28 rounded-md border border-black" src={playlist?.images[0]?.url || noImageURL} alt="Playlist Cover" />
        </div>

        <hr className="my-4 border-black mb-5" />

        <div className="flex flex-col"> {/* Track List */}
          <h1 className="font-bold lg:text-4xl md:text-2xl text-xl mb-2">Playlist Tracks</h1>
          <hr className="border border-black mb-4" />
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 w-[95%]">
            {tracks?.map((track, i) => (
              <div key={i} className="flex gap-2 group">
                <img className="sm:w-20 w-16 sm:h-20 h-16 rounded-md border-black border" src={track.track.album.images[0]?.url || playlist?.images[0]?.url || noImageURL} alt="Album Image" />
                <div className="flex justify-between w-full">
                  <div className="flex flex-col">
                    <HashLink
                      to={`/dashboard/pages/track?id=${track.track.id}#top`}
                      className="font-bold sm:text-xl text-lg hover:cursor-pointer hover:underline mb-[-7px]"
                    >
                      {track.track.name}
                    </HashLink>
                    <p>
                      {track.track.artists.map((artist, i) => (
                        <HashLink
                          key={i}
                          to={`/dashboard/pages/artist?id=${artist.id}#top`}
                          className="hover:underline text-sm font-semibold italic text-gray-700"
                        >
                          {artist.name}{i < track.track.artists.length - 1 ? ", " : ""}
                        </HashLink>
                      ))}
                    </p>
                    <p className="sm:text-sm text-xs font-semibold italic text-gray-700">
                      {msToMinutesSeconds(track.track.duration_ms)}
                    </p>
                  </div>
                  <button
                    onClick={() => removeTrack(track.track.id)}
                    className="bg-white text-black px-2 py-1 rounded-md BOBorder font-bold border-2 hidden group-hover:block
                  drop-shadow-md h-[35px] w-[35px] self-center hover:outline-none hover:ring-2 focus:border-black hover:ring-opacity-50"
                  >
                    X
                  </button>
                </div>
              </div>
            ))}
          </div>
          {tracks?.length !== playlistLength &&
            <button
              onClick={loadMore}
              className="bg-white text-black px-2 py-1 rounded-xl BOBorder font-bold border-2 drop-shadow-md
          hover:border- focus:outline-none hover:ring-2 focus:border-black hover:ring-opacity-50 mt-5"
            >
              Load More
            </button>
          }
        </div>
      </div>
    </div>
  );
}

function msToMinutesSeconds(durationMs: number | undefined): string {
  if (!durationMs) return "0:00";
  const durationSec: number = Math.floor(durationMs / 1000);
  const minutes: number = Math.floor(durationSec / 60);
  const seconds: number = durationSec % 60;

  // Ensure seconds are displayed with leading zero if needed
  return `${minutes}:${seconds.toString().padStart(2, '0')}`;
}

export default PlaylistPage;