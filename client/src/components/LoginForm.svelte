<script lang="ts">
  import { toast } from "@zerodevx/svelte-toast";
  import type { AxiosError } from "axios";
  import { navigate } from "svelte-navigator";
  import AuthAPI from "../API/Auth";
  import { API } from "../constants/API";
  import type { UserLogin } from "../models/Auth";
  import { AuthState } from "../store/Auth";

  export let toggleRegister;
  export let show = null;
  export let togglePasswordModal;

  let email = "";
  let password = "";

  function LoginUser() {
    const user: UserLogin = {
      email,
      password,
    };

    AuthAPI.post(API.Routes.Auth.Login, JSON.stringify(user))
      .then((res) => {
        AuthState.update((value) => {
          value.accessToken = res.data["access_token"];
          value.idToken = res.data["id_token"];
          return value;
        });
        navigate("/");
      })
      .catch((err: AxiosError) => {
        if (err.response!.status === 401) {
          toast.push("check your credentials.");
          return;
        }
        if (err.response!.status === 400) {
          toast.push("bad request.");
          return;
        }
        if (err.response!.status === 500) {
          toast.push("something went wrong.");
          return;
        }

        toast.push("an unknown error occured.");
      });
  }
</script>

{#if show}
  <form class="auth-form" on:submit|preventDefault={LoginUser}>
    <div class="title-con">
      <h2>Login</h2>
    </div>
    <div class="inpt-con">
      <input
        type="email"
        name="email"
        id="email"
        placeholder="Email"
        bind:value={email}
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
      <button class="btn btn-auth" type="submit"> Login </button>
      <button
        type="button"
        class="btn need_help"
        on:click={togglePasswordModal}
      >
        forgot password?
      </button>
      <div class="auth_toggle">
        <button type="button" on:click={toggleRegister}>
          Don't have an account yet? Register
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

    .need_help {
      color: #800080;
      padding: 1rem;
      font-size: 1rem;
    }
  }
</style>
