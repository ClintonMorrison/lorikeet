import React from 'react';

export default class Modal extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      passwordLength: 12,
      includeNumbers: true,
      includeSymbols: true,
    };

    this.handleGenerate = this.handleGenerate.bind(this);
  }

  componentDidMount() {
    window.M.Modal.init(document.querySelectorAll('.modal'), {
      onOpenEnd: () => { }
    });
    window.M.Range.init(document.querySelectorAll("input[type=range]"));
  }

  handleGenerate(e) {
    e.preventDefault();

    const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ";
    const numeric = '0123456789';
    const punctuation = '!@%^&*_+~:;?,.-=';
    let password = "";

    while (password.length < this.state.passwordLength) {
      let letterIndex = Math.floor(letters.length * Math.random());
      password += letters[letterIndex];

      if (this.state.includeNumbers) {
        let numberIndex = Math.floor(numeric.length * Math.random());
        password += numeric[numberIndex];
      }

      if (this.state.includeSymbols) {
        let puncIndex = Math.floor(punctuation.length * Math.random());
        password += punctuation[puncIndex];
      }
    }

    this.props.onPasswordGenerated(password);
  }

  render() {
    return (
      <div className="cp-modal">
        <a href="#generate-modal"
          className="waves-effect waves-light btn modal-trigger">
          <i className="material-icons left">loop</i>
          Generate
        </a>

        <div id="generate-modal" className="modal">
          <div className="modal-content">
            <h4>Generate Password</h4>
            <div>
              <p className="range-field">
                <label htmlFor="password-length">Password Length ({this.state.passwordLength})</label>
                <input
                  type="range"
                  id="password-length"
                  min="6"
                  max="25"
                  step="1"
                  value={this.state.passwordLength}
                  onChange={(e) => this.setState({ passwordLength: e.target.value })}
                />
              </p>

              <p>
                <label htmlFor="include-numbers">
                  <input
                    id="include-numbers"
                    type="checkbox"
                    className="filled-in"
                    checked={this.state.includeNumbers}
                    onChange={(e) => this.setState({ includeNumbers: e.target.checked })}
                  />
                  <span>Include Numbers</span>
                </label>
              </p>

              <p>
                <label htmlFor="include-symbols">
                  <input
                    id="include-symbols"
                    type="checkbox"
                    className="filled-in"
                    checked={this.state.includeSymbols}
                    onChange={(e) => this.setState({ includeSymbols: e.target.checked })}
                  />
                  <span>Include Symbols</span>
                </label>
              </p>
            </div>
          </div>
          <div className="modal-footer">
            <a href="#!" className="modal-close waves-effect waves-green btn-flat">Cancel</a>
            <a href="#!" className="modal-close waves-effect waves-green btn-flat" onClick={this.handleGenerate}>Generate</a>
          </div>
        </div>

      </div>
    );
  }
}
