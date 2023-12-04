import { useEffect, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";
import { Album, Artist, DataType, SearchResult, SearchType, Track } from "./types.ts";
import { msToMinutesSeconds, sortResults } from "./util.ts";
import { HashLink } from "react-router-hash-link";

const noImageURL = "https://static.thenounproject.com/png/1554489-200.png";

function Search() {
  const albums = useRef<HTMLInputElement>(null);
  const artists = useRef<HTMLInputElement>(null);
  const tracks = useRef<HTMLInputElement>(null);
  const searchbar = useRef<HTMLInputElement>(null);
  const searchButton = useRef<HTMLButtonElement>(null);
  const resultsDiv = useRef<HTMLDivElement>(null);

  const [results, setResults] = useState<DataType[] | null>(null);

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

    let endpoint = `/api/spotify/search/${q.replaceAll(" ", "+")}?type=` + type;
    let resp = await fetch(endpoint, { credentials: "include" });
    let data: SearchResult = await resp.json();

    // set a unique identifier to cache the previous search
    let randomID = Math.random().toString(36).substring(7);
    url.searchParams.set("cache", randomID);

    let sortedResults = sortResults(data.artists?.items, data.albums?.items, data.tracks?.items);
    localStorage.setItem("search", JSON.stringify({
      data: sortedResults,
      id: randomID
    }))

    window.history.pushState({}, "", url.toString());
    setResults(sortedResults);
    resultsDiv.current!.scrollTo({ top: 0, behavior: "smooth" });
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
            <div ref={resultsDiv} className="w-[95%] h-[500px] overflow-scroll bg-white bg-opacity-70 BOBorder border-2 rounded-lg">
              <div className="m-5">
                {!results &&
                  <h1 className="text-2xl font-bold text-black text-center italic">
                    Just search for something!
                  </h1>
                }
                {results &&
                  <div className="grid grid-cols-1 lg:grid-cols-2 w-full gap-4">
                    <DisplayResults results={results} />
                  </div>
                }
              </div>
            </div>
          </div>
        </div>

      </div>
    </div>
  );
}

function DisplayResults({ results }: { results: (DataType[] | null) }) {
  const navigate = useNavigate();
  return (
    <>
      {results?.map((result) => {
        switch (result.type) {
          case SearchType.Album:
            const album = result.data as Album;
            return (
              <div key={album.id} className="bg-white bg-opacity-70 border-black border-2 rounded-lg p-4 flex justify-between lg:gap-1 gap-4">
                <div className="flex flex-col justify-between">
                  <div>
                    <HashLink
                      to={`/dashboard/pages/album?id=${album.id}#top`}
                      className="lg:text-lg text-base font-bold hover:cursor-pointer"
                    >
                      {album.name}
                    </HashLink>
                    <p className="text-sm font-bold text-gray-700">{album.artists.map((artist) => artist.name).join(', ')}</p>
                    <p className="text-sm text-gray-500">{album.release_date.slice(0, 4)}</p>
                  </div>
                  <p className="text-sm text-gray-700 font-bold">album</p>
                </div>
                <HashLink to={`/dashboard/pages/album?id=${album.id}#top`}>
                  <img
                    src={album.images[0]?.url || noImageURL} alt={album.name}
                    className="lg:h-32 h-24 lg:w-32 w-24 rounded-xl hover:cursor-pointer border-black border"
                  />
                </HashLink>
              </div>
            )
          case SearchType.Artist:
            const artist = result.data as Artist;
            return (
              <div
                key={artist.id}
                className="bg-white bg-opacity-70 border-black border-2
                rounded-lg p-4 flex justify-between lg:gap-1 gap-4"
              >
                <div className="flex flex-col justify-between">
                  <div>
                    <HashLink
                      to={`/dashboard/pages/artist?id=${artist.id}#top`}
                      className="lg:text-lg text-base font-bold hover:cursor-pointer"
                    >
                      {artist.name}
                    </HashLink>
                    <p className="text-sm text-gray-500">{artist.genres.join(', ')}</p>
                  </div>
                  <p className="text-sm text-gray-700 font-bold">artist</p>
                </div>
                <HashLink to={`/dashboard/pages/artist?id=${artist.id}#top`}>
                  <img
                    src={artist.images[0]?.url || noImageURL} alt={artist.name}
                    className="lg:h-32 h-24 lg:w-32 w-24 rounded-xl hover:cursor-pointer border-black border"
                  />
                </HashLink>
              </div>
            );
          case SearchType.Track:
            const track = result.data as Track;
            return (
              <div
                key={track.id}
                className="bg-white bg-opacity-70 border-black border-2 rounded-lg p-4 flex justify-between lg:gap-1 gap-4"
              >
                <div className="flex flex-col justify-between">
                  <div>
                    <HashLink
                      to={`/dashboard/pages/track?id=${track.id}#top`}
                      className="lg:text-lg text-base font-bold hover:cursor-pointer"
                    >
                      {track.name}
                    </HashLink>
                    <p className="text-sm font-bold text-gray-700">{track.artists.map((artist) => artist.name).join(', ')}</p>
                    <p className="text-sm font-bold text-gray-500">{track.album.name}</p>
                    <p className="text-sm text-gray-500">{msToMinutesSeconds(track.duration_ms)}</p>
                  </div>
                  <p className="text-sm text-gray-700 font-bold">track</p>
                </div>
                <HashLink to={`/dashboard/pages/track?id=${track.id}#top`}>
                  <img
                    src={track.album.images[0]?.url || noImageURL} alt={track.album.name}
                    className="lg:h-32 h-24 lg:w-32 w-24 rounded-xl hover:cursor-pointer border-black border"
                  />
                </HashLink>
              </div>
            );
        }
      })}
    </>
  );
}

export default Search;
