import { FormEvent, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";
import { getCookie } from "@/util.ts";
import { useAtom } from "jotai";
import { csrfTokenAtom, emailAtom, loggedInAtom, usernameAtom } from "Atoms";
import { HashLink } from "react-router-hash-link";

function Login() {
  const navigate = useNavigate();
  const form = useRef<HTMLFormElement>(null)
  const button = useRef<HTMLButtonElement>(null)
  const [error, setError] = useState("")

  const [loggedIn, setLoggedIn] = useAtom(loggedInAtom);
  const [, setUsername] = useAtom(usernameAtom);
  const [, setEmail] = useAtom(emailAtom);
  const [, setCsrfToken] = useAtom(csrfTokenAtom);

  if (loggedIn) navigate("/");

  const login = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    button.current!.disabled = true;
    setError("");

    let username = form.current?.username.value;
    let password = form.current?.password.value;

    if (!username || !password) {
      setError("Please fill out all fields");
      button.current!.disabled = false;
      return;
    }

    const data = {
      username: username,
      password: password,
    };

    let resp = await fetch("/api/login", {
      method: "POST",
      credentials: "include",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data)
    })
    let json;

    switch (resp.status) {
      case 400: // failed validation
      case 500: // internal server error
        json = await resp.json();
        setError(json.message);
        button.current!.disabled = false;
        break;
      case 201: // success
        json = await resp.json();
        setLoggedIn(true);
        setUsername(json.username);
        setEmail(json.email);
        setCsrfToken(getCookie("Csrf")); // the cookie will be there
        navigate("/");
        break;
      case 308: // already authenticated
        window.location.href = "/";
        break;
    }
  }

  return (
    <>
      <div className="mt-10 mb-5 sm:mx-[20%] mx-[10%] flex justify-center">
        <form
          id="login-form"
          ref={form}
          onSubmit={login}
          className="Shadow w-[400px] min-h-[440px] rounded-2xl border-[1px] BOBorder"
        >
          <div className="mx-10 my-5">

            <h1 className="md:text-4xl text-3xl font-bold pt-5 text-black">
              Welcome Back
            </h1>
            <p className="font-bold mt-1 mb-5 text-brandOrange">Did you miss us?</p>

            <p className="text-lg text-black font-bold">
              <i className="fa-regular fa-user"></i> {" "}
              Username:
            </p>
            <Input name="username" type="text" />

            <p className="text-lg text-black font-bold">
              <i className="fa-solid fa-key"></i> {" "}
              Password:
            </p>
            <Input name="password" />

            {error && <p className="text-red-500 min-w-fit font-bold text-center">{error}</p>}
            <button
              ref={button}
              type="submit"
              className="BOBorder text-black bg-white border-2 w-full md:text-2xl text-xl
              font-bold rounded-3xl py-2 mt-5 mb-5 hover:ring-2
              focus:border-black hover:ring-opacity-50 disabled:opacity-50"
            >
              Login {" "}
              <i className="fa-solid fa-arrow-right-to-bracket"></i>
            </button>

          </div>
        </form>
      </div>

      <p className="text-center font-bold text-base mb-20">
        Don't have an account yet? {" "}
        <HashLink
          smooth
          to="/register#top"
          className="underline text-blue-500 hover:text-blue-700 focus:text-blue-900"
        >
          Register
        </HashLink>
      </p>
    </>
  );
}

export default Login;

function Input({ name, type }: { name: string, type?: string }) {
  return (
    <input
      name={name} type={type || name} placeholder={"your " + name + "..."} autoComplete="off"
      className="mb-5 w-full mx-auto px-3 py-2 rounded-md BOBorder md:text-xl text-lg border-2
      focus:ring-2 focus:ring-blue-500 focus:border-transparent"
    />
  );
}