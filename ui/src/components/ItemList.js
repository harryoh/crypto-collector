import React, { Component } from 'react';
import Item from './Item';
import Rule from './Rule';
import axios from "axios";

import TableContainer from '@material-ui/core/TableContainer';
import Paper from '@material-ui/core/Paper';

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
        const response = res.data?.data
        const data = parse_data(response);
        this.setState({
          response: response,
          data: data,
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
    const items = data.filter(k => k.Name && k.Name !== "currency")
      .map(({Id, Name, Price, Timestamp}) => (
        <Item
          id={Id}
          key={Id}
          name={Name}
          price={Price}
          data={response}
          timestamp={Timestamp}
        />
      )
    );
    return (
        <TableContainer component={Paper}>
          <Rule
            data={response['Rule']}
          />
          {items}
        </TableContainer>
    );
  }
}

export default ItemList;
