import { Component, h } from "preact";

class TextContent extends Component {
  render({ text = "" }, {}) {
    return (
      <section class="pa1 pa2-ns">
        <p>{text}</p>
      </section>
    );
  }
}

export default TextContent;
