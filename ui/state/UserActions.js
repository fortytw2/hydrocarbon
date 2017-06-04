export const ACTIVATE_API_KEY = "ACTIVATE_API_KEY";
export const REMOVE_API_KEY = "REMOVE_API_KEY";

export function activateApiKey(email, apiKey) {
  return { type: ACTIVATE_API_KEY, email: email, apiKey: apiKey };
}
