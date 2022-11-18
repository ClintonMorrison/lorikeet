import React from 'react';

export default class TextField extends React.Component {
  constructor(props) {
    super(props);
    this.ref = React.createRef();
  }

  componentDidMount() {
    setTimeout(() => {
      window.M.updateTextFields();
      if (this.props.autoFocus && this.ref.current) {
        this.ref.current.focus();
      }
    }, 0);
  }

  render() {
    const {
      id,
      label,
      onChange,
      value,
      error,
      type,
      icon,
      className
    } = this.props;

    return (
      <div className={`cp-text-field input-field ${className}`}>
        {icon && (<i className="material-icons prefix">{icon}</i>)}
        <input
          id={id}
          type={type}
          className={error ? 'invalid' : ''}
          autoComplete={this.props.autoComplete}
          value={value}
          onChange={(e) => onChange(e.target.value)}
          ref={this.ref} />
        <label htmlFor={id}>{label}</label>
        <span className="helper-text" data-error={error} />
      </div>
    );
  }
}

TextField.defaultProps = {
  type: 'text'
};
