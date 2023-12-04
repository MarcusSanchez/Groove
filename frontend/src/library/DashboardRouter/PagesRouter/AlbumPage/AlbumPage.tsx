import { useNavigate } from "react-router-dom";
import { useEffect, useState } from "react";
import { Album, Albums } from "./types.ts";
import { HashLink } from "react-router-hash-link";

const noImageURL = "https://static.thenounproject.com/png/1554489-200.png";

function AlbumPage() {

  const params = new URLSearchParams(window.location.search);
  const navigate = useNavigate();

  const [album, setAlbum] = useState<Album | null>(null);
  const [otherAlbums, setOtherAlbums] = useState<Albums | null>(null);

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
    displayAlbum().then();
  }, [])

  const swapAlbum = async (id: string) => {
    params.set("id", id);

    setOtherAlbums(null);
    setAlbum(null);

    await displayAlbum();
    window.scrollTo({ top: 0, behavior: "smooth" });
  }

  return (
    <div className="m-5 sm:p-8 p-4 rounded-lg Shadow BOBorder border-2 w-[95%] h-min">
      <div className="flex gap-4 justify-between"> {/* Album name and Image */}
        <div className="flex flex-col">
          <h1 className="text-3xl font-bold mb-1">{album?.name}</h1>
          <HashLink
            to={`/dashboard/pages/artist?id=${album?.artists[0].id}#top`}
            className="text-xl font-semibold hover:pointer hover:underline">
            {album?.artists[0].name}
          </HashLink>
          <h2 className="text-base font-semibold italic text-gray-700">{album?.total_tracks} Tracks</h2>
          <h2 className="text-base font-semibold italic text-gray-700">{album?.release_date.slice(0, 4)}</h2>
        </div>
        <img className="w-40 h-40 rounded-md border border-black" src={album?.images[0].url || noImageURL} alt="Album Art" />
      </div>

      <hr className="my-4 border-black mb-5j" />

      <div className="flex flex-col"> {/* Track List */}
        <h1 className="font-bold lg:text-4xl md:text-2xl text-xl mb-2">Track-list</h1>
        <hr className="border border-black mb-4" />
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 w-[95%]">
          {album?.tracks.items.map((track, i) => (
            <div key={i} className="flex gap-2 ">
              <img className="w-20 h-20 rounded-md border-black border" src={album?.images[0].url || noImageURL} alt="Album Image" />
              <div className="flex flex-col">
                <HashLink
                  to={`/dashboard/pages/track?id=${track.id}#top`}
                  className="font-bold text-xl hover:cursor-pointer hover:underline"
                >
                  {track.name}
                </HashLink>
                <p className="text-sm font-semibold italic text-gray-700">
                  {msToMinutesSeconds(track?.duration_ms)}
                </p>
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
              <img className="w-20 h-20 rounded-md hover:cursor-pointer border-black border" src={album.images[0].url || noImageURL} alt="Album Image" />
              <div className="flex flex-col">
                <p
                  onClick={() => swapAlbum(album.id)}
                  className="font-bold text-xl hover:cursor-pointer hover:underline"
                >
                  {album.name}
                </p>
                <p>{album.release_date.slice(0, 4)}</p>
              </div>
            </div>
          ))}
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

export default AlbumPage;