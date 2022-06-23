<script lang="ts">
  import { toast } from "@zerodevx/svelte-toast";
  import type { AxiosError } from "axios";

  import AuthAPI from "../API/Auth";
  import { API } from "../constants/API";

  export let show;
  export let toggleLogin;
  export let togglePasswordModal;

  let stage: 1 | 2 = 1;
  let email;
  let token;
  let new_password;
  let checkNewPassword;
  let title = "Reset Password";
  let disableSendTokenBtn = false;

  function SendResetToken() {
    disableSendTokenBtn = true;
    title = "Please wait...";
    AuthAPI.post(
      API.Routes.Account.ResetPasswordToken,
      JSON.stringify({
        email,
      })
    )
      .then(() => {
        disableSendTokenBtn = false;
        toast.push(
          "if your email is registered, and email containing a verification code is underway."
        );
        stage = 2;
      })
      .catch((err: AxiosError) => {
        toast.push("An unknown error occured.");
      });
  }

  function ResetPassword() {
    let isInfoValid = true;
    if (new_password !== checkNewPassword) {
      isInfoValid = false;
      toast.push("passwords don't match");
    }

    if (email === "") {
      isInfoValid = false;
      toast.push("email is empty.");
    }

    if (token === "") {
      isInfoValid = false;
      toast.push("token is empty.");
    }

    if (!isInfoValid) {
      return;
    }

    AuthAPI.post(
      API.Routes.Account.ResetPassword,
      JSON.stringify({
        email,
        token,
        new_password,
      })
    )
      .then(() => {
        email = "";
        token = "";
        new_password = "";
        checkNewPassword = "";
        stage = 1;
        toast.push("your password was successfully reset.");
        togglePasswordModal();
        toggleLogin();
      })
      .catch((err: AxiosError) => {
        if (err.response!.status === 403) {
          toast.push(
            "check your inputs. note that password must be at least 8 characters with 1 numeric digit"
          );
          return;
        }

        toast.push("an unknown error occured.");
      });
  }
</script>

{#if show}
  <div on:click|self={togglePasswordModal} class="reset_password-modal">
    {#if stage === 1}
      <form on:submit|preventDefault={SendResetToken}>
        <h2>{title}</h2>
        <div class="inputs">
          <input type="email" placeholder="Email" bind:value={email} />
        </div>
        <button disabled={disableSendTokenBtn} type="submit"
          >Send verification code</button
        >
      </form>
    {/if}
    {#if stage === 2}
      <form on:submit|preventDefault={ResetPassword}>
        <h2>Reset Password</h2>
        <div class="inputs">
          <input type="email" placeholder="Email" bind:value={email} />
          <input
            type="text"
            placeholder="Verification code"
            bind:value={token}
          />
          <input
            type="password"
            placeholder="Password"
            bind:value={new_password}
          />
          <input
            type="password"
            placeholder="Re-enter password"
            bind:value={checkNewPassword}
          />
        </div>
        <button type="submit">Reset</button>
      </form>
    {/if}
  </div>
{/if}

<style lang="scss">
  .reset_password-modal {
    width: 100%;
    height: 100%;
    background-color: rgba($color: #000000, $alpha: 0.3);
    position: fixed;
    display: flex;
    align-items: center;
    justify-content: center;

    form {
      width: 20rem;
      padding: 1rem;
      border-radius: 0.3rem;
      display: flex;
      align-items: center;
      justify-content: space-between;
      flex-direction: column;
      background-color: #fff;

      .inputs {
        width: 100%;
        padding: 1rem;
        display: flex;
        align-items: center;
        flex-direction: column;
      }

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

      button {
        border: none;
        outline: none;
        padding: 1rem;
        width: 100%;
        background-color: #800080;
        border-radius: 0.3rem;
        color: #fff;
        font-size: 1rem;
      }

      button:hover {
        background-color: #500050;
      }

      button:disabled,
      button[disabled] {
        background-color: #505050;
      }
    }
  }
</style>
