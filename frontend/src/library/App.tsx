import Nav from "./Nav/Nav.tsx";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import Landing from "./Landing/Landing.tsx";
import Login from "./Login/Login.tsx";
import Register from "./Register/Register.tsx";
import { useEffect } from "react";
import { getCookie } from "@/util.ts";
import { useAtom } from "jotai";
import { csrfTokenAtom, emailAtom, loggedInAtom, usernameAtom } from "Atoms";

function App() {
  const [, setLoggedIn] = useAtom(loggedInAtom);
  const [, setUsername] = useAtom(usernameAtom);
  const [, setEmail] = useAtom(emailAtom);
  const [, setCsrfToken] = useAtom(csrfTokenAtom);

  useEffect(() => {
    (async () => {
      let csrfToken = getCookie("Csrf")
      if (csrfToken === undefined) return;

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

      switch (resp.status) {
        case 200: // success
          let data = await resp.json();
          setLoggedIn(true);
          setUsername(data.username)
          setEmail(data.email)
          setCsrfToken(csrfToken)
          break;
        case 401: // unauthorized (already logged out)
          return;
        case 403: // forbidden (invalid csrf token)
          window.location.href = "/";
          return;
        case 500: // internal server error
          console.error("Internal Server Error Authenticating User");
          return;
      }
    })();
  }, [])

  return (
    <BrowserRouter>
      <Nav />
      <Routes>
        <Route path="/" element={<Landing />} />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App
