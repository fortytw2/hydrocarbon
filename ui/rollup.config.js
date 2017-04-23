// rollup.config.js
import resolve from "rollup-plugin-node-resolve";
import commonjs from "rollup-plugin-commonjs";
import closure from "rollup-plugin-closure-compiler-js";

export default {
  entry: "./main.js",
  format: "iife",

  plugins: [
    resolve({ jsnext: true, main: true }),
    commonjs(),
    closure()
  ]
};
