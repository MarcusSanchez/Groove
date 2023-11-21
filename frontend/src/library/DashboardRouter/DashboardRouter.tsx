import { Link, Route, Routes, useNavigate } from "react-router-dom";
import { useAtom } from "jotai";
import { loadedAtom, loggedInAtom, sidebarAtom } from "Atoms";
import Profile from "./Profile/Profile";
import Home from "./Home/Home.tsx";
import Search from "./Search/Search.tsx";
import React from "react";

function DashboardRouter() {
  const [isLoaded] = useAtom(loadedAtom);
  const [isLoggedIn] = useAtom(loggedInAtom);
  const [sidebar, setSidebar] = useAtom(sidebarAtom);
  const navigator = useNavigate();

  if (isLoaded && !isLoggedIn) navigator("/");

  return (
    <div className="flex gap-3 sm:mx-[7%] mt-[15px] max-w-full">
      <Sidebar />
      <main className="md:mx-0 mx-[3%] w-full ">
        <div className="md:hidden w-full flex justify-end">
          <button
            onClick={() => setSidebar(sb => !sb)}
            className=" bg-white text-black px-2 py-1 rounded-xl BOBorder font-bold border-2
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
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/profile" element={<Profile />} />
          <Route path="/search" element={<Search />} />
          <Route path="*" element={<h1>404 Not Found</h1>} />
        </Routes>
      </main>
    </div>
  );
}

function Sidebar() {
  const [sidebar] = useAtom(sidebarAtom);
  const isHidden = sidebar ? "" : "md:flex hidden";

  return (
    <>
      <div className={`${isHidden} flex flex-col lg:w-64 md:w-56 w-40 bg-white max-h-fit border-gray-400 border-r`}>
        <div className="flex items-center justify-center h-14 border-gray-400 border-b">
          <div className="text-xl font-bold text-BrandBlue">Dashboard</div>
        </div>
        <div className="flex-grow p-4">
          <ul className="space-y-2">
            <li>
              <Link to="/dashboard" className="flex items-center space-x-2 text-BrandBlue font-bold hover:text-gray-800">
                <i className="text-BrandOrange fa-solid fa-house"></i>
                <span>Home</span>
              </Link>
            </li>
            <li>
              <Link to="/dashboard/search" className="flex items-center space-x-2 text-BrandBlue font-bold hover:text-gray-800">
                <i className="text-BrandOrange fa-solid fa-magnifying-glass"></i>
                <span>Search</span>
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