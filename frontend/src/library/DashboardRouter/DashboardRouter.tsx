import { Link, Route, Routes, useNavigate } from "react-router-dom";
import { useAtom } from "jotai";
import { loadedAtom, loggedInAtom, sidebarAtom } from "Atoms";
import Profile from "./Profile/Profile";
import Search from "./Search/Search.tsx";
import PagesRouter from "./PagesRouter/PagesRouter.tsx";
import AllPlaylists from "./AllPlaylists/AllPlaylists.tsx";
import NotFound from "@/library/NotFound/NotFound.tsx";

function DashboardRouter() {
  const [isLoaded] = useAtom(loadedAtom);
  const [isLoggedIn] = useAtom(loggedInAtom);
  const [sidebar, setSidebar] = useAtom(sidebarAtom);
  const navigate = useNavigate();

  if (isLoaded && !isLoggedIn) navigate("/");

  return (
    <div className="flex gap-3 sm:mx-[7%] mt-[15px] max-w-full">
      <Sidebar />
      <main className="md:mx-0 mx-[3%] w-full ">
        <div className="md:hidden w-full flex justify-end">
          <button
            onClick={() => setSidebar(sb => !sb)}
            className="bg-white text-black px-2 py-1 rounded-xl BOBorder font-bold border-2
            hover:border-brandTeal focus:outline-none hover:ring-2 focus:border-black hover:ring-opacity-50"
          >
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              {sidebar
                ? /* X Icon */
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12"></path>
                : /* Hamburger Icon */
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4 6h16M4 12h16m-7 6h7"></path>
              }
            </svg>
          </button>
        </div>
        <div className={`flex justify-end md:hidden ${sidebar ? "" : "hidden"}`}>
          <div className="w-min flex justify-center items-centerbg-white text-black px-2 py-1 rounded-xl BOBorder font-bold border-2">
            <ul className="flex gap-5">
              <li>
                <Link to="/dashboard" className="flex items-center space-x-2 text-BrandBlue font-bold hover:text-gray-800">
                  <i className="text-BrandOrange fa-solid fa-magnifying-glass"></i>
                  <span>Home</span>
                </Link>
              </li>
              <li>
                <Link to="/dashboard/playlists" className="flex items-center space-x-2 text-BrandBlue font-bold hover:text-gray-800">
                  <i className="text-BrandOrange fa-solid fa-play"></i>
                  <span>Playlists</span>
                </Link>
              </li>
              <li>
                <Link to="/dashboard/profile" className="flex items-center space-x-2 text-BrandBlue font-bold hover:text-gray-800">
                  <i className="text-BrandOrange fa-solid fa-user"></i>
                  <span>Profile</span>
                </Link>
              </li>
            </ul>
          </div>
        </div>
        <Routes>
          <Route path="/" element={<Search />} />
          <Route path="/playlists" element={<AllPlaylists />} />
          <Route path="/profile" element={<Profile />} />
          <Route path="/pages/*" element={<PagesRouter />} />
          <Route path="*" element={<NotFound />} />
        </Routes>
      </main>
    </div>
  );
}

function Sidebar() {
  return (
    <>
      <div className={`md:flex hidden flex-col lg:w-64 md:w-56 w-40 bg-white max-h-fit border-gray-400 border-r`}>
        <div className="flex items-center justify-center h-14 border-gray-400 border-b">
          <div className="text-xl font-bold text-BrandBlue">Dashboard</div>
        </div>
        <div className="flex-grow p-4">
          <ul className="space-y-2">
            <li>
              <Link to="/dashboard" className="flex items-center space-x-2 text-BrandBlue font-bold hover:text-gray-800">
                <i className="text-BrandOrange fa-solid fa-magnifying-glass"></i>
                <span>Home</span>
              </Link>
            </li>
            <li>
              <Link to="/dashboard/playlists" className="flex items-center space-x-2 text-BrandBlue font-bold hover:text-gray-800">
                <i className="text-BrandOrange fa-solid fa-play"></i>
                <span>Playlists</span>
              </Link>
            </li>
            <li>
              <Link to="/dashboard/profile" className="flex items-center space-x-2 text-BrandBlue font-bold hover:text-gray-800">
                <i className="text-BrandOrange fa-solid fa-user"></i>
                <span>Profile</span>
              </Link>
            </li>
          </ul>
        </div>
      </div>
    </>
  );
}

export default DashboardRouter;
