import { h, Component } from "preact";
import { bind } from "decko";
import { createFeed, listPlugins } from "@/http";

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
  async createFeed(e) {
    e.preventDefault();

    this.setState({ submitting: true });

    let id;
    let title;
    try {
      const resp = await createFeed({
        url: this.state.url,
        folderId: this.props.folderId,
        apiKey: this.props.apiKey
      });
      id = resp.id;
      title = resp.title;
    } catch (e) {
      this.setState({
        submitting: false,
        url: "",
        error: e
      });
      return;
    }

    this.props.onSubmit({
      title: title,
      id: id,
      folderId: this.props.folderId
    });

    this.setState(initialState);
  }

  render({}, { submitting, error, url }) {
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
        <label for="feedURL">URL you would like to add</label>

        <button>Create Feed</button>
      </form>
    );
  }
}
