import React, { Component } from 'react';
import Item from './Item';
import Summary from './Summary';
import axios from "axios";
import { withStyles } from '@material-ui/core/styles';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableContainer from '@material-ui/core/TableContainer';
import Paper from '@material-ui/core/Paper';

const styles = theme => ({
  table: {
    minWidth: 320
  }
});

const parse_data = (data) => {
  return Object.keys(data)
    .filter(k => typeof(data[k])==="object")
    .map((k, index) => {
      return {Id: index, ...data[k]}
    }
  )
}

class ItemList extends Component {
  state = {
    response: {},
    data: [],
  }
  async componentDidMount() {
    const baseurl = (process.env.NODE_ENV === "development") ? "http://localhost:8080":"";
    const getPrices = async () => {
      try {
        const res = await axios.get(`${baseurl}/api/prices`);
        const data = parse_data(res.data);
        this.setState({
          response: res.data,
          data: data
        })
      } catch (error) {
        console.error(error);
      } finally {
        setTimeout(getPrices, 1000 * 1);
      }
    };
    getPrices();
  }

  render() {
    const { data, response } = this.state;
    const { classes } = this.props;
    const prices = data.map(
      ({Id, Name, Price, Timestamp}) => (
        <Item
          id={Id}
          key={Id}
          name={Name}
          price={Price}
          timestamp={Timestamp}
        />
      )
    );

    return (
      <TableContainer component={Paper}>
        <Table className={classes.table} size="small" aria-label="a dense table">
          <TableBody>
            {prices}
            <Summary data={response}/>
          </TableBody>
        </Table>
      </TableContainer>
    );
  }
}

export default withStyles(styles)(ItemList);
