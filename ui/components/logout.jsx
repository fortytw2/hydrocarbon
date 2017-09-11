import { Component, h } from "preact";

class Logout extends Component {
  constructor() {
    super();
    this.handleClick = this.handleClick.bind(this);
  }

  handleClick(event) {
    event.preventDefault();
    console.log(this.props.apiKey);
    // RequestLoginToken(this.state.emailValue);
  }

  render(props, state) {
    return (
      <a
        href="/"
        id="logout"
        class={this.props.class}
        onClick={this.handleClick}
      >
        logout
      </a>
    );
  }
}

export default Logout;
