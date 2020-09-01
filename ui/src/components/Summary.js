import React, { Component } from 'react';
import './Item.css';
import TableCell from '@material-ui/core/TableCell';
import TableRow from '@material-ui/core/TableRow';

const fixrate = 1200;
let _data;

class Summary extends Component {
  state = {
    upbit_premium_fix: 0,
    upbit_premium_curr: 0,
    bithumb_premium_fix: 0,
    bithumb_premium_curr: 0,
  };

  componentDidUpdate(prevProps, prevState) {
    const { data } = this.props;
    if (data === "undefined" || data === null || Object.keys(data).length === 0) {
      return
    }

    delete data.CreatedAt;
    if (JSON.stringify(data) === JSON.stringify(_data)) {
      return;
    }

    const getPremium = (src, desc, rate) => {
      return parseFloat(Math.floor(((src-(desc*rate))/src * 100)*1000)/1000).toFixed(3);
    }

    this.setState({
      upbit_premium_fix: getPremium(data.UpbitPrice.Price, data.BybitPrice.Price, fixrate),
      upbit_premium_curr: getPremium(data.UpbitPrice.Price, data.BybitPrice.Price, data.UsdKrw.Price),
      bithumb_premium_fix: getPremium(data.BithumbPrice.Price, data.BybitPrice.Price, fixrate),
      bithumb_premium_curr: getPremium(data.BithumbPrice.Price, data.BybitPrice.Price, data.UsdKrw.Price)
    });
    _data = data;

    const title=`upbit:${this.state.upbit_premium_fix}% bithumb:${this.state.bithumb_premium_fix}`;
    document.title = title;
  }

  render() {
    return (
      <TableRow>
        <TableCell component="th" scope="row">
          <strong>Premium</strong><br />(Fix:1200)
        </TableCell>
        <TableCell align="right">
          <strong>Upbit</strong><br/>
          Fix: {this.state.upbit_premium_fix} %<br/>
          Cur: {this.state.upbit_premium_curr} %<br/>
        </TableCell>
        <TableCell align="right">
          <strong>Bithumb</strong><br/>
          Fix: {this.state.bithumb_premium_fix} %<br/>
          Cur: {this.state.bithumb_premium_curr} %<br/>
        </TableCell>
      </TableRow>
    );
  }
}

export default Summary;