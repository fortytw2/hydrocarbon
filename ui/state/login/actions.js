import { ACTIVATE_API_KEY } from "./types";

export function activateApiKey(email, apiKey) {
  return { type: ACTIVATE_API_KEY, email: email, apiKey: apiKey };
}
