import logo from "@/assets/logo.png";
import { Link } from "react-router-dom";
import { useAtom } from "jotai";
import { csrfTokenAtom, emailAtom, loggedInAtom, usernameAtom } from "Atoms";

function Nav() {
  const [loggedIn, setLoggedIn] = useAtom(loggedInAtom);
  const [csrfToken, setCsrfToken] = useAtom(csrfTokenAtom);
  const [, setUsername] = useAtom(usernameAtom);
  const [, setEmail] = useAtom(emailAtom);

  const handleLogout = async () => {
    let resp = await fetch("/api/logout", {
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
      case 204: // success
        setLoggedIn(false)
        setUsername('')
        setEmail('')
        setCsrfToken('')
        break;
      case 401: // unauthorized (already logged out)
      case 403: // forbidden (invalid csrf token)
        window.location.href = "/";
        return;
      case 500: // internal server error
        console.error("Internal Server Error Logging Out User");
        return;
    }
  }

  return (
    <nav className="border-y-2 border-b-BrandOrange border-t-BrandBlue
      flex justify-between content-center sm:mx-[50px] mt-[15px]"
    >
      <Link to="/">
        <img
          src={logo} alt="groove-guru-logo.png"
          className="sm:h-[5rem] h-[4rem] mx-[20px] my-[15px]"
        />
      </Link>
      <div className="flex items-center space-x-5 mx-4">
        {!loggedIn &&
          <>
            <Link
              to="/login"
              className="bg-white text-black px-4 py-2 rounded-xl BOBorder font-bold border-2
              hover:border-brandTeal focus:outline-none hover:ring-2 focus:border-black hover:ring-opacity-50"
            >
              Login
            </Link>
            <Link
              to="/register"
              className="bg-white text-black px-4 py-2 rounded-xl BOBorder font-bold border-2
              hover:border-brandTeal focus:outline-none hover:ring-2 focus:border-black hover:ring-opacity-50"
            >
              Register
            </Link>
          </>
        }
        {loggedIn &&
          <button
            onClick={handleLogout}
            className="bg-white text-black px-4 py-2 rounded-xl BOBorder font-bold border-2
            hover:border-brandTeal focus:outline-none hover:ring-2 focus:border-black hover:ring-opacity-50"
          >
            Logout
          </button>
        }
      </div>
    </nav>
  )
}

export default Nav;