import type { Auth } from "../models/Auth";
import { writable } from "svelte/store";

export const AuthState = writable<Auth>({
  username: null,
  email: null,
  accessToken: null,
  idToken: null,
});
