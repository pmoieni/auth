export const API = {
  Routes: {
    Base: "http://localhost:8080",
    Auth: {
      Login: "http://localhost:8080/api/v1/auth/direct/login",
      Register: "http://localhost:8080/api/v1/auth/direct/register",
      Logout: "http://localhost:8080/api/v1/auth/direct/logout",
      Token: "http://localhost:8080/api/v1/auth/direct/token",
    },
    Account: {
      Me: "http://localhost:8080/api/v1/account/me/",
      ResetPasswordToken: "http://localhost:8080/api/v1/account/password/token",
      ResetPassword: "http://localhost:8080/api/v1/account/password/reset",
    },
  },
};
