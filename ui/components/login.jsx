import { Component, h } from "preact";

import { RequestLoginToken } from "../http/user";

class Login extends Component {
  constructor(props) {
    super(props);
    this.setState({ emailValue: "" });
  }

  handleChange(event) {
    this.setState({ emailValue: event.target.value });
  }

  handleSubmit(event) {
    event.preventDefault();
    RequestLoginToken(this.state.emailValue);
  }

  render(props, state) {
    return (
      <div class="fl w-100 pa5 h-auto">
        <form class="measure center" onSubmit={this.handleSubmit.bind(this)}>
          <fieldset class="ba b--transparent ph0 mh0">
            <div class="mt3">
              <label for="email-address" class="db fw6 lh-copy f6">
                we'll send you a link to login
              </label>
              <input
                value={state.emailValue}
                onChange={this.handleChange.bind(this)}
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
              value="Submit"
            />
          </div>
        </form>
      </div>
    );
  }
}

export default Login;
