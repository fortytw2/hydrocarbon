import { Component, h, render } from "preact";

export default ({ type, url }, {}) => (
  <section class="pa1 pa2-ns">
    <h3>404 - Page not Found</h3>
    <p>It looks like we hit a snag.</p>
    <pre>{url}</pre>
  </section>
);
