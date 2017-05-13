import { Component, h, render } from "preact";

import { Store } from "../state/Store";

class Login extends Component {
  render(props, state) {
    return (
      <div class="fl w-100 pa2 h-auto">
        <form class="measure center">
          <fieldset class="ba b--transparent ph0 mh0">
            <div class="mt3">
              <label for="email-address" class="db fw6 lh-copy f6">
                we'll send you a link to login
              </label>
              <input
                placeholder="example@example.com"
                id="email"
                type="email"
                name="email-address"
                class="pa2 input-reset ba bg-transparent w-100"
              />
            </div>
          </fieldset>
          <div>
            <input
              type="submit"
              class="b ph3 pv2 input-reset ba b--black bg-transparent grow pointer f6 dib"
            />
          </div>
        </form>
      </div>
    );
  }
}

export default Login;
