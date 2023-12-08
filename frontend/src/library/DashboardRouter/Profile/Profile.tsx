import { useAtom } from "jotai";
import { csrfTokenAtom, emailAtom, spotifyAtom, usernameAtom } from "Atoms";

function Profile() {
  const [username] = useAtom(usernameAtom);
  const [email] = useAtom(emailAtom);
  const [spotify] = useAtom(spotifyAtom);
  const [csrfToken] = useAtom(csrfTokenAtom);

  const handleSpotifyLink = async () => {
    if (!spotify) {
      // Link Spotify
      let resp = await fetch("/api/spotify/link", {
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
          window.location.href = await resp.text();
          break;
        case 401: // unauthorized (logged out)
          window.location.href = "/";
          break;
        case 403: // forbidden (invalid already linked or invalid csrf token)
          window.location.reload();
          break;
        case 500: // internal server error
          console.log("Internal Server Error Linking Spotify");
          break;
      }
    } else {
      // Unlink Spotify
      let resp = await fetch("/api/spotify/unlink", {
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
          window.location.reload();
          break;
        case 401: // unauthorized (logged out)
          window.location.href = "/";
          break;
        case 403: // forbidden (invalid already linked or invalid csrf token)
          window.location.reload();
          break;
        case 500: // internal server error
          console.log("Internal Server Error Unlinking Spotify");
          break;
      }
    }
  }

  return (
    <div className="w-full flex justify-center content-center">
      <div className="m-5 p-8 rounded-lg Shadow BOBorder border-2 lg:w-[70%] w-[95%] h-min">
        <h1 className="text-4xl font-bold mb-8 text-black italic">Your Profile</h1>
        <div className="bg-white bg-opacity-70 px-2 py-6 BOBorder border-y-2">
          <p className="md:text-lg text-base text-gray-800 mb-2">
            <span className="font-bold italic">Username: </span>{username}
          </p>
          <p className="md:text-lg text-base text-gray-800 mb-2">
            <span className="font-bold italic">Email: </span>{email}
          </p>
          <p className="md:text-lg text-base text-gray-800">
            <span className="font-bold italic">Spotify Status: </span>{spotify ? "Linked" : "Unlinked"}
          </p>
          {/* Spotify Link or Unlink button */}
          {/* In Spotify style with a spotify Icon*/}
          <button
            onClick={handleSpotifyLink}
            className="bg-green-500 px-2 py-1 rounded-xl border-black font-bold border-2 mt-4
            hover:border-brandTeal focus:outline-none hover:ring-2 focus:border-black hover:ring-opacity-50"
          >
            {spotify ? "Unlink" : "Link"} Spotify
            <i className="text-black fa-brands fa-spotify ml-2"></i>
          </button>
        </div>
      </div>
    </div>
  )
}


export default Profile;