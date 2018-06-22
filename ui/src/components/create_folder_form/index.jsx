import { h, Component } from "preact";
import { bind } from "decko";
import { createFolder } from "@/http";

import style from "./style.css";
import inputStyle from "@/styles/textbox.css";

const initialState = {
  nameVal: "",
  submitting: false,
  error: null
};

export default class CreateFolderForm extends Component {
  constructor(props) {
    super(props);

    this.setState(initialState);
  }

  @bind
  handleInput(e) {
    this.setState({ nameVal: e.target.value });
  }

  @bind
  async createFolder(e) {
    e.preventDefault();

    this.setState({ submitting: true });

    let id;
    try {
      resp = await createFolder({
        name: this.state.nameVal,
        apiKey: this.props.apiKey
      });
      id = resp.id;
    } catch (e) {
      this.setState({ submitting: false, nameVal: "", error: e });
      return;
    }

    this.props.onSubmit({ name: this.state.nameVal, id: id });
    this.setState(initialState);
  }

  render({}, { submitting, error, nameVal }) {
    if (error) {
      return <h3>{error}</h3>;
    }

    if (submitting) {
      return <h3>Processing...</h3>;
    }

    return (
      <form class={style.folderForm} onSubmit={this.createFolder}>
        <h4>New Folder</h4>

        <input
          id="folderName"
          class={inputStyle.input}
          type="text"
          value={nameVal}
          onChange={this.handleInput}
          placeholder="Folder Name"
        />
        <label for="folderName">idk some rando placeholder text</label>

        <button>Create Folder</button>
      </form>
    );
  }
}
