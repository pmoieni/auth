<script lang="ts">
  import { SvelteToast } from "@zerodevx/svelte-toast";

  import { Route, Router } from "svelte-navigator";
  import Auth from "./pages/Auth.svelte";
  import Home from "./pages/Home.svelte";
</script>

<main>
  <SvelteToast />
  <Router>
    <Route path="/" component={Home} />
    <Route path="auth" component={Auth} />
  </Router>
</main>

<style>
  @font-face {
    font-family: raleway;
    src: url(./assets/fonts/Raleway-Regular.ttf);
  }

  :global(body) {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
    font-family: raleway;
  }

  :global(.btn) {
    border: none;
    outline: none;
    background: none;
    font-size: 1.2rem;
  }
</style>
