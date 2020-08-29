import React, { Component } from 'react';
import Item from './Item';
import { withStyles } from '@material-ui/core/styles';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableContainer from '@material-ui/core/TableContainer';
import Paper from '@material-ui/core/Paper';

const styles = theme => ({
  root: {
    width: "100%",
    marginTop: theme.spacing(3),
    overflowX: "auto",
  },
  table: {
//    minWidth: 1080
  }
});

class ItemList extends Component {
  render() {
    const { data, classes } = this.props;
    const prices = data.map(
      ({id, name, price, timestamp}) => (
        <Item
          id={id}
          key={id}
          name={name}
          price={price}
          timestamp={timestamp}
        />
      )
    );

    return (
      <TableContainer component={Paper}>
        <Table className={classes.table} size="small" aria-label="a dense table">
          <TableBody>
            {prices}
          </TableBody>
        </Table>
      </TableContainer>
    );
  }
}

export default withStyles(styles)(ItemList);
