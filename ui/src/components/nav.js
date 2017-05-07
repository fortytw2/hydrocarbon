import m from "mithril";
import user from "../models/user";

var navClass = "pa3 pa4-ns";
var titleClass = "link dim black b f6 f5-ns dib mr3";
var linkClass = "link dim gray    f6 f5-ns dib mr3";

export default {
  view(vnode) {
    return m(
      "nav",
      { class: navClass },
      m(
        "a",
        { class: titleClass, href: "/", oncreate: m.route.link },
        "hydrocarbon"
      ),
      user.loggedIn()
        ? m(
            "a",
            {
              class: linkClass,
              href: "/settings",
              oncreate: m.route.link
            },
            "settings"
          )
        : m(
            "a",
            { class: linkClass, href: "/login", oncreate: m.route.link },
            "login"
          )
    );
  }
};
