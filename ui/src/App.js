import React, { Component } from 'react';
import ListTemplate from './components/ListTemplate';
import ItemList from './components/ItemList';

class App extends Component {
  // state = {
  //   UpbitPrice: {
  //     Name: "upbit",
  //     Symbol: "KRW-BTC",
  //     Price: 13467000,
  //     Timestamp: 1598706857
  //   },
  //   BithumbPrice: {
  //     Name: "bithumb",
  //     Symbol: "BTC_KRW",
  //     Price: 13448000,
  //     Timestamp: 1598706858
  //   },
  //   BybitPrice: {
  //     Name: "bybit",
  //     Symbol: "BTCUSD",
  //     Price: 11469,
  //     Timestamp: 1598706858
  //   },
  //   UsdKrw: {
  //     Name: "usdkrw",
  //     Symbol: "USDKRW",
  //     Price: 1180.48,
  //     Timestamp: 1598706268
  //   },
  //   CreatedAt: 1598706861
  // }
  state = {
    data: [
      { id: 0, name: 'upbit', price: '13467000', timestamp: '1598715994'},
      { id: 1, name: 'bithumb', price: '13448000', timestamp: '1598706858'},
      { id: 2, name: 'bybit', price: '11469', timestamp: '1598706858'},
      { id: 3, name: 'usdkrw', price: '1180.48', timestamp: '1598706268'},
    ]
  }
  render() {
    const { data } = this.state;
    return (
      <ListTemplate>
        <ItemList data={data}/>
      </ListTemplate>
    );
  }
}

export default App;
