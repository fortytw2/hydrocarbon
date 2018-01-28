import { h, Component } from "preact";
import Button from "preact-material-components/Button";
import "preact-material-components/Button/style.css";
import "preact-material-components/Theme/style.css";
import { route } from "preact-router";
import LinearProgress from "preact-material-components/LinearProgress";
import "preact-material-components/LinearProgress/style.css";

export default class Login extends Component {
  constructor(props) {
    super(props);
    this.setState({
      email: "",
      success: {
        submitted: false,
        status: "",
        message: ""
      }
    });

    this.update = this.update.bind(this);
    this.submit = this.submit.bind(this);
  }

  componentDidMount() {
    if (this.props.callback) {
      let params = new URLSearchParams(location.search.slice(1));
      let token = params.get("token");

      fetch("/v1/key/create", {
        method: "POST",
        body: JSON.stringify({
          token: token
        })
      })
        .then(response => {
          if (response.ok) {
            return response.json();
          }
        })
        .then(json => {
          window.localStorage.setItem("hydrocarbon-key", json.key);
          window.localStorage.setItem("email", json.email);

          route("/folders");
        });
    }
  }

  update(e) {
    e.preventDefault();

    let email = e.target.value;
    this.setState({ email });
  }

  submit = e => {
    e.preventDefault();

    fetch("/v1/token/create", {
      method: "POST",
      body: JSON.stringify({
        email: this.state.email
      })
    })
      .then(response => {
        if (response.ok) {
          return response.json();
        }
      })
      .then(json => {
        this.setState({
          success: {
            submitted: true,
            message: json.note
          }
        });
      });
  };

  render({ callback }, { success, email }) {
    if (callback) {
      return (
        <div>
          <h4> Logging you in... </h4>
        </div>
      );
    }

    if (success.submitted) {
      return (
        <div>
          <h4>{success.message}</h4>
        </div>
      );
    }

    return (
      <div>
        <label>
          <input type="email" value={email} onChange={this.update} />
        </label>
        <Button ripple raised onClick={this.submit}>
          Login (or Register)
        </Button>
      </div>
    );
  }
}
