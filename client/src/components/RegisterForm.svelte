<script lang="ts">
  import { toast } from "@zerodevx/svelte-toast";
  import type { AxiosError } from "axios";
  import AuthAPI from "../API/Auth";
  import { API } from "../constants/API";
  import type { UserLogin, UserRegister } from "../models/Auth";

  export let toggleLogin;
  export let show = null;

  let isPasswordModalOpen = false;

  let username = "";
  let email = "";
  let password = "";
  let checkPassword = "";

  function RegisterUser() {
    let isInfoValid = true;
    if (password !== checkPassword) {
      isInfoValid = false;
      toast.push("passwords don't match");
    }

    if (email === "") {
      isInfoValid = false;
      toast.push("email is empty.");
    }

    if (username === "") {
      isInfoValid = false;
      toast.push("username is empty.");
    }

    if (!isInfoValid) {
      return;
    }

    const user: UserRegister = {
      username,
      email,
      password,
    };

    AuthAPI.post(API.Routes.Auth.Register, JSON.stringify(user))
      .then(() => {
        toast.push(
          "your account was successfully registered. now you can log in to your account."
        );
        toggleLogin();
      })
      .catch((err: AxiosError) => {
        if (err.response!.status === 400) {
          toast.push("invalid information");
          return;
        }
        if (err.response!.status === 409) {
          toast.push("user with the same email or username already exists.");
          return;
        }
        if (err.response!.status === 500) {
          toast.push("something went wrong.");
          return;
        }

        toast.push("An unknown error occured");
      });
  }
</script>

{#if show}
  <form class="auth-form" on:submit|preventDefault={RegisterUser}>
    <div class="title-con">
      <h2>Register</h2>
    </div>
    <div class="inpt-con">
      <input
        type="username"
        name="username"
        id="username"
        placeholder="Username"
        bind:value={username}
      />
      <input
        type="email"
        name="email"
        id="email"
        placeholder="Email"
        bind:value={email}
      />
      <input
        type="password"
        name="checkPassword"
        id="checkPassword"
        placeholder="Re-enter your password"
        bind:value={checkPassword}
      />
      <input
        type="password"
        name="password"
        id="password"
        placeholder="Password"
        bind:value={password}
      />
    </div>
    <div class="btn-con">
      <button class="btn btn-auth" type="submit"> Register </button>
      <div class="auth_toggle">
        <button type="button" on:click={toggleLogin}>
          Already have an account? Login
        </button>
      </div>
    </div>
  </form>
{/if}

<style lang="scss">
  .auth-form {
    display: flex;
    align-items: center;
    justify-content: space-between;
    flex-direction: column;
    padding: 2rem;
    background-color: #fff;
    border-radius: 0.3rem;
    box-shadow: 0 0 50px rgba($color: #000000, $alpha: 0.3);

    & > div {
      width: 100%;
      padding: 0.5rem;
    }

    .title-con {
      height: 3rem;
      display: flex;
      align-items: center;
      justify-content: center;
    }

    .inpt-con {
      display: flex;
      align-items: center;
      justify-content: flex-start;
      flex-direction: column;

      input {
        border: none;
        outline: none;
        width: 100%;
        background-color: transparent;
        padding: 1rem;
        margin: 0.3rem 0;
        border-top-left-radius: 0.3rem;
        border-top-right-radius: 0.3rem;
        border-bottom: 2px solid rgb(180, 180, 180);
      }

      input:focus {
        background-color: rgba(132, 0, 255, 0.1);
        border-bottom: 2px solid #800080;
      }
    }

    .btn-con {
      display: flex;
      align-items: center;
      flex-direction: column;

      .auth_toggle {
        width: 100%;
        padding: 1rem;
        display: flex;
        align-items: center;
        justify-content: center;

        button {
          border: none;
          outline: none;
          background-color: transparent;
          color: #800080;
          font-size: 1rem;
        }

        button:hover {
          text-decoration: underline;
          cursor: pointer;
        }
      }
    }

    .btn-auth {
      border: none;
      outline: none;
      padding: 1rem;
      width: 100%;
      background-color: #800080;
      border-radius: 0.3rem;
      color: #fff;
      font-size: 1rem;
    }

    .btn-auth:hover {
      background-color: #500050;
    }
  }
</style>
