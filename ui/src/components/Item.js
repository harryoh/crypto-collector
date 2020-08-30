import React, { Component } from 'react';
import './Item.css';
import TableCell from '@material-ui/core/TableCell';
import TableRow from '@material-ui/core/TableRow';

class Item extends Component {
  render() {
    const { id, name, price, timestamp } = this.props;
    const numberWithCommas = (x) => {
      return x.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
    }
    const toDateStr = (timestamp) => {
      const tz_kr = 9*60*60;
      const iso = new Date((Number(timestamp)+tz_kr)*1000).toISOString();
      return `${iso.slice(0, 10)} ${iso.slice(-13, -5)}`
    }

    return (
      <TableRow key={id}>
        <TableCell component="th" scope="row"><strong>{name.toUpperCase()}</strong></TableCell>
        <TableCell align="right">{numberWithCommas(price)}</TableCell>
        <TableCell align="right">{toDateStr(timestamp)}</TableCell>
      </TableRow>
    );
  }
}

export default Item;