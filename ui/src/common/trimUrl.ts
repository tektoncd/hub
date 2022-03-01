/**
 * This trims the `/` character from url if it has as a last character
 */
const trimChar = `/`;

export const TrimUrl = (url: string) => {
  if (url === '') {
    return url;
  } else {
    return url.charAt(url.length - 1) !== trimChar ? url : url.slice(0, -1);
  }
};
