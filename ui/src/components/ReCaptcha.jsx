import { useRef, useEffect } from "react";
import ReCAPTCHA from "react-google-recaptcha";

const SITE_KEY = '6LdmALIiAAAAAK-Kgn9zd7ohPIUmC3K0QZfOI_e5';

// See https://www.npmjs.com/package/react-google-recaptcha
export default function ReCaptcha({ onChange, reset, darkMode }) {
  const recaptchaRef = useRef();

  useEffect(() => {
    if (reset && recaptchaRef.current) {
      recaptchaRef.current.reset();
    }
  }, [reset]);

  return (
    <div className="input-field">

      <ReCAPTCHA
        sitekey={SITE_KEY}
        onChange={onChange}
        ref={recaptchaRef}
        theme={darkMode ? 'dark' : 'light'}
      />
    </div>
  )
}