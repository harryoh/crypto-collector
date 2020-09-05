import React, { Component } from 'react';
import './ListTemplate.css';

class ListTemplate extends Component {
  state = {
    now: ""
  }

  async componentDidMount() {
    const clock = () => {
      const tz_kr = 9*60*60;
      const iso = new Date(new Date().getTime()+tz_kr*1000).toISOString();
      this.setState({
        now: iso.slice(-13, -5)
      });
      setTimeout(clock, 1000);
    };

    clock();
  }

  render() {
    const { children } = this.props;

    return (
      <div>
        <div className="clock">
          { this.state.now }
        </div>
        <main className="list-template">
          <section className="list-wrapper">
            { children }
          </section>
        </main>
      </div>
    );
  }
};

export default ListTemplate;
