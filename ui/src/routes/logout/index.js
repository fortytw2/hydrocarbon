import { route } from "preact-router";

export default ({ loginSwapper }, {}) => {
  window.localStorage.clear();
  loginSwapper();
  route("/login");
};
