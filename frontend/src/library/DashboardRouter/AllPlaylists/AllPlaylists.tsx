import React, { useEffect, useState } from "react";
import { Playlists } from "./types.ts";
import { useAtom } from "jotai";
import { spotifyAtom, spotifyIDAtom } from "Atoms";
import { HashLink } from "react-router-hash-link";

const noImageURL = "https://static.thenounproject.com/png/1554489-200.png";

function AllPlaylists() {
  const [playlists, setPlaylists] = useState<Playlists | null>(null);
  const [spotify] = useAtom(spotifyAtom);
  const [spotifyID] = useAtom(spotifyIDAtom);

  useEffect(() => {
    (async () => {
      const resp = await fetch("/api/spotify/playlists");
      if (resp.status === 500) {
        console.error("Internal Server Error Fetching Playlists");
        return;
      }

      const playlists = await resp.json() as Playlists;
      setPlaylists(playlists);
    })();
  }, [])

  return (
    <>
      {spotify &&
        <div className="w-full flex justify-center content-center">
          <div className="m-5 sm:p-8 p-4 rounded-lg Shadow BOBorder border-2 w-[95%] min-h-[440px] h-min">
            <div className="flex flex-col"> {/* All Playlists */}
              <h1 className="font-bold lg:text-4xl md:text-2xl text-xl mb-2 text-center">Playlists</h1>
              <hr className="border border-black mb-4" />

              {playlists?.items.length === 0 &&
                <p className="text-center">
                  You have no playlists
                </p>
              }

              <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 w-[95%]">
                {playlists?.items.map((playlist, i) => (
                  playlist.owner.id === spotifyID &&
                  <div key={i} className="flex gap-2 ">
                    <img className="sm:w-20 w-16 sm:h-20 h-16 rounded-md hover:cursor-pointer border-black border" src={playlist.images[0]?.url || noImageURL} alt="playlist Image" />
                    <div className="flex flex-col">
                      <HashLink
                        to={`/dashboard/pages/playlist?id=${playlist.id}#top`}
                        className="font-bold text-xl hover:cursor-pointer hover:underline"
                      >
                        {playlist.name}
                      </HashLink>
                      <p className="sm:text-sm text-xs font-semibold italic text-gray-700">
                        {playlist.tracks.total} Tracks
                      </p>
                    </div>
                  </div>
                ))}
              </div>

            </div>
          </div>
        </div>
      }
      {!spotify && // ask user to login into spotify before viewing playlists
        <div className="w-full flex justify-center content-center">
          <div className="m-5 sm:p-8 p-4 rounded-lg Shadow BOBorder border-2 w-[95%] h-[500px]">
            <div className="flex flex-col"> {/* All Playlists */}
              <h1 className="font-bold lg:text-4xl md:text-2xl text-xl mb-2 text-center">Playlists</h1>
              <hr className="border border-black mb-4" />
              <p className="text-center">
                Please {" "}
                <HashLink
                  to="/dashboard/profile#top"
                  className="hover:cursor-pointer text-BrandBlue hover:text-blue-700"
                >
                  login to Spotify
                </HashLink>
                {" "} to view your playlists
              </p>
            </div>
          </div>
        </div>
      }
    </>
  );
}

export default AllPlaylists;