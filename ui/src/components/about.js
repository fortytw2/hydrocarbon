import m from "mithril";
import * as Mithril from "mithril";
import nav from "./nav";
import footer from "./footer";

export default {
  view(vnode) {
    return m(".page", [
      m(nav),
      m("h1", "about"),
      m("p", "this is the about page."),
      m(footer)
    ]);
  }
};
