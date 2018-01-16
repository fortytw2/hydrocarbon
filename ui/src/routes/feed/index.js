import {h, Component} from 'preact';
import Drawer from 'preact-material-components/Drawer';
import List from 'preact-material-components/List';
import 'preact-material-components/Drawer/style.css';
import 'preact-material-components/List/style.css';

export default class Feed extends Component {
  render(){
    return (
      <div>
		  <h1 class='mdc-typography--display1'> feed content goes here </h1>
      </div>
    );
  }
}