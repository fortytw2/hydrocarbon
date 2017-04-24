import m from "mithril";

var footerClass = "ph3 pv2 ph5-m w-100 mid-gray absolute bottom-1";
var linkClass = "f6 dib ph2 link mid-gray dim";
var svgLinkClass = "f6 dib ph2 link mid-gray dim bottom-0 lh0";

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
        { class: "tc mt2" },
        m(
          "a",
          {
            class: svgLinkClass,
            href: "https://github.com/fortytw2/hydrocarbon?ref=hc",
            oncreate: m.route.link
          },
          m(
            "svg",
            {
              fill: "currentColor",
              xmlns: "http://www.w3.org/2000/svg",
              width: "16px",
              height: "16px",
              "view-box": "0 0 16 16",
              "fill-rule": "evenodd",
              "clip-rule": "evenodd",
              "stroke-linejoin": "round",
              "stroke-miterlimit": "1.414"
            },
            m("path", {
              d: "M8 0C3.58 0 0 3.582 0 8c0 3.535 2.292 6.533 5.47 7.59.4.075.547-.172.547-.385 0-.19-.007-.693-.01-1.36-2.226.483-2.695-1.073-2.695-1.073-.364-.924-.89-1.17-.89-1.17-.725-.496.056-.486.056-.486.803.056 1.225.824 1.225.824.714 1.223 1.873.87 2.33.665.072-.517.278-.87.507-1.07-1.777-.2-3.644-.888-3.644-3.953 0-.873.31-1.587.823-2.147-.083-.202-.358-1.015.077-2.117 0 0 .672-.215 2.2.82.638-.178 1.323-.266 2.003-.27.68.004 1.364.092 2.003.27 1.527-1.035 2.198-.82 2.198-.82.437 1.102.163 1.915.08 2.117.513.56.823 1.274.823 2.147 0 3.073-1.87 3.75-3.653 3.947.287.246.543.735.543 1.48 0 1.07-.01 1.933-.01 2.195 0 .215.144.463.55.385C13.71 14.53 16 11.534 16 8c0-4.418-3.582-8-8-8"
            })
          )
        ),
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
          { class: linkClass, href: "/promise", oncreate: m.route.link },
          "promise"
        ),
        m(
          "a",
          {
            class: linkClass,
            href: "/privacy-policy",
            oncreate: m.route.link
          },
          "privacy"
        )
      )
    ]);
  }
};
