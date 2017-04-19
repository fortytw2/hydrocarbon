// rollup.config.js
import typescript from 'rollup-plugin-typescript';
import resolve from 'rollup-plugin-node-resolve';
import commonjs from 'rollup-plugin-commonjs';
import closure from 'rollup-plugin-closure-compiler-js';
import typescript2 from 'typescript';

export default {
  entry: './main.ts',
  format: 'iife',

  plugins: [
    resolve({ jsnext: true, main: true }),
    commonjs(),
    typescript({
        typescript: typescript2
    }),
    closure()
  ]
}
