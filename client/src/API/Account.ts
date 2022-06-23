import axios, { Axios, AxiosRequestConfig } from "axios";
import { API } from "../constants/API";
import { AuthState } from "../store/Auth";
import { navigate } from "svelte-navigator";
import AuthAPI from "./Auth";
import { toast } from "@zerodevx/svelte-toast";

const AccountAPI = axios.create({
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
  },
});

declare module "axios" {
  export interface AxiosRequestConfig {
    __isRetryRequest?: boolean;
  }
}

AccountAPI.interceptors.request.use((config: AxiosRequestConfig) => {
  let accessToken;
  AuthState.subscribe((value) => (accessToken = value.accessToken));

  if (config.headers === undefined) {
    config.headers = {};
  }

  config.headers["Authorization"] = `${accessToken}`;

  return config;
});

AccountAPI.interceptors.response.use(
  (response) => {
    return response;
  },
  async (err) => {
    const originalConfig = err.config;
    if (axios.isAxiosError(err)) {
      if (err.response!.status === 401 && !originalConfig.__isRetryRequest) {
        originalConfig.__isRetryRequest = true;
        try {
          const res = await AuthAPI.get(API.Routes.Auth.Token);
          AuthState.update((value) => {
            value.accessToken = res.data["access_token"];
            return value;
          });
        } catch (err) {
          navigate("/auth");
        }

        return AccountAPI(originalConfig);
      }
      // if (err.response!.status === 401 && originalConfig.__isRetryRequest) {
      //   window.location.href = "/auth";
      // }
      if (err.response!.status === 500) {
        toast.push("An unknown error occured.");
      }
    }

    return Promise.reject(err);
  }
);

export default AccountAPI;
