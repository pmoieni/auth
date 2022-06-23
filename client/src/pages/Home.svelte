<script lang="ts">
  import type { AxiosError } from "axios";
  import UserCard from "../components/UserCard.svelte";
  import AccountAPI from "../API/Account";

  import { API } from "../constants/API";

  import type { User } from "../models/Auth";
  import { onMount } from "svelte";

  let userInfo: User = {
    username: "",
    email: "",
  };
  onMount(() => {
    AccountAPI.get(API.Routes.Account.Me).then((res) => {
      userInfo = {
        username: res.data["username"],
        email: res.data["email"],
      };
    });
  });
</script>

<div class="home-page">
  <UserCard user={userInfo} />
</div>

<style lang="scss">
  .home-page {
    width: 100vw;
    height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
  }
</style>
