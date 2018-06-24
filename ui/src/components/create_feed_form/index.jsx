import { h, Component } from "preact";
import { bind } from "decko";
import { createFeed, listPlugins } from "@/http";

import style from "./style.css";
import inputStyle from "@/styles/textbox.css";

const initialState = {
  url: "",
  loadingPlugins: true,
  plugins: [],
  plugin: "",
  submitting: false,
  error: null
};

export default class CreateFeedForm extends Component {
  constructor(props) {
    super(props);

    this.setState(initialState);
  }

  async componentDidMount(props) {
    try {
      const plugins = await listPlugins({ apiKey: this.props.apiKey });
      this.setState({ plugins: plugins, loadingPlugins: false });
    } catch (e) {
      this.setState({ error: e, loadingPlugins: false });
    }
  }

  @bind
  reset(e) {
    e.preventDefault();
    this.setState(initialState);
    this.componentDidMount(this.props);
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
      const resp = await createFeed({
        url: this.state.url,
        plugin: this.state.plugin,
        folderId: this.props.folderId,
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
      title: this.state.url,
      id: id,
      folderId: this.props.folderId
    });

    this.setState(initialState);
  }

  render({}, { submitting, error, url, plugin, plugins, pluginsLoading }) {
    if (error) {
      return (
        <div>
          <h3>{error}</h3>
          <button onClick={this.reset}>Reset</button>
        </div>
      );
    }

    if (submitting) {
      return <h3>Processing...</h3>;
    }

    if (pluginsLoading) {
      return <h3>Loading Plugins...</h3>;
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

        <select
          id="feedPlugin"
          onChange={this.handlePluginInput}
          value={plugin}
        >
          {plugins.map(p => <option value={p}>{p}</option>)}
        </select>

        <label for="feedPlugin">actual plugin to use</label>

        <button>Create Feed</button>
      </form>
    );
  }
}
