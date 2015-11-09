var React = require('react')

module.exports = React.createClass({
  render() {
    return (
      <header>
        <span>{this.props.name}</span>
      </header>
    );
  }
});
