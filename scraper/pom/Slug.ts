export function slugify(str: string) {
  str = str.replace(/^\s+|\s+$/g, ''); // trim
  str = str.toLowerCase();

  // remove accents, swap ñ for n, etc
  var from = "áäâàãåčçćďéěëèêẽĕȇğíìîïıňñóöòôõøðřŕšşťúůüùûýÿžþÞĐđßÆa·/_,:;";
  var to   = "aaaaaacccdeeeeeeeegiiiiinnooooooorrsstuuuuuyyzbBDdBAa------";
  for (var i = 0, l = from.length; i < l; i++) {
    str = str.replace(new RegExp(from.charAt(i), 'g'), to.charAt(i));
  }

  str = str.replace(/[^a-z0-9 -]/g, '') // remove invalid chars
           .replace(/\s+/g, '-') // collapse whitespace and replace by -
           .replace(/-+/g, '-'); // collapse dashes

  return str;
}

export function dateToString(date: Date) {
  let year = date.getFullYear().toString();
  if (year == "2001") {
    const today = new Date();
    year = today.getFullYear().toString();
  }

  let month = (date.getMonth() + 1).toString();
  if (date.getMonth() < 9) {
    month = `0${month}`;
  }
  let day = date.getDate().toString();
  if (date.getDate() < 10) {
    day = `0${day}`;
  }
  return `${year}-${month}-${day}`;
}

export function escapeString(str: string) {
  return str.replace(/'/g, `''`);
}
