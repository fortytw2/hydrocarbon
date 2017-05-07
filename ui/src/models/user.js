// src/models/User.js
import m from "mithril";
import config from "../config";
import flash from "./flash";
import raven from "raven-js";

export default {
  activateToken: function(token) {
    return m
      .request({
        method: "POST",
        url: config.URL + "/token/activate",
        withCredentials: true,
        data: {
          token: token
        }
      })
      .then(function(result) {
        window.localStorage.setItem("hydrocarbon-key", result.key);
        flash.flashMessage = "logged in ok";
      });
  },
  requestToken: function(email) {
    return m
      .request({
        method: "POST",
        url: config.URL + "/token/request",
        withCredentials: true,
        data: {
          email: email
        }
      })
      .then(function(result) {
        console.log(result);
        flash.flashMessage = result.note;
      })
      .catch(function(error) {
        raven.captureException(error);
      });
  },
  loggedIn: function() {
    if (window.localStorage.getItem("hydrocarbon-key") !== null) {
      return true;
    }
    return false;
  }
};
