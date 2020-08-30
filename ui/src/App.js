import React, { Component } from 'react';
import ListTemplate from './components/ListTemplate';
import ItemList from './components/ItemList';

class App extends Component {
  render() {
    return (
      <ListTemplate>
        <ItemList/>
      </ListTemplate>
    );
  }
}

export default App;
