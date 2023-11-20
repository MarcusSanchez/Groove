import { useAtom } from "jotai";
import { emailAtom, spotifyAtom, usernameAtom } from "Atoms";

function Profile() {
  const [username] = useAtom(usernameAtom);
  const [email] = useAtom(emailAtom);
  const [spotify] = useAtom(spotifyAtom);

  return (
    <div className="w-full flex justify-center content-center">
      <div className="m-5 p-8 rounded-lg Shadow BOBorder border-2 lg:w-[50%] min-w-min h-min">
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
        </div>
      </div>
    </div>
  )
}


export default Profile;