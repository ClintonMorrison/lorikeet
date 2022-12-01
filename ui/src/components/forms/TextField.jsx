import _ from 'lodash';
import React from 'react';

const MASK = '•'

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
      masked,
      icon,
      className,
    } = this.props;

    return (
      <div className={`cp-text-field input-field ${className}`}>
        {icon && (<i className="material-icons prefix">{icon}</i>)}
        <input
          id={id}
          type={type}
          className={error ? 'invalid' : ''}
          autoComplete={this.props.autoComplete}
          value={masked ? _.repeat(MASK, value.length) : value}
          onChange={(e) => {
            // Below code is based on https://github.com/karaggeorge/react-better-password/blob/master/src/index.js
            if (!masked) {
              onChange(e.target.value)
              return
            }

            // Save the current cursor position to restore after masking
            const input = this.ref.current;
            if (!input) {
              return;
            }

            const cursorPos = input.selectionEnd;

            // Partially masked value from event
            const newRawValue = e.target.value;

            // This is going to be the new original value (unmasked)
            const newValue = newRawValue.replace(new RegExp(`${cursorPos ? `(^\\${MASK}{1,${cursorPos}})|` : ''}(\\${MASK}+)`, 'g'), (match, part, offset, string) => {
              if (!offset && cursorPos) return value.substr(0, match.length);
              else return value.substr(-match.length);
            });


            // Restore cursor position
            setTimeout(() => {
              input.selectionStart = cursorPos;
              input.selectionEnd = cursorPos;
            }, 0);

            onChange(newValue);
          }}
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
