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

class Rule extends Component {
  state = {
    use: 0,
    rule_alarm_min: 0,
    rule_alarm_max: 0,
  };

  // shouldComponentUpdate(nextProps, nextState) {
  //   return this.props.timestamp !== nextProps.timestamp;
  // }

  render() {
    const { data, classes } = this.props;
    let item;
    if (data) {
      item = (
        <TableRow key>
          <TableCell>Use: {data.Use.toString()}</TableCell>
          <TableCell>AlarmMin: {data.AlarmMin}%</TableCell>
          <TableCell>AlarmMax: {data.AlarmMax}%</TableCell>
          <TableCell><button>Change</button></TableCell>
        </TableRow>
      )
    }

    return (
      <Table className={classes.table} size="small" aria-label="a dense table">
        <TableBody>
          {item}
        </TableBody>
      </Table>
    );
  }
}

export default withStyles(styles)(Rule);
