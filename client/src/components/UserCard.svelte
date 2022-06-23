<script lang="ts">
  import AuthAPI from "../API/Auth";
  import { API } from "../constants/API";
  import { navigate } from "svelte-navigator";

  import type { User } from "../models/Auth";

  export let user: User;

  function Logout() {
    AuthAPI.post(API.Routes.Auth.Logout).then(() => {
      navigate("/auth");
    });
  }
</script>

<div class="user-card">
  <h2>Welcome</h2>
  <div class="info">
    <p><b>Username</b>: {user.username}</p>
    <p><b>Email</b>: {user.email}</p>
  </div>
  <button class="btn" on:click={Logout}>Logout</button>
</div>

<style lang="scss">
  .user-card {
    width: 20rem;
    height: 20rem;
    padding: 1rem;
    display: flex;
    align-items: center;
    justify-content: space-between;
    flex-direction: column;
    border-radius: 0.3rem;
    box-shadow: 0 0 50px rgba($color: #000000, $alpha: 0.3);

    p {
      padding: 1rem;
      border-radius: 0.3rem;
      background-color: rgba($color: #800080, $alpha: 0.1);
    }

    button {
      width: 100%;
      padding: 1rem;
      color: #fff;
      background-color: #800080;
      border-radius: 0.3rem;
    }

    button:hover {
      background-color: #500050;
    }
  }

  b {
    color: #800080;
  }
</style>
