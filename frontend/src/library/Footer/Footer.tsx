import text from "@/assets/text.png";

function Footer() {
  return (
    <footer>
      <div className="border-y-2 BOBorder flex justify-between content-center sm:mx-[7%] mt-[60px] h-[500px]">

      </div>
      <div className="flex justify-center my-5">
        <img
          src={text} alt="logo-text.png"
          className="w-[200px]"
        />
      </div>
    </footer>
  );
}

export default Footer;