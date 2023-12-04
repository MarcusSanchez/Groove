import { useNavigate } from "react-router-dom";
import React, { useEffect, useState } from "react";
import { Track } from "./types.ts";
import { HashLink } from "react-router-hash-link";

const noImageURL = "https://static.thenounproject.com/png/1554489-200.png";

function TrackPage() {
  const params = new URLSearchParams(window.location.search);
  const navigate = useNavigate();

  const [track, setTrack] = useState<Track | null>(null);

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

  return (
    <div className="m-5 sm:p-8 p-4 rounded-lg Shadow BOBorder border-2 w-[95%] h-min">
      <div className="flex gap-4"> {/* Track Name and Image */}
        <img src={track?.album.images[0].url ?? noImageURL} alt="Album Art" className="h-32 w-32 border border-black rounded-lg" />
        <div className="flex flex-col">
          <h1 className="text-2xl font-bold">{track?.name}</h1>
          <h2 className="text-lg">
            <HashLink
              to={`/dashboard/pages/artist?id=${track?.artists[0].id}#top`}
              className="text-BrandBlue hover:text-blue-700"
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
          </h2>
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