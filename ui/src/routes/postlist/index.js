import { h, Component } from "preact";
import style from "./style";

export default class PostList extends Component {
  render({ id }, {}) {
    return (
      <div class={style.content}>
        <h1>Feed ID {id}</h1>
      </div>
    );
  }
}
