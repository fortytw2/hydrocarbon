import { h, Component } from "preact";
import { bind } from "decko";
import { createFeed } from "@/http";

import style from "./style.css";
import inputStyle from "@/styles/textbox.css";

const initialState = {
  url: "",
  submitting: false,
  error: null
};

export default class CreateFeedForm extends Component {
  constructor(props) {
    super(props);

    this.setState(initialState);
  }

  @bind
  handleURLInput(e) {
    this.setState({ url: e.target.value });
  }

  @bind
  handlePluginInput(e) {
    this.setState({ plugin: e.target.value });
  }

  @bind
  async createFeed(e) {
    e.preventDefault();

    this.setState({ submitting: true });

    let id;
    try {
      resp = await createFeed({
        url: this.state.url,
        plugin: this.state.plugin,
        folderID: this.props.folderID,
        apiKey: this.props.apiKey
      });
      id = resp.id;
    } catch (e) {
      this.setState({
        submitting: false,
        url: "",
        error: e
      });
      return;
    }

    this.props.onSubmit({
      name: this.state.url,
      id: id
    });
    this.setState(initialState);
  }

  render({}, { submitting, error, url, plugin }) {
    if (error) {
      return <h3>{error}</h3>;
    }

    if (submitting) {
      return <h3>Processing...</h3>;
    }

    return (
      <form class={style.feedForm} onSubmit={this.createFeed}>
        <h4>New Feed</h4>

        <input
          id="feedURL"
          class={inputStyle.input}
          type="text"
          value={url}
          onChange={this.handleURLInput}
          placeholder="Feed URL"
        />
        <label for="feedURL">base URL for the plugin</label>

        <input
          id="feedPlugin"
          class={inputStyle.input}
          type="text"
          value={plugin}
          onChange={this.handlePluginInput}
          placeholder="Plugin"
        />
        <label for="feedPlugin">base URL for the plugin</label>

        <button>Create Feed</button>
      </form>
    );
  }
}
