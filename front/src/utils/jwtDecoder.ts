import {jwtDecode} from 'jwt-decode';

interface JwtPayloadType {
  username?: string;
  email?: string;
  role?: string;
  id?: number;
}

export const getUsernameFromToken = () => {
  const authToken = localStorage.getItem("authToken");
  if (authToken) {
    const decoded = jwtDecode<JwtPayloadType>(authToken);
    return decoded?.username;
  }
  return null;
};

export const getEmailFromToken = () => {
  const authToken = localStorage.getItem("authToken");
  if (authToken) {
    const decoded = jwtDecode<JwtPayloadType>(authToken);
    return decoded?.email;
  }
  return null;
};

export const getUserIdFromToken = () => {
  const authToken = localStorage.getItem("authToken");
  if (authToken) {
    const decoded = jwtDecode<JwtPayloadType>(authToken);
    return decoded?.id;
  }
  return null;
};

export const getRoleFromToken = () => {
  const authToken = localStorage.getItem("authToken");
  if (authToken) {
    const decoded = jwtDecode<JwtPayloadType>(authToken);
    return decoded?.role;
  }
  return null;
};
