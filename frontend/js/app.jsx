var React = require('react');
var ReactDOM = require('react-dom')
var Header = require('./components/header.jsx')

let App = React.createClass({
  render() {
    return (
      <Header name="fueleconomy.io" />
    );
  }
});

ReactDOM.render(<App />, document.body)
