import opener from '@/assets/opener.png';
import s from './Landing.module.css';
import React from "react";
import { HashLink } from 'react-router-hash-link';

function Landing() {
  return (
    <div className="mx-[12%] md:pt-[60px] pt-[30px]">
      <div className="flex justify-between content-center gap-[5%]">
        <div className="lg:w-[63%] w-full">
          <h2 className={`Shadow w-fit p-5 rounded-2xl lg:text-left text-center border-[3px] BOBorder font-bold ${s.LandingText}`}>
            Empowering Your
            <span className="text-BrandOrange italic"> Musical</span>
            <span className="text-BrandBlue italic"> Journey</span>:
            Where Artistry Meets Technology
          </h2>
          <div className="flex justify-center content-center gap-[2%]">
            <p className={`border-y-[2px] BOBorder w-full mt-[3%]
               p-5 text-gray-800 lg:text-left text-center font-bold ${s.LandingParagraph}`}
            >
              <span className="text-black font-extrabold italic">Welcome to your music discovery oasis!</span> Dive into the world of melodies
              and rhythms with our intuitive and feature-packed music exploration app.
              Seamlessly search for your favorite artists, playlists, or albums to uncover
              detailed biographies, top tracks, and captivating albums.
            </p>
          </div>
        </div>
        <div className="lg:w-[40%] lg:flex hidden justify-center items-center">
          <img
            src={opener} alt="opener.png"
            className="h-[30vw]"
          />
        </div>
      </div>
      <div className="lg:hidden">
        <img
          src={opener} alt="opener.png"
          className="w-full mt-[8%]"
        />
        {
          <div className="mt-[5%] mb-[50px]">
            <h1 className="text-center text-4xl font-bold BOBorder border-y-2 m-5 p-3">Interested?</h1>
            <p></p>
            <div className="flex justify-center content-center gap-[2%]">
              <HashLink
                smooth
                to="/register#top"
                className="bg-white text-black px-4 py-2 rounded-xl BOBorder font-bold border-2
                hover:border-brandTeal focus:outline-none hover:ring-2 focus:border-black hover:ring-opacity-50"
              >
                Click here to get started!
              </HashLink>
            </div>
          </div>
        }
      </div>
    </div>
  );
}

export default Landing;