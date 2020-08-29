import React from 'react';
import './ListTemplate.css';

const ListTemplate = ({children}) => {
  return (
    <main className="list-template">
      <section className="list-wrapper">
        { children }
      </section>
    </main>
  );
};

export default ListTemplate;