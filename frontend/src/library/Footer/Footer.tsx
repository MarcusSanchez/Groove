import text from "@/assets/text.png";

function Footer() {
  return (
    <footer>
      <div className="border-y-2 BOBorder flex flex-wrap gap-2 justify-center sm:mx-[7%] md:mt-[150px] mt-[60px] p-5">
        <p className="font-semibold italic text-gray-700">
          Made by Marcus Sanchez
        </p>
        <p className="font-semibold italic text-gray-700">
          |
        </p>
        <p className="font-semibold italic text-gray-700">
          Repository:
          <a
            href="https://github.com/MarcusSanchez/Groove"
            className="text-blue-500 hover:text-blue-700"
          > Here</a>
        </p>
      </div>
      <div className="flex justify-center my-5">
        <img
          src={text} alt="logo-text.png"
          className="sm:w-[200px] w-[150px]"
        />
      </div>
    </footer>
  );
}

export default Footer;