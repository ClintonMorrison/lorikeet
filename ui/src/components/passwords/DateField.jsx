import React from 'react';
import dayjs from 'dayjs';
import BasicField from "./BasicField";

export default function DateField({ title, value }) {
  const date = dayjs(value).format("MMMM D, YYYY");
  return (
    <BasicField title={title} value={date} />
  );
}