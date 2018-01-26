import { h, Component } from "preact";

export default class PostList extends Component {
  render({ id }, {}) {
    return (
      <div>
        <h1>Feed ID {id}</h1>
      </div>
    );
  }
}
