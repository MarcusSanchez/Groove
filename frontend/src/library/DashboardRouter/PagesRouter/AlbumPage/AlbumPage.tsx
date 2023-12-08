import { useNavigate, useParams } from "react-router-dom";
import { useEffect, useState } from "react";
import { Album, Albums } from "./types.ts";
import { HashLink } from "react-router-hash-link";
import { Playlist, Playlists } from "@/library/DashboardRouter/AllPlaylists/types.ts";
import { useAtom } from "jotai";
import { csrfTokenAtom, spotifyAtom, spotifyIDAtom } from "Atoms";

const noImageURL = "https://static.thenounproject.com/png/1554489-200.png";
let TrackID: string;

function AlbumPage() {
  const navigate = useNavigate();
  const params = new URLSearchParams(window.location.search);
  const hookedParams = useParams();

  const [album, setAlbum] = useState<Album | null>(null);
  const [otherAlbums, setOtherAlbums] = useState<Albums | null>(null);

  const [playlists, setPlaylists] = useState<Playlists | null>(null);
  const [modal, setModal] = useState<boolean>(false);
  const [spotify] = useAtom(spotifyAtom);
  const [spotifyID] = useAtom(spotifyIDAtom);
  const [csrf] = useAtom(csrfTokenAtom);

  const displayAlbum = async () => {
    const id = params.get("id");
    if (!id) return;

    let localData = JSON.parse(localStorage.getItem(`albumJSON-${id}`) || "{}");
    if (localData.id === id && localData.expires > Date.now()) {
      setAlbum(localData.album);
      setOtherAlbums(localData.otherAlbums);
      return;
    }

    let resp = await fetch(`/api/spotify/albums/${id}`);
    switch (resp.status) {
      case 400:
      case 404:
        navigate("/404");
        return;
      case 500:
        console.error("Internal Server Error Fetching Album");
        return;
    }
    const album = await resp.json() as Album;
    setAlbum(album);

    resp = await fetch(`/api/spotify/artists/${album?.artists[0].id}/albums`);
    const otherAlbums = await resp.json() as Albums;
    setOtherAlbums(otherAlbums);

    localStorage.setItem(`albumJSON-${id}`, JSON.stringify({
      id: id,
      album: album,
      otherAlbums: otherAlbums,
      expires: Date.now() + (1000 * 60 * 60 * 24) // 24 hours
    }));
  }

  useEffect(() => {
    displayAlbum().then(
      () => window.scrollTo({ top: 0, behavior: "smooth" })
    );
  }, [hookedParams]);

  const swapAlbum = async (id: string) => {
    // update url (required for back button since react router doesn't care for query params)
    window.history.pushState({}, "", `${window.location.pathname}?${params.toString()}`);
    params.set("id", id);
    navigate(`?${params.toString()}`, { replace: true });

    setOtherAlbums(null);
    setAlbum(null);
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
          <div className="flex gap-4 justify-between"> {/* Album name and Image */}
            <div className="flex flex-col">
              <h1 className="sm:text-3xl text-xl font-bold">{album?.name}</h1>
              <HashLink
                to={`/dashboard/pages/artist?id=${album?.artists[0].id}#top`}
                className="sm:text-xl font-semibold text- hover:pointer over:underline text-lg text-BrandBlue hover:text-blue-700">
                {album?.artists[0].name}
              </HashLink>
              <h2 className="sm:text-base text-sm font-semibold italic text-gray-700">{album?.total_tracks} Tracks</h2>
              <h2 className="sm:text-base text-sm font-semibold italic text-gray-700">{album?.release_date.slice(0, 4)}</h2>
            </div>
            <img className="w-40 h-40 rounded-md border border-black" src={album?.images[0]?.url || noImageURL} alt="Album Art" />
          </div>

          <hr className="my-4 border-black mb-5" />

          <div className="flex flex-col"> {/* Track List */}
            <h1 className="font-bold lg:text-4xl md:text-2xl text-xl mb-2">Track-list</h1>
            <hr className="border border-black mb-4" />
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 w-[95%]">
              {album?.tracks.items.map((track, i) => (
                <div key={i} className="flex gap-2 group">
                  <img className="sm:w-20 w-16 sm:h-20 h-16 rounded-md border-black border" src={album?.images[0]?.url || noImageURL} alt="Album Image" />
                  <div className="flex justify-between w-full">
                    <div className="flex flex-col">
                      <HashLink
                        to={`/dashboard/pages/track?id=${track.id}#top`}
                        className="font-bold sm:text-xl text-base hover:cursor-pointer hover:underline mb-[-5px]"
                      >
                        {track.name}
                      </HashLink>
                      <p>
                        {track?.artists.map((artist, i) => (
                          <HashLink
                            key={i}
                            to={`/dashboard/pages/artist?id=${artist.id}#top`}
                            className="hover:underline text-sm font-semibold italic text-gray-700"
                          >
                            {artist.name}{i < track?.artists.length - 1 ? ", " : ""}
                          </HashLink>
                        ))}
                      </p>
                      <p className="text-sm font-semibold italic text-gray-700">
                        {msToMinutesSeconds(track?.duration_ms)}
                      </p>
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
                </div>
              ))}
            </div>
          </div>

          <div className="flex flex-col mt-10"> {/* Other Albums */}
            <h1 className="font-bold lg:text-4xl md:text-2xl text-xl mb-2">Other Albums</h1>
            <hr className="border border-black mb-4" />
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 w-[95%]">
              {otherAlbums?.items.map((album, i) => (
                album.id !== params.get("id") &&
                <div key={i} className="flex gap-2 ">
                  <img className="sm:w-20 w-16 sm:h-20 h-16 rounded-md hover:cursor-pointer border-black border" src={album.images[0]?.url || noImageURL} alt="Album Image" />
                  <div className="flex flex-col">
                    <p
                      onClick={() => swapAlbum(album.id)}
                      className="font-bold sm:text-xl text-base hover:cursor-pointer hover:underline"
                    >
                      {album.name}
                    </p>
                    <p className="text-sm font-semibold italic text-gray-700">{album.release_date.slice(0, 4)}</p>
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

function msToMinutesSeconds(durationMs: number | undefined): string {
  if (!durationMs) return "0:00";
  const durationSec: number = Math.floor(durationMs / 1000);
  const minutes: number = Math.floor(durationSec / 60);
  const seconds: number = durationSec % 60;

  // Ensure seconds are displayed with leading zero if needed
  return `${minutes}:${seconds.toString().padStart(2, '0')}`;
}

export default AlbumPage;