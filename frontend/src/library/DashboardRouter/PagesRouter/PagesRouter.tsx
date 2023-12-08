import { Route, Routes } from "react-router-dom";
import ArtistPage from "./ArtistPage/ArtistPage";
import TrackPage from "./TrackPage/TrackPage.tsx";
import AlbumPage from "./AlbumPage/AlbumPage";
import PlaylistPage from "./PlaylistPage/PlaylistPage";

function PagesRouter() {
  return (
    <Routes>
      <Route path="/artist" element={<ArtistPage />} />
      <Route path="/album" element={<AlbumPage />} />
      <Route path="/track" element={<TrackPage />} />
      <Route path="playlist" element={<PlaylistPage />} />
      <Route path="*" element={<h1>404 Not Found</h1>} />
    </Routes>
  );
}

export default PagesRouter;