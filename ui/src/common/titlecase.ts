/**
 * To convert the first letter of the input word to Upper Case
 */
export const titleCase = (filterName: string) => {
  return filterName.replace(/\w\S*/g, function (word) {
    return word.charAt(0).toUpperCase() + word.substr(1).toLowerCase();
  });
};
