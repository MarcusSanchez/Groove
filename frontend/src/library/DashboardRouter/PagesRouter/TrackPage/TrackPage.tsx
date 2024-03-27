import { useNavigate } from "react-router-dom";
import React, { useEffect, useState } from "react";
import { Track } from "./types.ts";
import { HashLink } from "react-router-hash-link";
import { csrfTokenAtom, spotifyAtom, spotifyIDAtom } from "Atoms";
import { useAtom } from "jotai";
import { Playlist, Playlists } from "@/library/DashboardRouter/AllPlaylists/types.ts";

const noImageURL = "https://static.thenounproject.com/png/1554489-200.png";
let TrackID: string;

function TrackPage() {
  const params = new URLSearchParams(window.location.search);
  const navigate = useNavigate();

  const [track, setTrack] = useState<Track | null>(null);

  const [playlists, setPlaylists] = useState<Playlists | null>(null);
  const [modal, setModal] = useState<boolean>(false);
  const [spotify] = useAtom(spotifyAtom);
  const [spotifyID] = useAtom(spotifyIDAtom);
  const [csrf] = useAtom(csrfTokenAtom);

  useEffect(() => {
    (async () => {
      const id = params.get("id");
      if (!id) return;

      const localData = JSON.parse(localStorage.getItem(`trackJSON-${id}`) || "{}");
      if (localData.id === id && localData.expires > Date.now()) {
        setTrack(localData.track);
        return;
      }

      let resp = await fetch(`/api/spotify/tracks/${id}`);
      switch (resp.status) {
        case 400:
        case 404:
          navigate("/404");
          return;
        case 500:
          console.error("Internal Server Error Fetching Track");
          return;
      }
      const track = await resp.json() as Track;
      setTrack(track);

      localStorage.setItem(`trackJSON-${id}`, JSON.stringify({
        id: id,
        track: track,
        expires: Date.now() + (1000 * 60 * 60 * 24) // 24 hours
      }));
    })()
  }, [])

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
          <div className="flex gap-4"> {/* Track Name and Image */}
            <img src={track?.album.images[0].url ?? noImageURL} alt="Album Art" className="sm:h-32 sm:w-32 h-28 w-28 border border-black rounded-lg" />
            <div className="flex flex-col justify-between">
              <div>
                <h1 className="sm:text-2xl text-lg font-bold">{track?.name}</h1>
                <div className="sm:text-lg text-sm">
                  <HashLink
                    to={`/dashboard/pages/artist?id=${track?.artists[0].id}#top`}
                    className="text-BrandBlue hover:text-blue-700 sm:text-lg text-base"
                  >
                    {track?.artists[0].name}
                  </HashLink>
                  <br />
                  {track?.artists.length && track?.artists.length > 1 &&
                    <span className="font-semibold italic">Featuring: </span>
                  }
                  {track?.artists.slice(1).map((artist, i) => (
                    <HashLink
                      key={i}
                      to={`/dashboard/pages/artist?id=${artist.id}#top`}
                      className="hover:underline"
                    >
                      {artist.name}{i < track?.artists.length - 2 ? ", " : ""}
                    </HashLink>
                  ))}
                </div>
              </div>
              {spotify &&
                <button
                  onClick={() => displayModal(track!.id)}
                  className="bg-white text-black px-2 py-1 rounded-md BOBorder font-bold border-2 sm:text-sm w-[150px]
                drop-shadow-md hover:outline-none hover:ring-2 focus:border-black hover:ring-opacity-50 text-xs"
                >
                  Add to Playlist <i className="fas fa-plus"></i>
                </button>
              }
            </div>
          </div>

          <hr className="my-4 border-black" />
          {/* TRACK DETAILS*/}
          <div className="flex flex-col gap-2">
            <h1 className="font-bold text-xl">Track Details</h1>
            <div className="flex flex-col gap-1">
              <p className="font-semibold">
                <span className="font-semibold italic">Album: </span>
                <HashLink
                  to={`/dashboard/pages/album?id=${track?.album.id}#top`}
                  className="hover:underline hover:cursor-pointer hover:text-blue-700 text-BrandBlue"
                >
                  {track?.album.name}
                </HashLink></p>
              <p>
                <span className="font-semibold italic">Duration: </span>
                {msToMinutesSeconds(track?.duration_ms)}
              </p>
              <p>
                <span className="font-semibold italic">Release Date: </span>
                {track?.album.release_date}
              </p>
              <p>
                <span className="font-semibold italic">Explicit: </span>
                {track?.explicit ? "Yes" : "No"}
              </p>
              <p>
                <span className="font-semibold italic">Track Number: </span>
                {track?.track_number}
              </p>
              <p>
                <span className="font-semibold italic">Popularity: </span>
                {track?.popularity}
              </p>
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
                          src={playlist.images?.[0]?.url || noImageURL}
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
      </div>
    </>
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

export default TrackPage;