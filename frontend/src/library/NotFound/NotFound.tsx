function NotFound() {
  return (
    <div className="flex gap-3 sm:mx-[7%] mt-[15px] max-w-full">
      <div className="w-full flex justify-center content-center">
        <div className="m-5 sm:p-8 p-4 rounded-lg Shadow BOBorder border-2 w-[95%] h-[500px]">
          <h1 className="text-center text-3xl font-bold">404 Not Found</h1>
          <p className="text-center text-lg font-semibold">The page you are looking for does not exist.</p>
        </div>
      </div>
    </div>
  )
}

export default NotFound;