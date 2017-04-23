import m from "mithril";
import * as Mithril from "mithril";

var footerClass = "pv4 ph3 ph5-m ph6-l mid-gray bottom-0";
var linkClass = "f6 dib ph2 link mid-gray dim";

export default {
  view(vnode) {
    return m("footer", { class: footerClass }, [
      m(
        "small",
        { class: "f6 db tc" },
        "© 2017 ",
        m("b", { class: "b" }, "hydrocarbon")
      ),
      m(
        "div",
        { class: "tc mt3" },
        m(
          "a",
          {
            class: linkClass,
            href: "/terms-and-conditions",
            oncreate: m.route.link
          },
          "terms and conditions"
        ),
        m(
          "a",
          { class: linkClass, href: "/privacy-policy", oncreate: m.route.link },
          "promise"
        ),
        m(
          "a",
          { class: linkClass, href: "/privacy-policy", oncreate: m.route.link },
          "privacy"
        )
      )
    ]);
  }
};
