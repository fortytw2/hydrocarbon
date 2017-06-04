import { ACTIVATE_API_KEY, REMOVE_API_KEY } from "./UserActions";
import { NOTIFICATION_ADD, NOTIFICATION_REMOVE } from "./Notifications";

import _ from "lodash";

const emptyLogin = { apiKey: "", email: "" };

export function login(state = emptyLogin, action) {
  switch (action.type) {
    case ACTIVATE_API_KEY:
      return { apiKey: action.apiKey, email: action.email };
    case REMOVE_API_KEY:
      return emptyLogin;
    default:
      return state;
  }
}

export function notifications(state = [], action) {
  switch (action.type) {
    case NOTIFICATION_ADD:
      return _.concat(state, [
        {
          key: action.key,
          message: action.message,
          level: action.level
        }
      ]);
    case NOTIFICATION_REMOVE:
      return _.filter(state, n => {
        return n.key !== action.key;
      });
    default:
      return state;
  }
}
