import { NOTIFICATION_ADD, NOTIFICATION_REMOVE } from "./types";

import _ from "lodash";

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
