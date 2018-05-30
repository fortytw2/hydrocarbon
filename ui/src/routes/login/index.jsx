import { h, Component } from "preact";
import { route } from "preact-router";
import { bind } from "decko";
import { requestToken } from "@/http";

import style from "./style.css";
import textBox from "@/styles/textbox.css";

const initialState = {
  email: "",
  presubmitError: null,
  success: {
    error: null,
    submitted: false,
    message: ""
  }
};

export default class Login extends Component {
  constructor(props) {
    super(props);

    this.setState(initialState);
  }

  @bind
  update(e) {
    e.preventDefault();

    let email = e.target.value;
    this.setState({ email });
  }

  @bind
  async submit(e) {
    e.preventDefault();

    var re = /\S+@\S+\.\S+/;
    if (!re.test(this.state.email)) {
      this.setState({
        presubmitError: "invalid email"
      });
      return;
    }

    try {
      const json = await requestToken({
        email: this.state.email
      });

      this.setState({
        success: {
          submitted: true,
          message: json.data
        }
      });
    } catch (e) {
      this.setState({
        success: {
          error: e
        }
      });
    }
  }

  @bind
  getPresubmitError() {
    if (this.state.presubmitError) {
      return <p class={style.labelText}>{this.state.presubmitError}</p>;
    }

    return null;
  }

  render({}, { email, success }) {
    if (success.error) {
      return (
        <div class={style.loginArea}>
          <div class={style.loginBox}>
            <div class={style.formOffset}>
              <h3>Could not email link</h3>
              <p class={style.labelText}>{success.error}</p>
              <div class={style.buttonBox}>
                <button
                  class={style.submitButton}
                  onClick={() => {
                    this.setState(initialState);
                  }}
                >
                  Try Again
                </button>
              </div>
            </div>
          </div>
        </div>
      );
    }

    if (success.submitted) {
      return (
        <div class={style.loginArea}>
          <div class={style.loginBox}>
            <div class={style.notifOffset}>
              <h3>Login email sent</h3>
              <p class={style.labelText}>{success.message}</p>
            </div>
          </div>
        </div>
      );
    }

    return (
      <div class={style.loginArea}>
        <div class={style.loginBox}>
          <div class={style.formOffset}>
            <h3>Login</h3>
            <p class={style.labelText}>
              Hydrocarbon sends an email with a link to login. There are no
              passwords to generate or remember.
            </p>
            <div class={style.loginInput}>
              <input
                class={textBox.input}
                type="email"
                placeholder="example@example.com"
                value={email}
                onChange={this.update}
              />
            </div>
            {this.getPresubmitError()}
            <div class={style.buttonBox}>
              <button class={style.submitButton} onClick={this.submit}>
                Login (or Register)
              </button>
            </div>
          </div>
        </div>
      </div>
    );
  }
}
