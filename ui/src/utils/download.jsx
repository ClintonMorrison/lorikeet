import _ from 'lodash';

const fields = [
  'title', 'username', 'email', 'password', 'website', 'notes',
  'created', 'updated', 'lastUsed'
];

function _formatPasswords(passwords) {
  console.log(passwords);
  return passwords.map(p => _.pick(p, fields));
}

function _downloadBlob(data, type, filename) {
  const blob = new Blob([data], { type });
  const e = document.createEvent('MouseEvents');
  const a = document.createElement('a');

  a.download = filename;
  a.href = window.URL.createObjectURL(blob);
  a.dataset.downloadurl =  [type, a.download, a.href].join(':');
  e.initMouseEvent('click', true, false, window, 0, 0, 0, 0, 0, false, false, false, false, 0, null);
  a.dispatchEvent(e)
}

function _formatVaueCSV(val) {
  if (val.includes(",")) {
    return `"${val.replace(/"/g, '""') || ''}"`;   
  }

  return val;
}

function _formatTextField(field, value) {
  return `${_.startCase(field)}: ${value}`
}

export function downloadAsJSON(document, filename){
  const passwords = _formatPasswords(document.passwords);
  const data = JSON.stringify(passwords, undefined, 4);
  _downloadBlob(data, 'text/json', filename);
}

export function downloadAsText(document, filename){
  const passwords = _formatPasswords(document.passwords);

  const rows = passwords.map(p => {
    return fields
      .map(field => _formatTextField(field, p[field]))
      .join('\n');
  });

  _downloadBlob(rows.join('\n\n-------\n\n'), 'text/plain', filename);
}

export function downloadAsCSV(document, filename){
  const passwords = _formatPasswords(document.passwords);
  console.log(passwords);

  const header = fields.map(_.startCase).join(', ');

  const rows = passwords.map(p => {
    return fields
      .map(field => _formatVaueCSV(p[field]))
      .join(', ');
  });

  _downloadBlob([header, ...rows].join('\n'), 'text/csv', filename);
}
