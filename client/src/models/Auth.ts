export interface UserLogin {
  email: string;
  password: string;
}

export interface UserRegister {
  username: string;
  email: string;
  password: string;
}

export interface Auth {
  username: string | null;
  email: string | null;
  accessToken: string | null;
  idToken: string | null;
}

export interface User {
  username: string;
  email: string;
}
