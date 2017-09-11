import { ACTIVATE_API_KEY_FAILURE, ACTIVATE_API_KEY_SUCCESS, REMOVE_API_KEY } from "./types";

const emptyLogin = { apiKey: "", email: "" };

export function login(state = emptyLogin, action) {
  console.log("called");
  switch (action.type) {
    case ACTIVATE_API_KEY_SUCCESS:
      return { apiKey: action.apiKey, email: action.email };
    case ACTIVATE_API_KEY_FAILURE:
      return emptyLogin
    case REMOVE_API_KEY:
      return emptyLogin;
    default:
      return state;
  }
}
