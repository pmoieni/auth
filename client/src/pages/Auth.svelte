<script lang="ts">
  import ResetPassword from "../components/ResetPassword.svelte";
  import LoginForm from "../components/LoginForm.svelte";
  import RegisterForm from "../components/RegisterForm.svelte";

  let isLoginOpen = true;
  let isResetPasswordOpen = false;

  function toggleAuth() {
    isLoginOpen = !isLoginOpen;
  }

  function toggleLogin() {
    isLoginOpen = true;
  }

  function togglePasswordModal() {
    isResetPasswordOpen = !isResetPasswordOpen;
  }
</script>

<div class="auth-page">
  <RegisterForm show={!isLoginOpen} toggleLogin={toggleAuth} />
  <LoginForm
    show={isLoginOpen}
    toggleRegister={toggleAuth}
    {togglePasswordModal}
  />
  <ResetPassword
    show={isResetPasswordOpen}
    {toggleLogin}
    {togglePasswordModal}
  />
</div>

<style lang="scss">
  .auth-page {
    widows: 100vw;
    height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background-image: url("../assets/cool-background.svg");
    background-position: center;
    background-repeat: no-repeat;
    background-size: cover;
    background-attachment: fixed;
  }
</style>
