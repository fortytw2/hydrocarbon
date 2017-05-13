export const SEND_LOGIN_TOKEN_SUCCESS = "SEND_LOGIN_TOKEN_SUCCESS";
export const SEND_LOGIN_TOKEN_FAILURE = "SEND_LOGIN_TOKEN_FAILURE";

export function sendLoginTokenSuccess() {
  return { type: SEND_LOGIN_TOKEN_SUCCESS };
}

export function sendLoginTokenFailure(message) {
  return { type: SEND_LOGIN_TOKEN_FAILURE, text: message };
}
