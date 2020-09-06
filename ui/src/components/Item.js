import React, { Component } from 'react';

import './Item.css';
import { withStyles } from '@material-ui/core/styles';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableRow from '@material-ui/core/TableRow';

const styles = theme => ({
  table: {
    minWidth: 320
  }
});

class Item extends Component {
  state = {
    premium_fix: 0,
    premium_cur: 0,
  };

  shouldComponentUpdate(nextProps, nextState) {
    return this.props.timestamp !== nextProps.timestamp;
  }

  render() {
    const { id, name, price, data, classes } = this.props;

    const numberWithCommas = (x) => {
      return x.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
    }
    const toDateStr = (timestamp) => {
      const tz_kr = 9*60*60;
      const iso = new Date((Number(timestamp)+tz_kr)*1000).toISOString();
      return iso.slice(-13, -5)
    }

    const timeCheck = (timestamp) => {
      const now = Math.floor(+ new Date() / 1000);
      return (timestamp + 20 <= now) ? "warning":"";
    }

    const getPremium = (symbol, src, desc, rate) => {
      const bybit = desc.reduce((n, p) => {
        n[p.Symbol] = {}
        n[p.Symbol]["Price"] = p.Price;
        n[p.Symbol]["Timestamp"] = p.Timestamp;
        return n;
      }, {});
      return parseFloat((src-bybit[symbol].Price*rate)/src * 100).toFixed(3);
    }

    let items;
    let items_name;
    switch (name.toUpperCase()) {
      case "CURRENCY":
        items = price.map(
          ({Symbol, Price, Timestamp}) => (
            <TableRow key={Symbol+'-'+Timestamp}>
              <TableCell component="th" scope="row"><strong>{Symbol}</strong></TableCell>
              <TableCell align="right">Fix: 1,200</TableCell>
              <TableCell align="right">Cur: {numberWithCommas(Price)}</TableCell>
              <TableCell align="right">{toDateStr(Timestamp)}</TableCell>
            </TableRow>
          )
        );
        break;
      case "BYBIT":
        let bybit = price.reduce((n, p) => {
          n[p.Symbol] = {}
          n[p.Symbol]["Price"] = p.Price;
          n[p.Symbol]["Timestamp"] = p.Timestamp;
          return n;
        }, {});
        items = (
          <TableRow className={timeCheck(this.props.timestamp)}>
            <TableCell align="right"><strong>{name.toUpperCase()}</strong></TableCell>
            <TableCell align="right">BTC: {numberWithCommas(bybit.BTC.Price.toFixed(1))}</TableCell>
            <TableCell align="right">ETH: {numberWithCommas(bybit.ETH.Price.toFixed(1))}</TableCell>
            <TableCell align="right">XRP: {bybit.XRP.Price.toFixed(4)}</TableCell>
            <TableCell align="right">{toDateStr(this.props.timestamp)}</TableCell>
          </TableRow>
        )
        break;
      default:
        items_name = (
          <TableRow>
            <TableCell align="center" colSpan={5}><strong>{name.toUpperCase()}</strong></TableCell>
          </TableRow>
        );

        for (let p of price) {
          p["PremiumFix"] = getPremium(p.Symbol, p.Price, data.BybitPrice.Price, 1200);
          p["PremiumCur"] = getPremium(p.Symbol, p.Price, data.BybitPrice.Price, data.Currency.Price[0].Price);
        }

        items = price.map(
          ({Symbol, Price, Timestamp, PremiumFix, PremiumCur}) => (
            <TableRow key={Symbol+'-'+Timestamp} className={timeCheck(Timestamp)}>
              <TableCell component="th" scope="row"><strong>{Symbol}</strong></TableCell>
              <TableCell align="right">{numberWithCommas(Price)}</TableCell>
              <TableCell align="right">Fix: {PremiumFix}%</TableCell>
              <TableCell align="right">Cur: {PremiumCur}%</TableCell>
              <TableCell align="right">{toDateStr(Timestamp)}</TableCell>
            </TableRow>
          )
        );
        break;
    }

    return (
      <Table key={id} className={classes.table} size="small" aria-label="a dense table">
        <TableBody>
          {items_name}
          {items}
        </TableBody>
      </Table>
    );
  }
}

export default withStyles(styles)(Item);