import m from "mithril";
import * as Mithril from "mithril";
import nav from "./nav";
import footer from "./footer";

export default {
  view(vnode) {
    return m(".page", [
      m(nav),
      m("h1", "login"),
      m("p", "login page will go here"),
      m(footer)
    ]);
  }
};
