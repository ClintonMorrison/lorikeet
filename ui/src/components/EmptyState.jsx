import React from 'react';

import './EmptyState.scss';

export default function EmptyState({ img, title, subtitle}) {
  return (
    <div className="cp-empty-state">
      <img className="graphic" src={img} alt="" />

      <p className="title"><strong>{title}</strong></p>
      <p>{subtitle}</p>
    </div>
  )
}