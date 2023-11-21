import { useEffect, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";

type Image = {
  height: number;
  url: string;
  width: number;
};

type ListedArtists = {
  href: string;
  id: string;
  name: string;
  type: string;
};

type Album = {
  album_type: string;
  artists: ListedArtists[];
  id: string;
  images: Image[];
  name: string;
  release_date: string;
  type: string;
};

type Track = {
  album: Album;
  artists: ListedArtists[];
  duration_ms: number;
  id: string;
  name: string;
  // uri: string;
};

type Artist = {
  genres: string[];
  id: string;
  images: Image[];
  name: string;
  popularity: number;
  type: string;
}

type SearchResult = {
  albums?: {
    items: Album[];
  };
  artists?: {
    items: Artist[];
  };
  tracks?: {
    items: Track[];
  };
};

enum SearchType {
  Album = "album",
  Artist = "artist",
  Track = "track",
}

function Search() {
  const albums = useRef<HTMLInputElement>(null);
  const artists = useRef<HTMLInputElement>(null);
  const tracks = useRef<HTMLInputElement>(null);
  const searchbar = useRef<HTMLInputElement>(null);
  const searchButton = useRef<HTMLButtonElement>(null);

  const navigate = useNavigate();
  const redirect = (id: string, type: SearchType) => {
    switch (type) {
      case SearchType.Album:
        navigate(`/dashboard/album?id=${id}`);
        break;
      case SearchType.Artist:
        navigate(`/dashboard/artist?id=${id}`);
        break;
      case SearchType.Track:
        navigate(`/dashboard/track?id=${id}`);
        break;
    }

  }

  const [results, setResults] = useState<SearchResult | null>(null);

  const handleSearch = async () => {
    let q = searchbar.current!.value;
    let typeArray = [] as string[];

    if (albums.current!.checked) typeArray.push("album");
    if (artists.current!.checked) typeArray.push("artist");
    if (tracks.current!.checked) typeArray.push("track");

    let type = typeArray.join(",");
    let url = new URL(window.location.href);
    url.searchParams.set("q", q);
    url.searchParams.set("type", type);

    let endpoint = `/api/spotify/search/${q}?type=` + type;
    let resp = await fetch(endpoint, { credentials: "include" });
    let data: SearchResult = await resp.json();

    // set a unique identifier to cache the previous search
    let randomID = Math.random().toString(36).substring(7);
    url.searchParams.set("cache", randomID);

    localStorage.setItem("search", JSON.stringify({
      data: data,
      id: randomID
    }))

    window.history.pushState({}, "", url.toString());
    setResults(data);
  }

  useEffect(() => {
    // grab the q param and the type param from the url if they exist.
    // if they do, set the checkboxes and search bar to the values.
    // then execute the search.

    let url = new URL(window.location.href);
    let q = url.searchParams.get("q");
    let type = url.searchParams.get("type");
    let cacheID = url.searchParams.get("cache");

    if (q) searchbar.current!.value = q;
    else return;

    if (type) {
      albums.current!.checked = type.includes("album");
      artists.current!.checked = type.includes("artist");
      tracks.current!.checked = type.includes("track");
    }

    if (cacheID) {
      let cacheData = localStorage.getItem("search");
      if (cacheData) {
        let cacheJSON = JSON.parse(cacheData);
        if (cacheJSON.id === cacheID) {
          setResults(cacheJSON.data);
          return;
        }
      }
    }

    handleSearch().then();
  }, []);


  return (
    <div className="w-full flex justify-center content-center">
      <div className="m-5 sm:p-8 p-4 rounded-lg Shadow BOBorder border-2 w-[95%] h-min">
        <h1 className="text-4xl font-bold mb-4 text-black text-center italic">Search</h1>

        <div className="bg-white bg-opacity-70 px-2 py-6 BOBorder border-y-2">
          {/* Checkboxes for including albums, artists, and tracks */}
          <div className="flex justify-center mb-4">
            <label className="inline-flex items-center mx-2">
              <input
                ref={albums}
                type="checkbox"
                name="albums"
                defaultChecked
                className="form-checkbox h-[17px] w-[17px] text-black font-bold"
              />
              <span className="ml-2 text-black">Albums</span>
            </label>

            <label className="inline-flex items-center mx-2">
              <input
                ref={artists}
                type="checkbox"
                name="artists"
                defaultChecked
                className="form-checkbox h-[17px] w-[17px] text-black font-bold"
              />
              <span className="ml-2 text-black">Artists</span>
            </label>

            <label className="inline-flex items-center mx-2">
              <input
                ref={tracks}
                type="checkbox"
                name="tracks"
                defaultChecked
                className="form-checkbox h-[17px] w-[17px] text-black font-bold"
              />
              <span className="ml-2 text-black">Tracks</span>
            </label>
          </div>

          {/* Search bar and button */}
          <div className=" flex ">
            <input
              ref={searchbar}
              type="text"
              onKeyDown={(e) => (e.key === "Enter") ? handleSearch() : null}
              placeholder="Search"
              className="w-full h-10 px-3 text-base placeholder-black border-2 border-gray-600 rounded-l-xl drop-shadow-md
                focus:shadow-outline focus:outline-none focus:border-"
            />
            <button
              ref={searchButton}
              onClick={handleSearch}
              className="bg-white text-black px-2 py-1 rounded-r-xl BOBorder font-bold border-2 drop-shadow-md
                hover:border- focus:outline-none hover:ring-2 focus:border-black hover:ring-opacity-50"
            >
              Search
            </button>
          </div>
        </div>

        {/*  Results section  */}
        <div className="mt-4">
          <h1 className="text-2xl font-bold mb-4 text-black text-center italic">Results</h1>
          <div className="flex justify-center">
            <div className="w-[95%] h-[500px] overflow-scroll bg-white bg-opacity-70 BOBorder border-2 rounded-lg">
              <div className="m-5">
                {!results &&
                  <h1 className="text-2xl font-bold text-black text-center italic">
                    Just search for something!
                  </h1>
                }
                {results && (
                  <div className="grid grid-cols-1 lg:grid-cols-2 w-full gap-4">
                    {results.albums?.items.map((album) => (
                      <div key={album.id} className="bg-white bg-opacity-70 border-black border-2 rounded-lg p-4 flex justify-between lg:gap-1 gap-4">
                        <div className="flex flex-col justify-between">
                          <div>
                            <h2
                              onClick={() => redirect(album.id, SearchType.Album)}
                              className="lg:text-lg text-base font-bold hover:cursor-pointer"
                            >
                              {album.name}
                            </h2>
                            <p className="text-sm font-bold text-gray-700">{album.artists.map((artist) => artist.name).join(', ')}</p>
                            <p className="text-sm text-gray-500">{album.release_date.slice(0, 4)}</p>
                          </div>
                          <p className="text-sm text-gray-700 font-bold">album</p>
                        </div>
                        <img
                          onClick={() => redirect(album.id, SearchType.Album)}
                          src={album.images[0]?.url} alt={album.name} className="lg:h-32 h-24 lg:w-32 w-24 rounded-xl hover:cursor-pointer border-black border"
                        />
                      </div>
                    ))}
                    {results.artists?.items.map((artist) => (
                      <div key={artist.id} className="bg-white bg-opacity-70 border-black border-2 rounded-lg p-4 flex justify-between lg:gap-1 gap-4">
                        <div className="flex flex-col justify-between">
                          <div>
                            <h2
                              onClick={() => redirect(artist.id, SearchType.Artist)}
                              className="lg:text-lg text-base font-bold hover:cursor-pointer"
                            >
                              {artist.name}
                            </h2>
                            <p className="text-sm text-gray-500">{artist.genres.join(', ')}</p>
                          </div>
                          <p className="text-sm text-gray-700 font-bold">artist</p>
                        </div>
                        <img
                          onClick={() => redirect(artist.id, SearchType.Artist)}
                          src={artist.images[0]?.url} alt={artist.name} className="lg:h-32 h-24 lg:w-32 w-24 rounded-xl hover:cursor-pointer border-black border"
                        />
                      </div>
                    ))}
                    {results.tracks?.items.map((track) => (
                      <div key={track.id} className="bg-white bg-opacity-70 border-black border-2 rounded-lg p-4 flex justify-between lg:gap-1 gap-4">
                        <div className="flex flex-col justify-between">
                          <div>
                            <h2
                              onClick={() => redirect(track.id, SearchType.Artist)}
                              className="lg:text-lg text-base font-bold hover:cursor-pointer"
                            >
                              {track.name}
                            </h2>
                            <p className="text-sm font-bold text-gray-700">{track.artists.map((artist) => artist.name).join(', ')}</p>
                            <p className="text-sm font-bold text-gray-500">{track.album.name}</p>
                            <p className="text-sm text-gray-500">{msToMinutesSeconds(track.duration_ms)}</p>
                          </div>
                          <p className="text-sm text-gray-700 font-bold">track</p>
                        </div>
                        <img
                          onClick={() => redirect(track.id, SearchType.Artist)}
                          src={track.album.images[0]?.url} alt={track.album.name} className="lg:h-32 h-24 lg:w-32 w-24 rounded-xl hover:cursor-pointer border-black border"
                        />
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>

      </div>
    </div>
  );
}

function msToMinutesSeconds(durationMs: number): string {
  // Convert milliseconds to seconds
  const durationSec: number = Math.floor(durationMs / 1000);

  // Calculate minutes and seconds
  const minutes: number = Math.floor(durationSec / 60);
  const seconds: number = durationSec % 60;

  // Format the result
  const result: string = `${minutes}:${seconds.toString().padStart(2, '0')}`; // Ensure seconds are displayed with leading zero if needed

  return result;
}


export default Search;
