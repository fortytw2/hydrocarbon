export const NOTIFICATION_ADD = "NOTIFICATION_ADD";
export const NOTIFICATION_REMOVE = "NOTIFICATION_REMOVE";

export const NOTIFICATION_LEVEL_INFO = "NOTIFICATION_LEVEL_INFO";
export const NOTIFICATION_LEVEL_WARNING = "NOTIFICATION_LEVEL_WARNING";

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
