import m from "mithril";
import nav from "./nav";
import footer from "./footer";
import user from "../models/user";

var inputClass = "pa2 input-reset ba bg-transparent w-100";

export default {
  view(vnode) {
    return m(
      "form",
      {
        class: "measure center",
        onsubmit: function(e) {
          e.preventDefault()
          user.requestToken(document.getElementById("email").value)
        }
      },
      [
        m(
          "fieldset",
          { class: "ba b--transparent ph0 mh0" },
          m("div", { class: "mt3" }, [
            m(
              "label",
              {
                class: "db fw6 lh-copy f6",
                for: "email-address"
              },
              "we'll send you a link to login"
            ),
            m("input", {
              class: inputClass,
              placeholder: "example@example.com",
              id: "email",
              type: "email",
              name: "email-address"
            })
          ])
        ),
        m(
          "div",
          m(
            "input",
            {
              class: "b ph3 pv2 input-reset ba b--black bg-transparent grow pointer f6 dib",
              type: "submit"
            },
            "send link"
          )
        )
      ]
    );
  }
};
