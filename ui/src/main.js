// App entry point

import m from "mithril";
import layout from "./components/layout";
import login from "./components/login";

m.route.prefix("");

m.route(document.body, "/", {
  "/": {
    render: function() {
      return m(layout, m("p", "this is the homepage"));
    }
  },
  "/about": {
    render: function() {
      return m(layout, m("p", "this is the about page"));
    }
  },
  "/privacy-policy": {
    render: function() {
      return m(layout, m("p", "this is privacy policy"));
    }
  },
  "/terms-and-conditions": {
    render: function() {
      return m(layout, m("p", "this is terms and conditions"));
    }
  },
  "/promise": {
    render: function() {
      return m(layout, m("p", "this is the hydrocarbon promise"));
    }
  },
  "/login": {
    render: function() {
      return m(layout, m(login));
    }
  }
});
