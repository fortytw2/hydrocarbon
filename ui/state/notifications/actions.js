import { NOTIFICATION_ADD, NOTIFICATION_REMOVE } from "./types";

export function addNotification(level, message) {
  return {
    type: NOTIFICATION_ADD,
    key: (Math.random() + 1).toString(36),
    level: level,
    message: message
  };
}

export function removeNotification(key) {
  return { type: NOTIFICATION_REMOVE, key: key };
}
