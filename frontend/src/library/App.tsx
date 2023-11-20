import Nav from "./Nav/Nav.tsx";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import Landing from "./Landing/Landing.tsx";
import Login from "./Login/Login.tsx";
import Register from "./Register/Register.tsx";
import { useEffect } from "react";
import { getCookie } from "@/util.ts";
import { useAtom } from "jotai";
import Footer from "./Footer/Footer.tsx";
import { csrfTokenAtom, emailAtom, loadedAtom, loggedInAtom, spotifyAtom, usernameAtom } from "Atoms";
import DashboardRouter from "./DashboardRouter/DashboardRouter.tsx";

function App() {
  const [, setLoggedIn] = useAtom(loggedInAtom);
  const [, setUsername] = useAtom(usernameAtom);
  const [, setEmail] = useAtom(emailAtom);
  const [, setCsrfToken] = useAtom(csrfTokenAtom);
  const [, setIsLoaded] = useAtom(loadedAtom);
  const [, setSpotify] = useAtom(spotifyAtom);

  useEffect(() => {
    (async () => {
      let csrfToken = getCookie("Csrf");
      if (csrfToken === "") {
        setIsLoaded(true);
        return;
      }

      let resp = await fetch("/api/authenticate", {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          csrf_: csrfToken,
        }),
      })

      interface Payload {
        user: {
          username: string,
          email: string,
          spotify: boolean
        }
      }

      switch (resp.status) {
        case 200: // success
          let data = await resp.json() as Payload;
          setLoggedIn(true);
          setUsername(data.user.username)
          setSpotify(data.user.spotify)
          setEmail(data.user.email)
          setCsrfToken(csrfToken)
          break;
        case 401: // unauthorized (already logged out)
          break;
        case 403: // forbidden (invalid csrf token)
          window.location.href = "/";
          break;
        case 500: // internal server error
          console.error("Internal Server Error Authenticating User");
          break;
      }
      setIsLoaded(true);
    })();
  }, [])

  return (
    <BrowserRouter>
      <Nav />
      <Routes>
        <Route path="/" element={<Landing />} />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/dashboard/*" element={<DashboardRouter />} />
        <Route path="*" element={<h1>404 Not Found</h1>} />
      </Routes>
      <Footer />
    </BrowserRouter>
  );
}

export default App
