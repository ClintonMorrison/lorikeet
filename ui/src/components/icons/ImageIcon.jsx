import React from 'react';

import './ImageIcon.scss';

export default function ImageIcon({ src }) {
  return (
    <img className="cp-image-icon graphic" src={src} alt="" />
  )
}