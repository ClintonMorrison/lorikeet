import React from 'react';

import './ImageIcon.scss';

export default function ImageIcon({ src, small }) {
  return (
    <img className={`cp-image-icon graphic ${small ? 'small' : ''}`} src={src} alt="" />
  )
}